package config

import (
	"log"

	"github.com/spf13/viper"
)

type config struct {
	App struct {
		Name string
		Port string
	}
	Database struct {
		Dsn          string
		MaxIdleConns int
		MaxOpenConns int
	}
}

var AppConfig *config

func InitConfig() {
	//viper 读取配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		//如果读取配置文件失败，记录错误并退出程序
		log.Fatalf("Error reading config file%v", err)
	}

	AppConfig = &config{}

	//unmarshal 将读取到的配置文件内容解码到结构体中
	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct%v", err)
	}

	InitDB()    //初始化数据库连接
	InitRedis() //初始化Redis连接
}
