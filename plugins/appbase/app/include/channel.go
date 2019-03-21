package include

import "github.com/eosspark/eos-go/libraries/asio"

type Channel struct {
	iosPtr *asio.IoContext
	signal Signal
}

func NewChannel (io *asio.IoContext) *Channel {
	ch := new(Channel)
	ch.iosPtr = io
	return ch
}


/**
* Publish data to a channel.  This data is *copied* on publish.
* @param data - the data to publish
*/
func (s *Channel) Publish(data ...interface{}) {
	s.iosPtr.Post(func(err error) {
		s.signal.Emit(data...)
	})
}

/**
* subscribe to data on a channel
* @tparam Callback the type of the callback (functor|lambda)
* @param cb the callback
* @return handle to the subscription
*/
func (s *Channel) Subscribe(f Caller) {
	s.signal.Connect(f)
}
