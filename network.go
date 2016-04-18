package redis

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/bluegol/errutil"
)

func newConn(host string, port int, passwd string, statusCh chan<- error ) ( net.Conn, error ) {
	for i := 0; i < numTryConn; i++ {
		conn, err := net.DialTimeout(
			"tcp", fmt.Sprintf("%s:%d", host, port), timeout)
		if err == nil {
			// set keepalive so that Write fails when conn is disconnected
			var tcpConn = conn.(*net.TCPConn)
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(10 * time.Second)

			if len(passwd) > 0 {
				// \todo send and check passwd
			}

			return conn, nil
		}
		time.Sleep(sleepBetweenRetryConn)
	}
	return nil, ErrCannotConnect
}

const numTryConn = 5
const timeout = 10 * time.Second
const sleepBetweenRetryConn = 2 * time.Second
const cmdBytesBufferSize = 4096

func send(conn net.Conn, b []byte) error {
	// so do I have to do this every time? it's very annoying!
	conn.SetWriteDeadline(time.Now().Add(timeout))
	n, err := conn.Write(b)
	if err != nil {
		if n == 0 {
			return ErrBeforeSending
		} else {
			return ErrWhileSending
		}
	}

	return nil
}

func readReply(r *bufio.Reader, job *Job) error {
	reply := make([]*Reply, job.numCmds)
	for i := 0; i < job.numCmds; i++ {
		r, err := parseReply(r)
		if err != nil {
			// even if network error occurs, it will be handled in sender.
			// but this is critical, so report it.
			// but continue. reconnecting is done at other places
			job.Err = errutil.Embed(ErrCannotRead, err,
				errutil.MoreInfo,
				fmt.Sprintf("error while reading reading %d-th reply", i+1),
				"job", job.DebugString() )
			return job.Err
		} else {
			reply[i] = r
			err := r.Error()
			if err != nil && job.Err == nil {
				job.Err = err
			}
		}
	}
	job.Reply = reply
	return job.Err
}

func readReplyAndReturn(reader *bufio.Reader, job *Job) error {
	readReply(reader, job)
	if job.resultCh == nil {
		return errutil.NewAssert(errutil.MoreInfo, "job resultCh is nil")
	}
	job.resultCh <- job
	close(job.resultCh)
	return job.Err
}
