# TinyPod
A simple http/tcp proxy util which can bring great fun!

### What can i do?
I can map a directory as an http content server with just one single command!

I can also proxy a remote port on my local machine!

You must install golang first!

#### Build
build on linux:
```shell
./make.sh
```
build on windows:
```shell
./make.cmd
```
output executable files:
- bin/http
- bin/proxy

#### How to use

http way:
```shell
http start -p 8080 -c /usr/share/html -a "admin:123456"
```

tcp way:
```shell
proxy start -l 2022 -r 192.168.1.100:22
```

Docker Image:
[http://hub.docker.com/u/hehety/tinypod](http://hub.docker.com/u/hehety/tinypod)
