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
	"context"
	"dataServer/pkg/myetcd"
	"dataServer/conf"
	"mime/multipart"
	"io/ioutil"
	"net/http"
	"io"
)


type FileInfo struct {
	*UserFileInfo
	*RealFileInfo
	Content string
	Size   int64
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


//添加文件
func (fileInfo *FileInfo)AddFile()error{
	var userfileInfoval string
	var err error
	ukey := fileInfo.UserFileInfo.CreateKey()
	rkey := fileInfo.RealFileInfo.CreateKey(fileInfo.Sha)
	if userfileInfoval,err = marshal(fileInfo.UserFileInfo);err != nil{
		return err
	}
	//文件信息是否存在于etcd中
	b,err := fileInfo.UserFileInfo.IsExist()
	if err != nil{
		return err
	}else if b == true {
		return fmt.Errorf("文件已经存在于该路径了：%s",ukey)
	}else{
		//userfileinfo不存在
		//realfileinfo是否存在
		b,err =  fileInfo.RealFileInfo.IsExist(rkey)
		if err != nil{
			return err
		}
		if b == false {
			var ch chan int = make(chan int)
			go fileInfo.writeFile(ch)
			<- ch
			//添加realfileinfo到etcd中
			relfilinfoval,err := marshal(fileInfo.RealFileInfo)
			if err != nil{
				return err
			}
			if _,err =myetcd.KV.Put(context.TODO(),rkey,relfilinfoval);err != nil{
				return err
			}
		}
		//realinfo在etcd中，但是userfileinfo不在
		if _,err =myetcd.KV.Put(context.TODO(),ukey,userfileInfoval); err != nil{
			return err
		}
	}
	return nil
}

//创建文件
func (fileInfo *FileInfo)writeFile(ch chan int)(error){
	//建立文件目录 todo:应该把dirName用环境变量传进来
	dirName := fmt.Sprintf("%s/%s","localDB-",time.Now().Format("2006-01-02"))

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
	//添加真实文件信息
	fileInfo.RealFileInfo = &RealFileInfo{
		FileName:uuid,
		FilePath:dirName,
		Host:conf.SERVER_ADDR,
	}
	ch <- 1
	close(ch)
	_,err = f.WriteString(fileInfo.Content)
	if err !=nil{
		return  err
	}
	return err
}

//获取文件的信息和内容(传入参数为etcd中保存userinfo的key，/prefix/user/path/filename)
func GetFileInfo(username,path,filename string)(*FileInfo,error){
	var fileinfo = &FileInfo{UserFileInfo:&UserFileInfo{},RealFileInfo:&RealFileInfo{}}
	if username == "" || path == "" || filename == ""{
		return nil,fmt.Errorf("传入的参数无效")
	}
	//去etcd中获取到userinfo数据
	key := fmt.Sprintf("%s/%s/%s/%s",conf.ETCD_FILEPATH_PREFIX,username,path,filename)
	var err error
	if err = fileinfo.UserFileInfo.GetUserInfo(key);err != nil{
		return nil,err
	}
	//然后通过sha值去获取到文件存放的信息
	if err = fileinfo.RealFileInfo.GetRealFileInfo(fileinfo.Sha);err != nil{
		return nil,err
	}
	fileUrl := fmt.Sprintf("%s/%s",fileinfo.RealFileInfo.FilePath,fileinfo.RealFileInfo.FileName)
	//判断主机和当前主机是否相同
	if fileinfo.Host == conf.SERVER_ADDR{
			//相同就直接去打开该主机路径下的文件
		if fileinfo.Content,err = readContent(fileUrl);err!=nil{
			return nil,err
		}
	}else {
		//文件不在该主机上，就要去发送http请求到另一个客户端，获取文件内容
		client := http.Client{Timeout:5*time.Second}  //todo:写进配置文件里
		url := fmt.Sprintf("%s?fileUrl=%s",fileinfo.Host,fileUrl)
		request,err := http.NewRequest(http.MethodGet,url,nil)
		if err != nil{

			return nil,err
		}
		resp,err := client.Do(request)
		if err != nil{
			return nil,err
		}
		if resp.StatusCode != http.StatusOK{
			return nil,fmt.Errorf("打开文件失败")
		}
		content,err := readContentFromBody(resp.Body)
		if err != nil{
			return nil,err
		}
		fileinfo.Content = content
	}
	return fileinfo,nil
}

func readContent(fileUrl string)(string,error){
	f,err := os.Open(fileUrl)
	defer f.Close()
	if err != nil{
		return "",err
	}
	//读取文件内容
 	data,err :=ioutil.ReadAll(f)
	if err != nil{
		return "",err
	}
	return string(data),nil
}

func readContentFromBody(r io.ReadCloser)(string,error){
	defer r.Close()
	data,err := ioutil.ReadAll(r)
	if err != nil{
		return "",err
	}
	//return strconv.Unquote(string(data))
	return string(data),err
}

//处理别的主机文件请求
func GetFileinByOther(fileUrl string)(string,error){
	f,err := os.Open(fileUrl)
	defer f.Close()
	if err != nil{
		return "",err
	}
	data,err := ioutil.ReadAll(f)
	if err != nil{
		return "",err
	}
	return string(data),nil
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
