/**
*FileName: file
*Create on 2018/11/22 下午10:58
*Create by mok
*/

package file

import (
	"fmt"
	"dataServer/conf"
	"context"
	"dataServer/pkg/myetcd"
	"go.etcd.io/etcd/clientv3"
)

type RealFileInfo struct {
	FileName string `json:"file_name"` //文件名称  文件名称为一个全局唯一id
	FilePath string `json:"file_path"` //文件路径 文件路径为：日期
	Host string `json:"host"`   //保存的主机位置，可能保存在多个主机上
}


//生成etcd保存的key
func(f *RealFileInfo)CreateKey(sha string)(string){
	key := fmt.Sprintf("%s/%s",conf.ETCD_FILEPATH_PREFIX,sha)
	return key
}

//通过sha获取到realinfo
func (f *RealFileInfo)GetRealFileInfo(sha string)error{
	key := fmt.Sprintf("%s/%s",conf.ETCD_FILEPATH_PREFIX,sha)
	getResp,err := myetcd.KV.Get(context.TODO(),key)
	if err != nil{
		return err
	}
	if err = unmarshal(getResp.Kvs[0].Value,f);err != nil{
		return  err
	}
	return nil
}

func(f *RealFileInfo)IsExist(key string)(bool,error){
	getResp,err := myetcd.KV.Get(context.TODO(),key,clientv3.WithCountOnly())
	if err != nil{
		return true,err
	}
	if getResp.Count > 0 {
		return true,nil
	}
	return false,nil
}

