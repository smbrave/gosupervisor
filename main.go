package main

import (
	"fmt"
	"log"
	"os"
)

func init() {
	procs = make(map[string]*Proc)

}

//todo stdlog config rpc_controll

func main() {

	f, err := os.OpenFile("/var/log/gosupervisor.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0755)
	if err != nil {
		panic(fmt.Sprintf("%s", err.Error()))
	}
	os.Stdin = f
	os.Stdout = f
	log.SetOutput(f)

	proc1 := newProc()
	proc2 := newProc()

	proc1.Command = "./test.sh"
	proc1.Args = "1200"
	proc1.Directory = "/home/vagrant/golang/src/github.com/smbrave/gosupervisor"

	proc2.Command = "./test.sh"
	proc2.Args = "1300"
	proc2.Directory = "/home/vagrant/golang/src/github.com/smbrave/gosupervisor"

	procs["jiangyong1"] = proc1
	procs["jiangyong2"] = proc2
	startProc()
	for {

	}
}
