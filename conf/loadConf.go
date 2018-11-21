/**
*FileName: pkg
*Create on 2018/11/4 上午9:00
*Create by mok
 */

package conf

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

var (
	NSQ_TCP_ADDRS         []string //nsqd集群结点地址
	NSQ_LOOKUPD_ADDR      string
	DATA_SERVER_TOPIC  string
	HEART_BEAT_TIME    int
	SERVER_ADDR           string
	SERVER_READ_TIME_OUT  int
	SERVER_WRITE_TIME_OUT int
	ETCD_FILEPATH_PREFIX    string
	
)

var (
	LogFile   string
	LogLevel  string
	LogFormat string
)

func Init(conf string) {
	viperSet(conf)
	loadConf()
	//initLog()
	logrus.Info("conf load success")
}

//使用viper加载配置文件
func loadConf() {

	//nsq
	NSQ_TCP_ADDRS = viper.GetStringSlice("nsq.nsqd_tcp_addrs")
	NSQ_LOOKUPD_ADDR = viper.GetString("nsq.nsqlookupd_addr")
	HEART_BEAT_TIME = viper.GetInt("nsq.heart_beat_time")  //心跳间隔时间
	DATA_SERVER_TOPIC = viper.GetString("nsq.data_server_topic")

	//server
	SERVER_ADDR = viper.GetString("server.server_addr")
	SERVER_READ_TIME_OUT = viper.GetInt("server.read_time_out")
	SERVER_WRITE_TIME_OUT = viper.GetInt("server.write_time_out")

	//ETCD
	ETCD_FILEPATH_PREFIX	= viper.GetString("etcd.filepath_prefix")
}

//对viper进行设置
func viperSet(confFile string) error {
	viper.SetConfigType("yml")
	if confFile == "" {
		viper.AddConfigPath("./conf")
		viper.SetConfigName("conf")
	}
	//通过全局变量获取viper的设置
	viper.AutomaticEnv()
	viper.SetEnvPrefix("BLOG1")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	err := viper.ReadInConfig()
	if err != nil {
		err = fmt.Errorf("regist viper is failed:%s", err.Error())
		return err
	}
	wathcing()
	return nil
}

//监听配置文件的改动
func wathcing() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logrus.Infof("config is changed:%s", e.Name)
	})
}
