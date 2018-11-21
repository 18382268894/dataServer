/**
*FileName: handlers
*Create on 2018/11/20 下午2:52
*Create by mok
 */

package file

import (
	"time"
	"fmt"
	"os"
	"dataServer/utils"
	"encoding/json"
	"go.etcd.io/etcd/clientv3"
	"context"
	"dataServer/pkg/myetcd"
	"dataServer/conf"
	"mime/multipart"
	"io/ioutil"
)


type FileInfo struct {
	*UserFileInfo
	*RealFileInfo
	Content string
	Size   int64
}

type RealFileInfo struct {
	FileName string `json:"file_name"` //文件名称  文件名称为一个全局唯一id
	FilePath string `json:"file_path"` //文件路径 文件路径为：日期
	Hosts []string `json:"host"`   //保存的主机位置，可能保存在多个主机上
}

type UserFileInfo struct {
	Username   string    `json:"username"`    //用户名
	FileName   string    `json:"file_name"`   //文件名
	Path       string    `json:"path"`        //文件路径
	CreateTime time.Time `json:"create_time"` //生成时间
	Sha string       	`json:"sha"` 		  //文件内容的散列值
}


//将realFileInfo序列化
func(f *RealFileInfo)Marshal()(string,error){
	return marshal(f)
}

//将realFileInfo序列化
func(f *UserFileInfo)Marshal()(string,error){
	return marshal(f)
}


func(f *RealFileInfo)GetKey(sha string)(string){
	key := fmt.Sprintf("%s/%s",conf.ETCD_FILEPATH_PREFIX,sha)
	return key
}

func(f *UserFileInfo)GetKey()(string){
	key := fmt.Sprintf("%s/%s/%s/%s",conf.ETCD_FILEPATH_PREFIX,f.Username,f.Path,f.FileName)
	return key
}


//从fheader从读取到文件内容
func GetFileFromFeader(username,path string,header *multipart.FileHeader)( *FileInfo,error){
	var fileinfo  = &FileInfo{
		UserFileInfo:&UserFileInfo{},
		RealFileInfo:&RealFileInfo{},
	}
	var err error
	if header == nil{
		err = fmt.Errorf("文件信息失效")
		return nil,err
	}
	f,err := header.Open()
	defer f.Close()
	if err != nil{
		return nil,err
	}
	data,err := ioutil.ReadAll(f)
	if err != nil{
		return nil,err
	}
	fileinfo.Content = string(data)
	fileinfo.Sha,_ = utils.HashEncode(fileinfo.Content)
	fileinfo.Path = path
	fileinfo.CreateTime =time.Now()
	fileinfo.Username = username
	fileinfo.UserFileInfo.FileName = header.Filename
	return fileinfo,nil
}


//添加文件(很完美)
func (fileInfo *FileInfo)AddFile()error{
	var relfilinfoval ,userfileInfoval string
	var err error
	var resp *clientv3.TxnResponse

	ukey := fileInfo.UserFileInfo.GetKey()
	rkey := fileInfo.RealFileInfo.GetKey(fileInfo.Sha)

	if userfileInfoval,err = fileInfo.UserFileInfo.Marshal();err != nil{
		return err
	}
	if relfilinfoval,err = fileInfo.RealFileInfo.Marshal();err != nil{
		return err
	}

	txn := myetcd.KV.Txn(context.TODO())
	txn.If(clientv3.Compare(clientv3.CreateRevision(ukey),"=",0),
		clientv3.Compare(clientv3.CreateRevision(rkey),"=",0)).
		Then(clientv3.OpPut(ukey,userfileInfoval),clientv3.OpPut(rkey,relfilinfoval))

	if resp,err = txn.Commit();err != nil{
		return err
	}

	if resp.Succeeded{
		fileInfo.writeFile()
		return nil
	}
	txn = myetcd.KV.Txn(context.TODO())
	txn.If(clientv3.Compare(clientv3.CreateRevision(ukey),"=",0)).
		Then(clientv3.OpPut(ukey,userfileInfoval))
	if resp,err = txn.Commit();err != nil{
		return err
	}
	if !resp.Succeeded{
		return fmt.Errorf("文件已经存在于该路径:%s",ukey)
	}
	return nil
}

//创建文件
func (fileInfo *FileInfo)writeFile()(error){
	//建立文件目录
	dirName := fmt.Sprintf("/Users/mok/localDB%s",time.Now().Format("2006-01-02"))  //可以写进配置文件里
	err := os.MkdirAll(dirName,0777)
	if err != nil{
		return err
	}
	//生成全局唯一文件名
	uuid,err := utils.CreateUUID()
	if err != nil{
		return err
	}
	f,err := os.OpenFile(dirName+"/"+uuid,os.O_CREATE | os.O_RDWR,0666)
	if err != nil{
		return err
	}
	_,err = f.WriteString(fileInfo.Content)
	if err !=nil{
		return  err
	}

	//添加真实文件信息
	fileInfo.RealFileInfo = &RealFileInfo{
		FileName:uuid,
		FilePath:dirName,
		Hosts:[]string{conf.SERVER_ADDR},
	}
	return err
}


func marshal(v interface{})(string,error){
	data,err := json.Marshal(v)
	if err != nil{
		return "",err
	}
	return string(data),err
}

func unmarshal(data []byte,v interface{})error{
	err := json.Unmarshal(data,v)
	return err
}
