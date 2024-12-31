package dubboex

import (
	"bytes"
	"fmt"
	"sync"

	"dubbo.apache.org/dubbo-go/v3"
	pathEx "github.com/illidaris/file/path"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

// 常量定义
const (
	CONFIG_SERVER = "_viper_config_server" // 配置服务器的键名
	TMP_DIR       = "tmp"                  // 临时文件路径
	TMP_PORT_FILE = "tmp/port.txt"         // 临时端口文件路径
)

// 全局变量定义
var (
	core         *AppConfig             // 应用配置实例
	instanceOpts []dubbo.InstanceOption // Dubbo实例选项
	ins          *dubbo.Instance        // Dubbo实例
	vipers       sync.Map               // 用于存储配置的Map
)

// SetOpts 设置Dubbo实例选项
func SetOpts(opts ...dubbo.InstanceOption) {
	instanceOpts = opts
}

// LoadConfig 加载配置文件并初始化应用配置和Dubbo实例
func LoadConfig(configPath string) error {
	// 检查配置文件是否存在
	b, err := pathEx.ExistOrNot(configPath)
	if err != nil {
		return err
	}
	if !b {
		return fmt.Errorf("[config]%s has no find", configPath)
	}

	// 初始化Viper实例并读取配置文件
	v := viper.New()
	v.SetConfigFile(configPath)
	if readErr := v.ReadInConfig(); readErr != nil {
		return readErr
	}

	// 解析配置文件内容到AppConfig结构体
	core = &AppConfig{}
	vipers.Store(CONFIG_SERVER, v)
	err = v.Unmarshal(core)
	if err != nil {
		return err
	}

	// 创建新的Dubbo实例并初始化
	err = NewInstance(instanceOpts...)
	if err != nil {
		return err
	}

	// 初始化Dubbo和Nacos
	return InitFrmDubboNacos()
}

// Get 读取配置
func Get(key string) interface{} {
	return GetById(core.Nacos.Clientconfig.Appname, key)
}

// GetById 读取配置
func GetById(name, key string) interface{} {
	return Load(name).Get(key)
}

// GetDubboInstance 获取Dubbo实例
func GetDubboInstance() *dubbo.Instance {
	return ins
}

// Store 存储配置到Map中
func Store(key, val string) {
	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(bytes.NewReader([]byte(cast.ToString(val)))); err != nil {
		v.SetConfigType("json")
		if err := v.ReadConfig(bytes.NewReader([]byte(cast.ToString(val)))); err != nil {
			return
		}
	}
	vipers.Store(key, v)
}

// Load 从Map中加载配置
func Load(key string) *viper.Viper {
	if v, ok := vipers.Load(key); ok {
		if vip, ok := v.(*viper.Viper); ok {
			return vip
		}
		return nil
	}
	return nil
}
