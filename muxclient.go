package redis

import (
	"bufio"
	"net"
	"sync"

	"github.com/bluegol/errutil"
)

type MuxClient struct {
	host     string
	port     int
	passwd   string
	numConn  int

	jobCh    chan *Job
	quitChs  []chan struct{}
	statusCh chan error

	// wg decrease when receiver returns
	wg       *sync.WaitGroup
	// q is used to quit once
	q        *sync.Once
}

func NewMuxClient(host string, port int, passwd string, numConn int,
	reporter func(error) ) (*MuxClient, error) {
	c := &MuxClient{
		host: host,
		port: port,
		passwd: passwd,
		numConn: numConn,

		jobCh: make(chan *Job, jobChSize),
		quitChs: make([]chan struct{}, numConn),
		statusCh: make(chan error, jobChSize),

		wg: new(sync.WaitGroup),
		q: new(sync.Once),
	}

	go statusHandler(c.statusCh, reporter)
	for i := 0; i < c.numConn; i++ {
		startupCh := make(chan error)
		quitCh := make(chan struct{})
		go c.sender(i, startupCh, quitCh)
		err := <- startupCh
		if err != nil {
			c.Quit()
			return nil, err
		}
		c.quitChs[i] = quitCh
	}

	return c, nil
}

func (c *MuxClient) Command(cmd Cmd, args ...interface{}) (*Reply, error, *Job) {
	return command(c.jobCh, cmd, args...)
}

func (c *MuxClient) Commands(args ...interface{}) ([]*Reply, error, *Job) {
 	return commands(c.jobCh, args...)
}

func (c *MuxClient) Quit() error {
	result := InfoAlreadyQuit
	c.q.Do(func() {
		c.statusCh <- InfoQuitStarted
		for _, q := range c.quitChs {
			if q != nil {
				q <- struct{}{}
				close(q)
			}
		}
		c.wg.Wait()
		c.statusCh <- InfoQuitDone
		close(c.statusCh)
		go func() {
			for job := range c.jobCh {
				job.Err = ErrAlreadyClosed
				job.sendResult()
			}
		}()
		result = nil
	})
	return result
}



/////////////////////////////////////////////////////////////////////

func (c *MuxClient) sender(senderNo int,
	startupCh chan error, quitCh chan struct{}) {

	conn, err := newConn(c.host, c.port, c.passwd, c.statusCh)
	if err != nil {
		startupCh <- err
		close(startupCh)
		return
	}
	waitCh := newReceiver(conn, c.wg, c.statusCh)
	close(startupCh)
	c.statusCh <- InfoConnected

	for {
		select {

		case <-quitCh:
			// send quit anyway, and don't care about response.
			conn.Write(cmdQuitBytes)
			close(waitCh)
			return

		case job := <-c.jobCh:
			ProcessJob:
			err = send(conn, job.b)
			if err == nil {
				waitCh <- job
				continue
			}

			// error occurred
			c.statusCh <- errutil.AddInfo(err, "job", job.DebugString())
			// stop the receiver
			close(waitCh)
			if errutil.CompareType(ErrBeforeSending, err) {
				// nothing was sent, so it could be timeout due to
				// long period of no communication
				var err2 error
				conn, err2 = newConn(
					c.host, c.port, c.passwd, c.statusCh)
				if err == nil {
					// let's resume
					waitCh = newReceiver(conn, c.wg, c.statusCh)
					goto ProcessJob
				} else {
					// cannot reconnect. critical. start shutdown.
					c.statusCh <- err2
					job.Err = errutil.Embed(ErrCannotConnect, err2)
					job.sendResult()
				}
			} else {
				job.Err = err
				job.sendResult()
			}

			go c.Quit()
			<-quitCh
			return
		}
	}
}

func newReceiver(conn net.Conn, wg *sync.WaitGroup, statusCh chan<- error) chan *Job {
	waitCh := make(chan *Job, waitChSize)
	go receiver(conn, waitCh, wg, statusCh)
	return waitCh
}

func receiver(conn net.Conn, waitCh <-chan *Job, wg *sync.WaitGroup, statusCh chan<- error) {

	wg.Add(1)
	defer wg.Done()

	reader := bufio.NewReader( conn )

	for job := range waitCh {
		err := readReplyAndReturn(reader, job)
		if err != nil {
			// critical, but sender will handle it.
			// just report it, and continue
			statusCh <- err
		}
	}
	// channel closed, meaning we are exiting.
	// just close and don't care about the error
	conn.Close()
}

const jobChSize = 10000
const waitChSize = 10000
