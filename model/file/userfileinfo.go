/**
*FileName: file
*Create on 2018/11/22 下午10:58
*Create by mok
*/

package file

import (
	"time"
	"fmt"
	"dataServer/conf"
	"dataServer/pkg/myetcd"
	"context"
	"go.etcd.io/etcd/clientv3"
)


type UserFileInfo struct {
	Username   string    `json:"username"`    //用户名
	FileName   string    `json:"file_name"`   //文件名
	Path       string    `json:"path"`        //文件路径
	CreateTime time.Time `json:"create_time"` //生成时间
	Sha string       	 `json:"sha"` 		  //文件内容的散列值
}


//获取到储存在etcd中的key
func(f *UserFileInfo)CreateKey()(string){
	key := fmt.Sprintf("%s/%s/%s/%s",conf.ETCD_FILEPATH_PREFIX,f.Username,f.Path,f.FileName)
	return key
}

//通过路径返回该路径下所有的文件信息(不返回文件内容)
func GetUserFileInfos(path string)(fs []*UserFileInfo,err error){
	getResp,err := myetcd.KV.Get(context.TODO(),path,clientv3.WithPrefix())
	if err != nil{
		return nil,err
	}
	if len(getResp.Kvs)==0{
		return nil,fmt.Errorf("该文件目录下为空")
	}

	for index,_:= range getResp.Kvs{
		var f = &UserFileInfo{}
		s:=string(getResp.Kvs[index].Value)
		unmarshal([]byte(s),f)  //todo:应该做一下判断，偷懒不做了,默认返回的数据格式是没问题的
		fs = append(fs,f)
		fmt.Println(fs)
	}
	return
}

//获取存放在etcd中的用户文件信息，传入的是保存在etcd中的key
func (f *UserFileInfo)GetUserInfo(key string)error{
	getResp,err := myetcd.KV.Get(context.TODO(),key)
	if err != nil{
		return err
	}
	if err = unmarshal(getResp.Kvs[0].Value,f);err != nil{
		return  err
	}
	return nil
}

//判断该文件信息是够存在etcd中
func(f *UserFileInfo)IsExist()(bool,error){
	key := f.CreateKey()
	getResp,err := myetcd.KV.Get(context.TODO(),key,clientv3.WithCountOnly())
	if err != nil{
		return true,err
	}
	if getResp.Count > 0 {
		return true,nil
	}
	return false,nil
}


