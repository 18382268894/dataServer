/**
*FileName: utils
*Create on 2018/11/20 下午6:24
*Create by mok
*/

package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/satori/go.uuid"
)

func CreateUUID()(string, error){
	UUID,err := uuid.NewV4()
	if err != nil{
		return "",err
	}
	return UUID.String(),nil
}


/*func NewConsumer(topic string,chanName string,h nsq.Handler)(consumer *nsq.Consumer,err error){
	if consumer,err = nsq.NewConsumer(topic,chanName,nsq.NewConfig());err != nil{
		return nil,err
	}
	consumer.ChangeMaxInFlight(6)
	err = consumer.ConnectToNSQLookupds([]string{"",""})
	if err != nil{
		return nil,err
	}
	consumer.AddHandler(h)
	return consumer,nil
}

func RanDomGetServer(addrs []string)string{
	var l int
	if  l = len(addrs);l == 0{
		return ""
	}
	rand.Seed(time.Now().UnixNano())
	return addrs[rand.Intn(l)]
}

//随机获取一个nsqds地址用于publish
func RadomGetNsqds([]string)string{
	rand.Seed(time.Now().UnixNano())
	return conf.NSQ_TCP_ADDRS[rand.Intn(len(conf.NSQ_TCP_ADDRS))]
}
*/

//哈希加密算法使内容唯一内容
func HashEncode(content string)(string,error){
	hash  := sha256.New()
	_,err := hash.Write([]byte(content))
	if err != nil{
		return "",err
	}
	md := hash.Sum(nil)
	sha := hex.EncodeToString(md)
	return sha,nil
}


