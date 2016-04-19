package redis

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/bluegol/errutil"
)

type Job struct {
	singleCmd *Cmd
	args      []interface{}
	resultCh  chan *Job

	numCmds   int
	b         []byte
	Reply     []*Reply
	Err       error
}

func (job *Job) DebugString() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("cmd(s): ")
	if job.singleCmd != nil {
		buf.Write([]byte(*job.singleCmd))
		for _, a := range job.args {
			buf.WriteString(" / ")
			writeDebugString(buf, a)
		}
	} else {
		if len(job.args) == 0 {
			buf.WriteString("<invalid empty cmd>")
		} else {
			for i, a := range job.args {
				if i > 0 {
					buf.WriteString(" / ")
				}
				writeDebugString(buf, a)
			}
		}
	}
	if len(job.Reply) > 0 {
		buf.WriteString(" --> reply: ")
		for i, r := range job.Reply {
			if i > 0 {
				buf.WriteString(" / ")
			}
			buf.WriteByte(r.Type)
			buf.Write(r.bytes)
		}
	}
	if job.Err != nil {
		buf.WriteString(" --> error: ")
		buf.WriteString(job.Err.Error())
	}

	s := strconv.QuoteToASCII(string(buf.Bytes()))
	return s[1:len(s)-1]
}



//////////////////////////////////////////////////////////////////////

func newSingleCmdJob(c Cmd, args ...interface{}) (*Job, error) {
	job := &Job{ singleCmd: &c, args: args, resultCh: make(chan *Job, 1) }
	err := job.prepare(bytes.NewBuffer(nil))
	return job, err
}

func newJob(args ...interface{}) (*Job, error) {
	job := &Job{ args: args, resultCh: make(chan *Job, 1) }
	err := job.prepare(bytes.NewBuffer(nil))
	return job, err
}

func newSingleCmdJobNoCh(c Cmd, args ...interface{}) (*Job, error) {
	job := &Job{ singleCmd: &c, args: args }
	err := job.prepare(bytes.NewBuffer(nil))
	return job, err
}

func newJobNoCh(args ...interface{}) (*Job, error) {
	job := &Job{ args: args }
	err := job.prepare(bytes.NewBuffer(nil))
	return job, err
}

func (job *Job) prepare(buf *bytes.Buffer) error {
	if job.singleCmd != nil {
		err := writeBytes(buf, *job.singleCmd, job.args...)
		if err != nil {
			job.Err = err
			return err
		}
		job.b = buf.Bytes()
		job.numCmds = 1
		return nil
	} else {
		if len(job.args) == 0 {
			job.Err = errutil.New(ErrInvalidCmd,
				errutil.MoreInfo, "empty cmd is invalid")
			return job.Err
		}
		var currentArgs []interface{}
		var currentCmd Cmd
		i := 0
		var a interface{}
		for i, a = range job.args {
			c, ok := a.(Cmd)
			if ok {
				if i >= 1 {
					// make previous cmd
					err := writeBytes(buf, currentCmd, currentArgs...)
					if err != nil {
						job.Err = errutil.AddInfo(err, "argNo", strconv.Itoa(i))
						return job.Err
					}
					job.numCmds++
				}
				currentCmd = c
				currentArgs = nil
			} else {
				currentArgs = append(currentArgs, a)
			}
		}
		// the last cmd
		err := writeBytes(buf, currentCmd, currentArgs...)
		if err != nil {
			job.Err = errutil.AddInfo(err, "argNo", "LAST")
			return job.Err
		}
		job.numCmds++

		job.b = buf.Bytes()
		return nil
	}

}

func writeBytes(buf *bytes.Buffer, cmd Cmd, args ...interface{}) error {
	buf.Write(bytesMarkArray)
	buf.Write( []byte(strconv.Itoa(1 + len(args))) )
	buf.Write( bytesCrlf )
	buf.Write([]byte(cmd))
	// check validity of args, and write bytes
	for i, a := range args {
		var err error
		switch a.(type) {
		case string:
			err = writeBulkString(buf, []byte(a.(string)))
		case int:
			err = writeBulkString(buf, []byte(strconv.Itoa(a.(int))))
		case int64:
			err = writeBulkString(buf, []byte(strconv.FormatInt(a.(int64), 10)))
		case []byte:
			err = writeBulkString(buf, a.([]byte))
		default:
			return errutil.New(ErrInvalidArgType,
				errutil.MoreInfo,
				fmt.Sprintf("argNo: %d, type: %T", i+1, a) )
		}
		if err != nil {
			return errutil.AssertEmbed(err,
				errutil.MoreInfo, fmt.Sprintf("argNo: %d, type: %T", i+1, a) )
		}
	}

	return nil
}

func writeBulkString(wr io.Writer, b []byte) error {
	l := len(b)
	lenBytes := []byte(strconv.Itoa(l))
	var err error
	_, err = wr.Write(bytesMarkBulkString)
	if err != nil {
		return err
	}
	_, err = wr.Write(lenBytes)
	if err != nil {
		return err
	}
	_, err = wr.Write(bytesCrlf)
	if err != nil {
		return err
	}
	_, err = wr.Write(b)
	if err != nil {
		return err
	}
	_, err = wr.Write(bytesCrlf)
	if err != nil {
		return err
	}
	return nil
}

func (job *Job) sendResult() {
	job.resultCh <- job
	close(job.resultCh)
}

func (job *Job) sendAndWait(ch chan<- *Job) {
	ch <- job
	<- job.resultCh
}

func writeDebugString(buf *bytes.Buffer, a interface{}) {
	switch a.(type) {

	case string:
		buf.WriteString(a.(string))

	case int:
		buf.WriteString(strconv.Itoa(a.(int)))

	case int64:
		buf.WriteString(strconv.FormatInt(a.(int64), 10))

	case []byte:
		buf.Write(a.([]byte))

	case Cmd:
		buf.Write([]byte(a.(Cmd)))

	default:
		buf.WriteString(fmt.Sprintf("<invalid type %T>", a))

	}
}