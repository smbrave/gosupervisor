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
	for {
		fmt.Printf("gosupervisor>")
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		cmdline := strings.Trim(line, "\r\n\t")
		if cmdline == "" {
			continue
		}
		err, ret := run(client, cmdline)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(ret)
		}
	}
}

func run1(client *rpc.Client, cmd string) (error, string) {
	var ret string
	switch cmd {
	case "reload":
		return client.Call("GoSupervisor.Reload", "", &ret), ret
	case "list":
		return client.Call("GoSupervisor.List", "", &ret), ret
	case "status":
		return client.Call("GoSupervisor.Status", "", &ret), ret
	}
	return errors.New("Unknown command"), ""
}

func run2(client *rpc.Client, cmd string, args1 string) (error, string) {
	var ret string
	switch cmd {
	case "start":
		return client.Call("GoSupervisor.Start", args1, &ret), ret

	case "stop":
		return client.Call("GoSupervisor.Stop", args1, &ret), ret

	case "reload":
		return client.Call("GoSupervisor.Reload", args1, &ret), ret

	case "restart":
		return client.Call("GoSupervisor.Restart", args1, &ret), ret

	}
	return errors.New("Unknown command"), ""
}

func run(client *rpc.Client, cmdline string) (error, string) {
	cmds := strings.Split(cmdline, " ")
	switch len(cmds) {
	case 1:
		return run1(client, cmds[0])

	case 2:
		return run2(client, cmds[0], cmds[1])

	}
	return errors.New("Unknown command"), ""
}
