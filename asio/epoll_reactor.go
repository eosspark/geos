//+build linux

package gosio

type epollReactor struct {

}
func (epollReactor) run() {
	//syscall.Epoll
}