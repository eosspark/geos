package asio

import "github.com/eapache/channels"

type GoroutineReactor struct {
	//opQueue  chan operation
	opQueue  *channels.InfiniteChannel
	sigQueue chan operation
	down 	 chan struct{}
}

func NewGoroutineReactor() *GoroutineReactor {
	r := new(GoroutineReactor)
	//r.opQueue = make(chan operation, 128)
	r.opQueue = channels.NewInfiniteChannel()
	r.sigQueue = make(chan operation, 1)
	r.down = make(chan struct{}, 1)
	return r
}

// use GoroutineReactor for each operating system
func (i *IoContext) GetService () ReactorService {
	if i.service == nil {
		i.service = NewGoroutineReactor()
	}
	return i.service
}

func (g *GoroutineReactor) run() {
LP:	for ;; {
		select {
		case <-g.down:
			break LP

		case sig := <-g.sigQueue:
			doReactor(sig)

		//case op := <-g.opQueue:
		//	doReactor(op)
		case op := <-g.opQueue.Out():
			doReactor(op.(operation))
		}
	}

}

func (g *GoroutineReactor) stop() {
	g.down <- struct{}{}
}

func (g *GoroutineReactor) post(op interface{}, args ...interface{}) {
	//g.opQueue <- operation{op, args}
	g.opQueue.In() <- operation{op, args}
}

func (g *GoroutineReactor) notify(op interface{}, args ...interface{}) {
	g.sigQueue <- operation{op, args}
}



