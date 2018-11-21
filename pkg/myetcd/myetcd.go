/**
*FileName: etcd_use
*Create on 2018/11/21 下午4:19
*Create by mok
 */

package myetcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var Client *clientv3.Client
var KV clientv3.KV

const (
	NewClientErr = "创建client失败"
)

func Init() {
	Client = InitClient()
	KV = initKv()
	go watch()
	fmt.Println(Client)
}

func InitClient() *clientv3.Client {
	var conf = clientv3.Config{
		Endpoints:   []string{"172.16.196.131:2380", "172.16.196.129:2380", "192.168.50.250:2380"},
		DialTimeout: 5 * time.Second,
	}
	var err error
	if Client, err = clientv3.New(conf); err != nil {
		panic(fmt.Errorf("%s:%s", NewClientErr, err.Error()))
	}
	return Client
}

func initKv() clientv3.KV {
	return clientv3.NewKV(Client)
}

func Close() {
	Client.Close()
}

func watch() {
	wc := Client.Watch(context.Background(), "fileproto", clientv3.WithPrefix(), clientv3.WithPrevKV())
	for {
		v := <-wc
		for _, e := range v.Events {
			fmt.Printf("type:%v\n kv:%v  prevKey:%v  ", e.Type, e.Kv, e.PrevKv)
		}
	}
}
