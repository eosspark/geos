//+build linux

package asio

type epollReactor struct {

}
func (epollReactor) run() {
	//syscall.Epoll
}