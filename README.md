# gosupervisor
gosupervisor是采用golang编写的一套进程管理监控工具，对进程状态异常自动恢复并回调告警
### 1.配置config.xml
```bash
$cat /etc/gosupervisor.xml
```

```xml
<gosupervisor>
	<program name="gosupervisor1">
		<command>./test.sh 200</command>
		<directory>/home/vagrant/golang/src/github.com/smbrave/gosupervisor/test</directory>
	</program>
  
	<program name="gosupervisor2">
		<command>./test.sh 201</command>
		<directory>/home/vagrant/golang/src/github.com/smbrave/gosupervisor/test</directory>
        <callbackurl>http://wb.smbrave.cn:8081/supervisor/report</callbackurl>
        <startsec>2</startsec>
	</program>
</gosupervisor>
```
#### 配置说明
* name: 监控程序名，配置中唯一即可
* command:程序运行命令，包括参数
* directory:程序运行路径
* environment:程序运行环境变量，如K1=V1;K2=V2
* stderrfile:标准错误输出文件，默认/dev/null
* stdoutfile:标准输出文件，默认/dev/null
* callbackUrl:进程异常回调地址，用户报警（用户自行开发api报警控制），默认为空
* startsec:程序稳定运行的时间，这个时间过后确定为启动正常，默认10s

### 2.启动服务
```bash
# cd $GOPATH/src/github.com/smbrave/gosupervisor;
# go build
# ./gosupervisor -server=true
```
####启动参数说明
* conf: 读取监控程序的配置文件，默认/etc/gosupervisor.xml
* log: 日志输入路径，默认/var/log/gosupervisor.log
* listen: 服务监听端口和地址，默认127.0.0.1:33870，启动客户端时必须和这个地址相同，否者连接不上
* server: 服务端启动，必须设置为true，默认为false

### 3.启动客户端
```bash
# ./gosupervisor
gosupervisor>list
proc:gosupervisor1   status:running    pid:4648   start:2016-10-18 19:05:37       uptime:44.901528285s       
proc:gosupervisor2   status:running    pid:4657   start:2016-10-18 19:05:37       uptime:44.882143426s       
proc:gosupervisor3   status:running    pid:4659   start:2016-10-18 19:05:37       uptime:44.882107908s       
proc:gosupervisor4   status:running    pid:4649   start:2016-10-18 19:05:37       uptime:44.897342696s

gosupervisor>help
list	: list all program
reload	: reload gosupervisor config
exit	: exit gosupervisor
start	: start program, eg: start procname
stop	: stop program, eg: stop procname
kill	: kill program, eg: kill procname
restart	: restart program, eg: restart procname
gosupervisor>
gosupervisor>
```