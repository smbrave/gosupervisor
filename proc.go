package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type Proc struct {
	Name        string
	Command     string
	Args        string
	Environmet  string
	Directory   string
	CallBackUrl string
	StartSec    int

	Cmd        *exec.Cmd
	StartTime  time.Time
	Status     string
	FirstStart bool
	ErrChan    chan error
	ExitChan   chan struct{}
}

var procs map[string]*Proc

func startProc() {
	for name, proc := range procs {
		proc.Name = name
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

func (proc *Proc) run() {
	cmd := exec.Command(proc.Command, proc.Args)
	cmd.Env = append(os.Environ(), proc.Environmet)
	cmd.Dir = proc.Directory
	proc.Cmd = cmd
	proc.StartTime = time.Now()
	proc.Status = "starting"
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
		if proc.FirstStart {
			proc.FirstStart = false
		} else {
			proc.callback(true)
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
	proc.callback(false)
}

//todo
func (proc *Proc) callback(start bool) {
	if proc.CallBackUrl == "" {
		return
	}
	log.Printf("proc:%s callback: %s %v", proc.Name, proc.CallBackUrl, start)
}

func (proc *Proc) stop() error {
	p := proc.Cmd.Process
	if p == nil {
		return nil
	}

	pgid, err := syscall.Getpgid(p.Pid)
	if err != nil {
		return err
	}

	pid := p.Pid
	if pgid == p.Pid {
		pid = -1 * pid
	}

	target, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return target.Signal(syscall.SIGHUP)
}

func (proc *Proc) status() string {
	return proc.Status
}
