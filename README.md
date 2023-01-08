# go-ros-fog

基于 Ubuntu18.04 ，ros melodic 版本的ros node

主要使用框架 

Singo （整合了许多开发 web Api 必要的组件）

https://github.com/Gourouting/singo

goroslib

https://github.com/aler9/goroslib

参考框架 

Zinx

https://github.com/aceld/zinx 




## 目的

本项目采用了一系列Golang中比较流行的组件，可以以本项目为基础快速搭建Restful Web API

## 特色

1.  ROS1 官方提供了 C++和 Python 编写的库，但是编译和运行环境超过 1GB，只能使用固定的编译工具。官方由于历史遗留问题，在 ROS1 中存在操作系统和语言倾向性，对 UBuntu 系统和 C++ 支持得最好（ ROSUDP 只有在 roscpp 中实现，rospy并不支持 ）。本项目可以交叉编译后在 Linux、Windows、MacOS 中运行，只需要 MySql 和 Redis 环境，并且编译出的可执行文件小。
2. 

## Godotenv

项目在启动的时候依赖以下环境变量，但是在也可以在项目根目录创建.env文件设置环境变量便于使用(建议开发环境使用)

```shell
MYSQL_DSN="db_user:db_password@/db_name?charset=utf8&parseTime=True&loc=Local" # Mysql连接地址
REDIS_ADDR="127.0.0.1:6379" # Redis端口和地址
REDIS_PW="" # Redis连接密码
REDIS_DB="" # Redis库从0到10
SESSION_SECRET="setOnProducation" # Seesion密钥，必须设置而且不要泄露
GIN_MODE="debug"
```

## Go Mod

本项目使用[Go Mod](https://github.com/golang/go/wiki/Modules)管理依赖。

```shell
go mod init go-ros-fog
export GOPROXY=https://goproxy.cn,direct 或者 go env -w GOPROXY=https://goproxy.cn,direct
go run main.go // 自动安装
```

## 运行

```shell
go run main.go
```

项目运行后启动在3000端口（可以修改，参考gin文档)