# TinyPod
[![Build Status](https://travis-ci.org/hetianyi/tinypod.svg?branch=master)](https://travis-ci.org/hetianyi/tinypod)


A simple http/tcp proxy util which can bring great fun!

### What can I do?
- I can map a directory as an http content server with just one single command!

- I can perform as an API gateway.

- I can also proxy a remote port on local machine!



## Table of Contents

* [Build](#build)
* [Usage](#how-to-use)
* [Docker Image](#docker-image)
* [Donation](#donation)

------



### Build

> You must install golang 1.8+ first!
build on linux:

```shell
./make.sh
```
build on windows:
```shell
./make.cmd
```
output executable file: ```bin/pod```



### Usage

**Full usage:**

```shell
http start \
-p 127.0.0.1:8080 \
-c /app
-w /usr/share/html \
-a "admin:123456" 
-i 
-b "/api:/http://test.host.com/api"
-s ":7878:192.168.1.100:22;127.0.0.1:9000:192.168.1.101:22" 
```

**description:**

```shell
-p 127.0.0.1:8080 
# start http server listening on 127.0.0.1:8080, if you want to disable http server, try: -p -

-c /app
# http context path

-w /usr/share/html
# the http server root directory.

-a "admin:123456" 
# http basic auth

-i
# whether index directory

-b "/api:/http://test.host.com/api"
#  this is typically used as API gateway.

-s ":7878:192.168.1.100:22;127.0.0.1:9000:192.168.1.101:22" 
# this parameter means that pod is listening on 0.0.0.0:7878 
# which is forwarding to 192.168.1.100:22 and 127.0.0.1:9000 is forwarding to 192.168.1.101:22.
```



### Docker Image:
[https://cloud.docker.com/u/hehety/repository/docker/hehety/tinypod](https://cloud.docker.com/u/hehety/repository/docker/hehety/tinypod)



### Donation

#### AliPay

![](doc/alipay.png)

#### Wechat Pay

![](doc/wechatpay.png)

#### Paypal

Donate money by [Paypal](https://www.paypal.me/hehety) to my account **hehety@outlook.com**.