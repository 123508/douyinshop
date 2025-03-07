package config

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type AppConfig struct {
	*MySQLConfig         `mapstructure:"mysql"`
	*RedisConfig         `mapstructure:"redis"`
	*EtcdConfig          `mapstructure:"etcd"`
	*Jwt                 `mapstructure:"jwt"`
	*ElasticSearch       `mapstructure:"elasticsearch"`
	*HertzConfig         `mapstructure:"hertz"`
	*AuthConfig          `mapstructure:"auth"`
	*UserConfig          `mapstructure:"user"`
	*CartConfig          `mapstructure:"cart"`
	*CheckoutConfig      `mapstructure:"checkout"`
	*OrderConfig         `mapstructure:"order"`
	*PaymentConfig       `mapstructure:"payment"`
	*ProductConfig       `mapstructure:"product"`
	*AddressConfig       `mapstructure:"address"`
	*BusinessOrderConfig `mapstructure:"business_order"`
	*ShopConfig          `mapstructure:"shop"`
	*AIConfig            `mapstructure:"ai"`
	*VolcengineConfig    `mapstructure:"volcengine"`
	*AliyunConfig        `mapstructure:"aliyun"`
	*RabbitmqConfig      `mapstructure:"rabbitmq"`
}

type VolcengineConfig struct {
	ApiKey      string `mapstructure:"api_key"`
	DouyinModel string `mapstructure:"douyin_model"`
	Timeout     int    `mapstructure:"timeout"`
	MaxTokens   int    `mapstructure:"max_tokens"`
	RetryTimes  int    `mapstructure:"retry_times"`
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

type ElasticSearch struct {
	Hosts    []string `mapstructure:"hosts"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
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
	Host        string       `mapstructure:"host"`
	Port        int          `mapstructure:"port"`
	ServiceName string       `mapstructure:"service_name"`
	Alipay      AlipayConfig `mapstructure:"alipay"`
}

type AlipayConfig struct {
	AppId      string `mapstructure:"app_id"`
	PrivateKey string `mapstructure:"app_private_key"`
}

type ProductConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	ServiceName string `mapstructure:"service_name"`
}

type AddressConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	ServiceName string `mapstructure:"service_name"`
}

type BusinessOrderConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	ServiceName string `mapstructure:"service_name"`
}

type ShopConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	ServiceName string `mapstructure:"service_name"`
}

type AIConfig struct {
	Host        string        `mapstructure:"host"`
	Port        int           `mapstructure:"port"`
	ServiceName string        `mapstructure:"service_name"`
	Timeout     time.Duration `mapstructure:"timeout"`     // 服务超时时间
	RetryTimes  int           `mapstructure:"retry_times"` // 重试次数

	// AI模型相关配置
	Model struct {
		Name        string  `mapstructure:"name"`        // 模型名称
		MaxTokens   int     `mapstructure:"max_tokens"`  // 最大token数
		Temperature float32 `mapstructure:"temperature"` // 温度参数
	} `mapstructure:"model"`

	// 订单相关配置
	Order struct {
		MaxAmount      float32 `mapstructure:"max_amount"`      // 单笔订单最大金额
		MaxItems       int     `mapstructure:"max_items"`       // 单笔订单最大商品数
		DefaultPayment int32   `mapstructure:"default_payment"` // 默认支付方式
	} `mapstructure:"order"`
}

type AliyunConfig struct {
	UploadPath int          `mapstructure:"uploadPath"`
	Oss        AliOssConfig `mapstructure:"oss"`
}

type AliOssConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyId     string `mapstructure:"accessKeyId"`
	AccessKeySecret string `mapstructure:"accessKeySecret"`
	BucketName      string `mapstructure:"bucketName"`
}

type RabbitmqConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	VirtualHost string `mapstructure:"virtual_host"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	Prefetch    int    `mapstructure:"prefetch"`
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
