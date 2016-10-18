package main

import (
	"fmt"
	"net"
	"net/rpc"
	"time"
)

type GoSupervisor int

func (r *GoSupervisor) Start(args string, ret *string) error {
	proc, ok := procs[args]
	*ret = "OK"
	if !ok {
		return fmt.Errorf("proc[%s] not exist", proc)
	}
	proc.start()
	return nil
}

func (r *GoSupervisor) Stop(args string, ret *string) error {
	proc, ok := procs[args]
	*ret = "OK"
	if !ok {
		return fmt.Errorf("proc[%s] not exist", proc)
	}
	return proc.stop()
}

func (r *GoSupervisor) Kill(args string, ret *string) error {
	proc, ok := procs[args]
	*ret = "OK"
	if !ok {
		return fmt.Errorf("proc[%s] not exist", proc)
	}
	return proc.kill()
}

func (r *GoSupervisor) Restart(args string, ret *string) error {
	proc, ok := procs[args]
	*ret = "OK"
	if !ok {
		return fmt.Errorf("proc[%s] not exist", proc)
	}

	return proc.restart()

}

func (r *GoSupervisor) Status(args string, ret *string) error {
	proc, ok := procs[args]
	*ret = "OK"
	if !ok {
		return fmt.Errorf("proc[%s] not exist", proc)
	}
	*ret = proc.status()
	return nil
}

func (r *GoSupervisor) Reload(args string, ret *string) (err error) {
	loadProgram()
	*ret = "OK"
	return nil
}

func (r *GoSupervisor) List(args string, ret *string) error {
	result := ""
	for name, proc := range procs {
		status := proc.status()
		if status == "running" {
			result += fmt.Sprintf("name:%-20s status:%-10s pid:%-6d start:%-25s uptime:%-20s\n",
				name, status, proc.Cmd.Process.Pid, proc.StartTime.Format("2006-01-02 15:04:05"), time.Since(proc.StartTime).String())
		} else {
			result += fmt.Sprintf("proc:%-15s status:%-10s\n",
				name, status)
		}
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
