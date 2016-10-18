package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/rpc"
	"os"
	"strings"
)

func startClient() {

	client, err := rpc.Dial("tcp", *flag_listen)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}
	defer client.Close()

	reader := bufio.NewReader(os.Stdin)
	cmdline := ""

	for {
		fmt.Printf("gosupervisor>")
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		cmdline = strings.Trim(line, "\r\n\t")
		if cmdline == "" {
			continue
		}

		err, ret := run(client, cmdline)
		if err != nil {
			fmt.Println(err.Error())
			if err.Error() == "connection is shut down" {
				client, err = rpc.Dial("tcp", *flag_listen)
				if err != nil {
					fmt.Printf("%s\n", err.Error())
					return
				}
			}
		} else {
			fmt.Println(ret)
		}
	}
}

func run0(client *rpc.Client, cmd string) (error, string) {
	var ret string
	ret = "OK"
	switch cmd {
	case "reload":
		return client.Call("GoSupervisor.Reload", "", &ret), ret
	case "list":
		return client.Call("GoSupervisor.List", "", &ret), ret
	case "exit":
		os.Exit(0)
	case "help":
		ret = "list\t: list all program\n"
		ret += "reload\t: reload gosupervisor config\n"
		ret += "exit\t: exit gosupervisor\n"
		ret += "start\t: start program, eg: start procname\n"
		ret += "stop\t: stop program, eg: stop procname\n"
		ret += "kill\t: kill program, eg: kill procname\n"
		ret += "restart\t: restart program, eg: restart procname\n"
		return nil, ret

	}
	return errors.New("command error"), ""
}

func run1(client *rpc.Client, cmd string, args1 string) (error, string) {
	var ret string
	ret = "OK"
	switch cmd {
	case "start":
		return client.Call("GoSupervisor.Start", args1, &ret), ret
	case "stop":
		return client.Call("GoSupervisor.Stop", args1, &ret), ret
	case "kill":
		return client.Call("GoSupervisor.Kill", args1, &ret), ret
	case "status":
		return client.Call("GoSupervisor.Status", args1, &ret), ret
	case "restart":
		return client.Call("GoSupervisor.Restart", args1, &ret), ret

	}
	return errors.New("command error"), ""
}

func run(client *rpc.Client, cmdline string) (error, string) {
	cmds := strings.Split(cmdline, " ")
	switch len(cmds) {
	case 1:
		return run0(client, cmds[0])

	case 2:
		return run1(client, cmds[0], cmds[1])

	}
	return errors.New("command error"), ""
}
