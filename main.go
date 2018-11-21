/**
*FileName: dataServer
*Create on 2018/11/22 上午12:34
*Create by mok
*/

package main

import (
	"flag"
	"runtime"
	"github.com/gin-gonic/gin"
	"dataServer/router"
	"net/http"
	"time"
	"dataServer/conf"
	"dataServer/pkg/myetcd"
)

var confPath string

func init(){
	flag.Parse()
	flag.StringVar(&confPath,"c","","config file")
	runtime.GOMAXPROCS(4)
	conf.Init(confPath)  //加载配置
}


func main(){
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	router.Load(r)
	server := http.Server{
		Addr:"192.168.50.250:8081",
		ReadTimeout:time.Duration(conf.SERVER_READ_TIME_OUT)*time.Second,
		WriteTimeout:time.Duration(conf.SERVER_READ_TIME_OUT)*time.Second,
		Handler:r,
	}
	//连接etcd
	myetcd.Init()
	defer myetcd.Close()
	if err := server.ListenAndServe();err != nil{
		panic(err)
	}
}