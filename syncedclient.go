package redis

import (
	"bytes"
	"bufio"
	"net"

	"github.com/bluegol/errutil"
)

type syncedClient struct {
	host     string
	port     int
	passwd   string

	conn     net.Conn
	reader   *bufio.Reader
	buf      *bytes.Buffer
	statusCh chan error

	hasQuit  bool
}

func NewSyncedClient(host string, port int, passwd string, dummy int,
    reporter func(error) ) (*syncedClient, error) {
	c := &syncedClient {
		host: host,
		port: port,
		passwd: passwd,

		statusCh: make(chan error, jobChSize),
	}

	go statusHandler(c.statusCh, reporter)

	err := c.newConn()
	if err != nil {
		return nil, err
	}
	c.buf = bytes.NewBuffer(nil)
	c.buf.Grow(cmdBytesBufferSize)

	return c, nil
}

// Command execute a single command and returns the result.
// guaranteed: err == nil <=> reply != nil
func (c *syncedClient) Command(cmd Cmd, args ...interface{}) (*Reply, error, *Job) {
	job, err := newSingleCmdJobNoCh(cmd, args...)
	if err != nil {
		return nil, err, nil
	}
	err = c.sendAndReceive(job)
	var reply *Reply
	if len(job.Reply) > 0 {
		reply = job.Reply[0]
	}
	return reply, err, job
}

func (c *syncedClient) Commands(args ...interface{}) ([]*Reply, error, *Job) {
	job, err := newJobNoCh(args...)
	if err != nil {
		return nil, err, nil
	}
	err = c.sendAndReceive(job)
	return job.Reply, err, job
}

func (c *syncedClient) Quit() error {
	if c.hasQuit {
		return InfoAlreadyQuit
	}
	c.conn.Write(cmdQuitBytes)
	c.conn.Close()
	c.conn = nil
	c.reader = nil
	c.hasQuit = true
	close(c.statusCh)
	return nil
}



//////////////////////////////////////////////////////////////////////

func (c *syncedClient) newConn() error {
	var err error
	c.conn, err = newConn(c.host, c.port, c.passwd, c.statusCh)
	if err != nil {
		return err
	}
	c.reader = bufio.NewReader(c.conn)
	return nil
}

func (c *syncedClient) sendAndReceive(job *Job) error {
	failCount := 0
	BeginProcess:
	err := send(c.conn, job.b)
	if err != nil {
		if errutil.CompareType(err, ErrBeforeSending) {
			c.statusCh <- err
			c.conn.Close()
			failCount++
			if failCount > 5 {
				job.Err = err
				return err
			}
			err = c.newConn()
			if err != nil {
				job.Err = err
				return err
			}
			goto BeginProcess
		} else {
			job.Err = err
			return err
		}
	}
	err = readReply(c.reader, job)
	return err
}
