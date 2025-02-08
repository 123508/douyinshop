package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type AppConfig struct {
	*MySQLConfig    `mapstructure:"mysql"`
	*RedisConfig    `mapstructure:"redis"`
	*EtcdConfig     `mapstructure:"etcd"`
	*Jwt            `mapstructure:"jwt"`
	*HertzConfig    `mapstructure:"hertz"`
	*AuthConfig     `mapstructure:"auth"`
	*UserConfig     `mapstructure:"user"`
	*CartConfig     `mapstructure:"cart"`
	*CheckoutConfig `mapstructure:"checkout"`
	*OrderConfig    `mapstructure:"order"`
	*PaymentConfig  `mapstructure:"payment"`
	*ProductConfig  `mapstructure:"product"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Password     string `mapstructure:"password"`
	Port         int    `mapstructure:"port"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

type EtcdConfig struct {
	Endpoints []string `mapstructure:"endpoints"`
	Username  string   `mapstructure:"username"`
	Password  string   `mapstructure:"password"`
}

type Jwt struct {
	AdminSecretKey string `mapstructure:"admin_secret_key"`
	AdminTtl       int    `mapstructure:"admin_ttl"`
	AdminSuv       int    `mapstructure:"admin_suv"`
}

type HertzConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type AuthConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	ServiceName string `mapstructure:"service_name"`
}

type UserConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	ServiceName string `mapstructure:"service_name"`
}

type CartConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	ServiceName string `mapstructure:"service_name"`
}

type CheckoutConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	ServiceName string `mapstructure:"service_name"`
}

type OrderConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	ServiceName string `mapstructure:"service_name"`
}

type PaymentConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	ServiceName string `mapstructure:"service_name"`
}

type ProductConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	ServiceName string `mapstructure:"service_name"`
}

var Conf AppConfig

func init() {
	v := viper.New()
	//v.SetConfigFile("config/config.yaml")

	v.AddConfigPath("config")
	v.AddConfigPath("../../config")
	v.SetConfigName("conf")
	v.SetConfigType("yaml")

	// viper.SetConfigFile("../pkg/config/config.yaml") // 指定配置文件
	err := v.ReadInConfig() // 读取配置信息
	if err != nil {
		fmt.Printf("viper Read Config failed, err:%v\n", err)
		return
	}

	// 把读取到的配置信息反序列化到 Conf 变量中
	if err := v.Unmarshal(&Conf); err != nil {
		fmt.Printf("viper Unmarshal failed, err:%v\n", err)
	}

	v.WatchConfig() // 对配置文件进行监视，若有改变就重新反序列到Conf中
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		}
	})
}
