package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

const APP = "gosupervisor"

var binaryVersion string
var buildTime string
var svnRevision string

func version() string {
	return fmt.Sprintf("%s v%s (built:%s git:%s %s)", APP, binaryVersion, buildTime, svnRevision, runtime.Version())
}

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

func init_singal() {
	signalChan := make(chan os.Signal, 1)
	signalIgnoreChan := make(chan os.Signal, 1)
	go func() {
		for {
			select {
			case <-signalChan:
				for _, proc := range procs {
					proc.stop()
					log.Printf("name:%s stoped!", proc.Name)
				}
				log.Printf("all proc stop!")
				os.Exit(0)
				return
			case <-signalIgnoreChan:
				continue
			}
		}
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(signalIgnoreChan, syscall.SIGPIPE)
}

var (
	flag_server  = flag.Bool("server", false, "gosupervisor run server")
	flag_log     = flag.String("log", "/var/log/gosupervisor.log", "gosupervisor log file")
	flag_conf    = flag.String("conf", "/etc/gosupervisor.xml", "gosupervisor config file")
	flag_listen  = flag.String("listen", "127.0.0.1:33870", "gosupgervisor listen socket")
	flag_version = flag.Bool("v", false, "print version")
)

func main() {
	flag.Parse()
	if *flag_version {
		fmt.Printf("%s\n", version())
		return
	}

	if *flag_server {
		init_singal()
		init_server()
		loadProgram()
		startServer()
		select {}
	}

	startClient()

}
