package redis

import (
	"sync"
	"time"
)

type ClientPool struct {
	host     string
	port     int
	passwd   string
	numConn  int

	clnts      []*RentableClient
	maxRentDur time.Duration
	rentCh     chan chan *RentableClient
	statusCh   chan error

	// wg decrease when controller returns
	wg       *sync.WaitGroup
	// q is used to quit once
	q        *sync.Once
}

func NewClientPool(host string, port int, passwd string, numConn int,
	reporter func(error) ) (*ClientPool, error) {

	p := &ClientPool{
		host: host,
		port: port,
		passwd: passwd,
		numConn: numConn,

		clnts: make([]*RentableClient, numConn),
		maxRentDur: maxRentDur,
	}
	var err error
	var i int = 0
	for i = 0; i < numConn; i++ {
		var c *syncedClient
		c, err = NewSyncedClient(host, port, passwd, 1, reporter)
		if err != nil {
			break
		}
		p.clnts[i] = &RentableClient{ c, make(chan struct{}), make(chan struct{}) }
	}
	if err != nil {
		// cleanup and return

		return nil, err
	}
	p.rentCh = make(chan chan *RentableClient)
	p.statusCh = make(chan error, jobChSize)
	p.wg = new(sync.WaitGroup)
	p.q = new(sync.Once)
	for _, clnt := range p.clnts {
		go controller(clnt, p.maxRentDur, p.wg, p.rentCh, p.statusCh)
	}

	return p, nil
}

func (p *ClientPool) Quit() error {
	result := InfoAlreadyQuit
	p.q.Do(func() {
		p.statusCh <- InfoQuitStarted
		for _, c := range p.clnts {
			c.quitCh <- struct{}{}
			close(c.quitCh)
		}
		p.wg.Wait()
		p.statusCh <- InfoQuitDone
		close(p.statusCh)
		result = nil
	})
	return result
}

func (p *ClientPool) Rent() *RentableClient {
	ch := make(chan *RentableClient)
	p.rentCh<- ch
	rc := <-ch
	return rc
}

type RentableClient struct {
	*syncedClient
	returnCh chan struct{}
	quitCh   chan struct{}
}

func (rc *RentableClient) Quit() error {
	return ErrRentedClientCannotQuit
}

func (rc *RentableClient) Return() {
	rc.returnCh<- struct{}{}
}



/////////////////////////////////////////////////////////////////////

const maxRentDur = 10 * time.Second

func controller(c *RentableClient, maxRentDur time.Duration,
	wg *sync.WaitGroup,
	rentCh <-chan chan *RentableClient,
	statusCh chan<- error ) {

	wg.Add(1)
	defer wg.Done()

	rented := false
	timedout := false
	errorOccurred := false
	mainLoop:
	for {
		if rented {
			// wait for return, timeout, or quit
			timeout := time.After(maxRentDur)
			for {
				select {

				case <-c.returnCh:
					rented = false
					if timedout {
						statusCh <- InfoReturnedAfterTimeout
					}
					continue mainLoop

				case <-timeout:
					// what should we do for timeout?
					statusCh <- ErrRentTimeout
					timedout = true
					break mainLoop

				case <-c.quitCh:
					statusCh <- ErrQuitWhileRented
					break mainLoop

				}
			}
		} else {
			// wait for job, rent request, or quit
			for {
				select {

				case ch := <-rentCh:
					rented = true
					timedout = false
					ch <- c
					continue mainLoop

				case <-c.quitCh:
					break mainLoop

				}
			}
		}
	}

	c.Quit()
	// if error occurred, wait till quit is signalled
	if errorOccurred {
		_ = <-c.quitCh
	}
}











