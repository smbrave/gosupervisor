package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func init_server() {
	//日志输出
	f, err := os.OpenFile(*flag_log, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0644)
	if err != nil {
		panic(fmt.Sprintf("%s", err.Error()))
	}
	log.SetOutput(f)

	//全局监控的进程
	procs = make(map[string]*Proc)
}

var (
	flag_server = flag.Bool("server", false, "gosupervisor run server")
	flag_log    = flag.String("log", "/var/log/gosupervisor.log", "gosupervisor log file")
	flag_conf   = flag.String("conf", "/etc/gosupervisor.conf", "gosupervisor config file")
	flag_listen = flag.String("listen", "127.0.0.1:33870", "gosupgervisor listen socket")
)

func main() {
	flag.Parse()
	if *flag_server {
		init_server()
		loadProgram()
		startServer()
		select {}
	}

	startClient()

}
