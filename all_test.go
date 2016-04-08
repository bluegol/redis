package redis

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

const host = "10.44.44.4"
const port = 36379
const numTry = 1
var numJobsArr = [...]int{1000000}
var numClientGoRoutines = [...]int{1000, 10000, 20000}

func TestTime(t *testing.T) {

	/*var numJobsArr = [...]int{100000}
	var numGoRoutinesArr = [...]int{1000}
	*/
	const numConn = 5
	testFunc := testPL6
	clientFunc := NewMuxClient

	var totalTimes [len(numJobsArr)][len(numClientGoRoutines)][numTry]float64
	for jjj, numJobs := range numJobsArr {
		for jj, numGoR := range numClientGoRoutines {

			for j := 0; j < numTry; j++ {
				fmt.Printf("jobs: %d    gor: %d    try: %d\n", numJobs, numGoR, j)
				c, err := clientFunc(host, port, "", numConn, DefaultStatusHandler)
				if err != nil {
					t.Error(err.Error())
					return
				}

				_, err, job := c.Command(CmdPing)
				if err != nil {
					t.Error(err.Error())
				} else {
					fmt.Printf("%s\n", job.DebugString())
				}

				_, err, job = c.Command(CmdSet, "TEST", 0)
				if err != nil {
					t.Error(err.Error())
				} else {
					fmt.Printf("%s\n", job.DebugString())
				}

				before := time.Now()
				var wg sync.WaitGroup
				perGoRoutine := numJobs / numGoR
				for i := 0; i < numGoR; i++ {
					wg.Add(1)
					go testFunc(t, &wg, c, i * perGoRoutine, perGoRoutine)
				}
				wg.Wait()
				dur := time.Now().Sub(before)
				c.Quit()
				fmt.Printf("all done. %f\n", dur.Seconds())
				totalTimes[jjj][jj][j] = dur.Seconds()
			}
		}
	}

	fmt.Println("final result:")
	for jjj, numJobs := range numJobsArr {
		for jj, numGoR := range numClientGoRoutines {
			fmt.Printf("%d,%d:", numJobs, numGoR)
			for j := 0; j < numTry; j++ {
				fmt.Printf("\t%f", totalTimes[jjj][jj][j])
			}
			fmt.Println("")
		}
	}
}

/*
func testSet(
t *testing.T, wg *sync.WaitGroup,
cli Client, startNo, count int ) {
	i := 0
	defer func() {
		fmt.Printf("done %d...%d\n", startNo, startNo+i-1)
		wg.Done()
	}()
	errCount := 0
	for ; i < count; i++ {
		job, reply := cli.Commands(
			CmdSet, startNo+i, startNo+i )
		if b, errCode := IsCriticalError(reply); b {
			t.Errorf( "critical error: %d %s, job: %s",
				errCode, ErrorString(reply), job.DebugString() )
			errCount++
		} else if job.CmdCount() != len(reply) {
			t.Errorf("reply mismatch1! %d %d", job.CmdCount(), len(reply))
			errCount++
		} else if reply[0].Type != MarkSimpleString {
			t.Errorf("unexpected reply. %s --> %s", job.DebugString(), reply[0].DebugString())
			errCount++
		}
		if errCount >= 2 {
			return
		}
	}
}

func testPL3(
t *testing.T, wg *sync.WaitGroup,
cli Client, startNo, count int ) {
	i := 0
	defer func() {
		fmt.Printf("done %d...%d\n", startNo, startNo+i-1)
		wg.Done()
	}()
	errCount := 0
	for ; i < count; i++ {
		job, reply := cli.Commands(
			CmdSet, startNo+i, startNo+i,
			CmdIncr, startNo+i,
			CmdGet, startNo+i )
		if b, errCode := IsCriticalError(reply); b {
			t.Errorf( "critical error: %d %s, job: %s",
				errCode, ErrorString(reply), job.DebugString() )
			errCount++
		} else if job.CmdCount() != len(reply) {
			t.Errorf("reply mismatch1! %d %d", job.CmdCount(), len(reply))
			errCount++
		} else if reply[0].Type != MarkSimpleString {
			t.Errorf("unexpected reply. %s --> %s", job.DebugString(), reply[0].DebugString())
			errCount++
		} else {
			r, err := reply[2].Int()
			if err != nil {
				t.Errorf("unexpected reply. %s --> %s", job.DebugString(), reply[2].DebugString())
				errCount++
			} else if r != startNo + i + 1 {
				t.Errorf("unexpected result %d. expecied %d", r, startNo + i + 1)
				errCount++
			}
		}
		if errCount >= 2 {
			return
		}
	}
}
*/

func testPL6(
t *testing.T, wg *sync.WaitGroup,
c Client, startNo, count int ) {
	i := 0
	defer func() {
		//fmt.Printf("done %d...%d\n", startNo, startNo+i-1)
		wg.Done()
	}()
	errCount := 0
	for ; i < count; i++ {
		_, err, job := c.Commands(
			CmdSet, startNo+i, startNo+i,
			CmdIncr, startNo+i,
			CmdGet, startNo+i,
			CmdHset, "h"+strconv.Itoa(startNo+i), "f"+strconv.Itoa(i), -startNo-2*i,
			CmdHincrby, "h"+strconv.Itoa(startNo+i), "f"+strconv.Itoa(i), 1,
			CmdHget, "h"+strconv.Itoa(startNo+i), "f"+strconv.Itoa(i) )
		if err != nil {
			t.Error(err.Error())
			return
		} else if job.Reply[0].Type != MarkSimpleString {
			t.Errorf("unexpected reply. %s", job.DebugString() )
			errCount++
		} else {
			r, err := job.Reply[2].Int()
			if err != nil {
				t.Errorf("unexpected reply. %s", job.DebugString())
				errCount++
			} else if r != startNo + i + 1 {
				t.Errorf("unexpected result %d. expecied %d", r, startNo + i + 1)
				errCount++
			}
		}
		if errCount >= 2 {
			return
		}
	}
}
