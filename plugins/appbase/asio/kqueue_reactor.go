//+build darwin
//not use

package asio

import (
	"syscall"
	"fmt"
	"strings"
	"errors"
	"github.com/eosspark/eos-go/plugins/appbase/asio/utils"
	)

type KqueueReactor struct {
	kqueueFd int
	shutdown bool
}

func NewKqueueReactor() *KqueueReactor {
	k := new(KqueueReactor)
	k.kqueueFd = doKqueueCreate()
	k.shutdown = false
	return k
}

func doKqueueCreate() int {
	kq, kqErr := syscall.Kqueue()
	if kqErr != nil {
		fmt.Println("Error creating Kqueue descriptor!", kqErr)
		// TODO: kq err
	}

	return kq
}


func (k *KqueueReactor) run() {

	fd, fdErr := CreateListenSocket("127.0.0.1:8080")
	if fdErr != nil {
		fmt.Println("Error create socket", fdErr)
		// TODO
	}

	change := syscall.Kevent_t{
		Ident: uint64(fd),
		Filter: syscall.EVFILT_READ,
		Flags: syscall.EV_ADD | syscall.EV_ENABLE,
		Fflags: 0,
		Data: 0,
		Udata: nil,
	}

	_, _err := syscall.Kevent(k.kqueueFd, []syscall.Kevent_t{change}, []syscall.Kevent_t{} , &syscall.Timespec{})
	if _err != nil {
		fmt.Println("Error subscribe events", _err)
		// TODO:
	}

	events := make([]syscall.Kevent_t, 10)
	timeout := syscall.Timespec{Sec: 10, Nsec: 0}

	nevent, evErr := syscall.Kevent(k.kqueueFd, []syscall.Kevent_t{}, events, &timeout)
	if evErr != nil {
		fmt.Println("Error getting kevents", evErr)
		// TODO:
	}

	fmt.Println("event num:", nevent)
	for i:=0; i<nevent; i++ {

	}
}


func CreateListenSocket(ipport string) (int, error) {
	socket, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	var flag = int(1)
	err := syscall.SetsockoptInt(socket, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, flag)
	if err != nil {
		fmt.Println("TcpServer Setsockopt failed")
		return 0, err
	}

	ipinfo := strings.Split(ipport, ":")
	if len(ipinfo) != 2 {
		fmt.Println("TcpServer invalid ipport:%s", ipport)
		return 0, errors.New("invalid ipport")
	}

	ip, err := utils.ParseIPv4(ipinfo[0])
	if err != nil {
		fmt.Println("TcpServer parseIPv4 failed")
		return 0, errors.New("invalid ip")
	}

	port, err := utils.ParsePort(ipinfo[1])
	if err != nil {
		fmt.Println("TcpServer parsePort failed")
		return 0, err
	}

	addr := &syscall.SockaddrInet4{
		//Family: syscall.AF_INET,
		Port: port,
		Addr: ip,
	}

	err = syscall.Bind(socket, addr)
	if err != nil {
		fmt.Println("TcpServer Bind failed")
		return 0, err
	}

	err = syscall.Listen(socket, 1024 * 100)
	if err != nil {
		fmt.Println("TcpServer listen failed")
		return 0, err
	}
	return socket, nil
}



