package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type Proc struct {
	Name        string
	Command     string
	Environment string
	Directory   string
	StderrFile  string
	StdoutFile  string
	CallBackUrl string
	StartSec    int

	Cmd        *exec.Cmd
	StartTime  time.Time
	OpTime     time.Time
	Status     string
	FirstStart bool
	Digest     string
	ErrChan    chan error
	ExitChan   chan struct{}
}

var procs map[string]*Proc

func startProc() {
	for _, proc := range procs {
		proc.start()
	}
}

func newProc() *Proc {
	return &Proc{
		ErrChan:    make(chan error),
		ExitChan:   nil,
		FirstStart: true,
		StartSec:   10,
	}
}

func (proc *Proc) start() {
	if proc.ExitChan != nil {
		log.Printf("proc name:%s already start ...", proc.Name)
		return
	}
	proc.OpTime = time.Now()
	proc.ExitChan = make(chan struct{})
	go func() {
		for {
			proc.run()
			ticker := time.NewTicker(time.Second)
			select {
			case <-ticker.C:
			case <-proc.ExitChan:
				goto exit
			}
		}
	exit:
		proc.Status = "stoped"
		proc.ExitChan = nil
	}()
}

func (proc *Proc) init() {
	tmp := strings.Split(proc.Command, " ")
	cmd := exec.Command(tmp[0], tmp[1:]...)
	cmd.Dir = proc.Directory

	//环境变量
	if proc.Environment != "" {
		envs := strings.Split(proc.Environment, ";")
		cmd.Env = os.Environ()
		for _, env := range envs {
			cmd.Env = append(cmd.Env, env)
		}
	}

	//标准输出
	if proc.StdoutFile != "" {
		path := ""
		if strings.HasPrefix(proc.StdoutFile, "/") {
			path = proc.StdoutFile
		} else {
			path = proc.Directory + "/" + proc.StdoutFile
		}
		stdout, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0644)
		if err != nil {
			log.Printf("open file:%s err:%s", proc.StdoutFile, err.Error())
		} else {
			cmd.Stdout = stdout
		}
	}

	//标准错误
	if proc.StderrFile != "" {
		path := ""
		if strings.HasPrefix(proc.StderrFile, "/") {
			path = proc.StderrFile
		} else {
			path = proc.Directory + "/" + proc.StderrFile
		}
		stderr, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0644)
		if err != nil {
			log.Printf("open file:%s err:%s", proc.StderrFile, err.Error())
		} else {
			cmd.Stderr = stderr
		}
	}

	proc.Cmd = cmd
}

func (proc *Proc) run() {
	proc.init()
	proc.StartTime = time.Now()
	proc.Status = "starting"

	cmd := proc.Cmd
	err := cmd.Start()
	if err != nil {
		log.Printf("start proc:%s err:%s", proc.Name, err.Error())
		return
	}
	go func() {
		proc.ErrChan <- cmd.Wait()
	}()

	ticker := time.NewTicker(time.Duration(proc.StartSec) * time.Second)
	defer ticker.Stop()

	select {
	case <-ticker.C:
		//首次启动不回调
		if proc.FirstStart {
			proc.FirstStart = false
		} else {
			proc.callback("start")
		}
	case err = <-proc.ErrChan:
		log.Printf("start proc:%s less than %ds", proc.Name, proc.StartSec)
		return
	}

	proc.Status = "running"
	log.Printf("proc:%s start ...", proc.Name)
	select {
	case err = <-proc.ErrChan:
		log.Printf("proc:%s exit, start at:%s, run span:%v, err:%v",
			proc.Name, proc.StartTime.Format("2006-01-02 15:04:05"), time.Since(proc.StartTime).String(), err)
	}
	proc.callback("stop")
}

//回调告警
func (proc *Proc) callback(status string) {
	if proc.CallBackUrl == "" {
		return
	}

	//人工操作的不回调
	span := time.Since(proc.OpTime).Seconds()
	if span < float64(proc.StartSec+10) {
		return
	}

	client := http.Client{
		Timeout: 10 * time.Second,
	}
	callbackData := make(map[string]interface{})
	callbackData["starttime"] = proc.StartTime.Format("2006-01-02 15:04:05")
	callbackData["sendtime"] = time.Now().Format("2006-01-02 15:04:05")
	callbackData["procname"] = proc.Name
	callbackData["command"] = proc.Command
	callbackData["status"] = status
	callbackData["hostname"] = hostname
	callbackBody, _ := json.Marshal(callbackData)
	resp, err := client.Post(proc.CallBackUrl, "application/json;charset=utf-8", bytes.NewBuffer(callbackBody))
	if err != nil {
		log.Printf("proc:%s callback url:%s err:%s", proc.Name, proc.CallBackUrl, err.Error())
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("proc:%s callback url:%s err:%s", proc.Name, proc.CallBackUrl, err.Error())
		return
	}

	log.Printf("proc:%s status:%s callback: %s %s",
		proc.Name, status, proc.CallBackUrl, string(data))
}

func (proc *Proc) stop() error {
	p := proc.Cmd.Process
	if p == nil {
		return nil
	}

	target, err := os.FindProcess(p.Pid)
	if err != nil {
		return err
	}
	if proc.ExitChan != nil {
		close(proc.ExitChan)
	}
	proc.OpTime = time.Now()
	return target.Signal(syscall.SIGTERM)
}

func (proc *Proc) restart() error {
	defer func() { proc.OpTime = time.Now() }()
	if proc.ExitChan == nil {
		proc.start()
		return nil
	}

	p := proc.Cmd.Process
	if p == nil {
		proc.start()
		return nil
	}

	target, err := os.FindProcess(p.Pid)
	if err != nil {
		proc.start()
		return err
	}

	return target.Signal(syscall.SIGTERM)

}

func (proc *Proc) kill() error {
	defer func() { proc.OpTime = time.Now() }()
	p := proc.Cmd.Process
	if proc.ExitChan != nil {
		close(proc.ExitChan)
	}
	return p.Kill()
}

func (proc *Proc) status() string {
	return proc.Status
}
