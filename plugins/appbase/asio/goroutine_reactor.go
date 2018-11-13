package asio

import "github.com/eapache/channels"

const Infinity = true

type GoroutineReactor struct {
	opq      chan operation
	opQueue  *channels.InfiniteChannel
	sigQueue chan operation
	down     chan struct{}
}

func NewGoroutineReactor() *GoroutineReactor {
	r := new(GoroutineReactor)

	if Infinity {
		r.opQueue = channels.NewInfiniteChannel()
	} else {
		r.opq = make(chan operation, 128)
	}
	r.sigQueue = make(chan operation, 1)
	r.down = make(chan struct{}, 1)
	return r
}

// use GoroutineReactor for each operating system
func (i *IoContext) GetService() ReactorService {
	if i.service == nil {
		i.service = NewGoroutineReactor()
	}
	return i.service
}

func (g *GoroutineReactor) run() {
	if Infinity {
	LP1:
		for {
			select {
			case <-g.down:
				break LP1

			case op := <-g.sigQueue:
				op.call()

			case op := <-g.opQueue.Out():
				op.(operation).call()

			default:

			}
		}
	} else {
	LP2:
		for {
			select {
			case <-g.down:
				break LP2

			case op := <-g.sigQueue:
				op.call()

			case op := <-g.opq:
				op.call()

			default:

			}
		}
	}

}

func (g *GoroutineReactor) stop() {
	g.down <- struct{}{}
}

func (g *GoroutineReactor) post(op operation) {
	//g.opQueue <- operation{op, args}
	if Infinity {
		g.opQueue.In() <- op
	} else {
		g.opq <- op
	}
}

func (g *GoroutineReactor) notify(op operation) {
	g.sigQueue <- op
}
