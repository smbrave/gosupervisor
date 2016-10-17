package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
)

type Program struct {
	Name        string `xml:"name,attr"`
	Command     string `xml:"command"`
	Environment string `xml:"environment"`
	Directory   string `xml:"directory"`
	StderrFile  string `xml:"stderrfile"`
	StdoutFile  string `xml:"stdoutfile"`
	CallBackUrl string `xml:"callback"`
	StartSec    int    `xml:"startsec"`
}

type GoSupervisorConf struct {
	Programs []*Program `xml:"program"`
}

var gosupervisorconf GoSupervisorConf

func loadProgram() {
	conf_file := *flag_conf
	data, err := ioutil.ReadFile(conf_file)

	if err != nil {
		log.Printf("load conf[%s] err:%s", conf_file, err.Error())
		return
	}

	err = xml.Unmarshal(data, &gosupervisorconf)
	if err != nil {
		log.Printf("parse conf[%s] err:%s", conf_file, err.Error())
		return
	}

	del := make(map[string]bool)
	for _, proc := range procs {
		del[proc.Name] = true
	}

	for _, program := range gosupervisorconf.Programs {
		if program.Command == "" {
			continue
		}

		proc, ok := procs[program.Name]
		if !ok {
			proc = newProc()
		}

		//验证配置是否更改
		data, _ := json.Marshal(&program)
		digest := HexDigest(string(data))
		if digest == proc.Digest {
			continue
		}

		proc.Digest = digest
		setProc(proc, program)
		delete(del, proc.Name)
		proc.start()
		procs[proc.Name] = proc
		log.Printf("load: %+v", program)
	}

	//停止删除配置的程序
	for name, _ := range del {
		proc, ok := procs[name]
		if !ok {
			continue
		}
		proc.stop()
		delete(procs, name)
	}
}

//初始化程序
func setProc(proc *Proc, program *Program) {
	proc.Name = program.Name
	proc.Command = program.Command
	proc.Directory = program.Directory
	proc.StderrFile = program.StderrFile
	proc.StdoutFile = program.StdoutFile
	proc.CallBackUrl = program.CallBackUrl
	proc.StartSec = program.StartSec

	//设置默认值
	if proc.StartSec == 0 {
		proc.StartSec = 10
	}
}

//摘要
func HexDigest(conf string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(conf))
	cipherStr := md5Ctx.Sum(nil)
	digest := hex.EncodeToString(cipherStr)
	return digest[0:16]
}
