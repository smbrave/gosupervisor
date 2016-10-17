package main

import (
	"fmt"
	"net"
	"net/rpc"
)

type GoSupervisor int

func (r *GoSupervisor) Start(proc string, ret *string) (err error) {
	*ret = "start ok"
	return nil
}

func (r *GoSupervisor) Stop(proc string, ret *string) (err error) {
	*ret = "stop ok"
	return nil
}

func (r *GoSupervisor) Reload(proc string, ret *string) (err error) {
	*ret = "reload ok"
	return nil
}

func (r *GoSupervisor) Status(proc string, ret *string) (err error) {
	*ret = "status ok"
	return nil
}

func (r *GoSupervisor) List(proc string, ret *string) (err error) {
	result := ""
	for name, proc := range procs {
		result += fmt.Sprintf("proc:%s status:%s start:%s\n",
			name, proc.status(), proc.StartTime.Format("2006-01-02 15:04:05"))
	}
	*ret = result
	return nil
}

func startServer() {
	gs := new(GoSupervisor)
	rpc.Register(gs)
	server, err := net.Listen("tcp", *flag_listen)
	if err != nil {
		panic(err.Error())
	}
	for {
		client, err := server.Accept()
		if err != nil {
			continue
		}
		rpc.ServeConn(client)
	}
}
