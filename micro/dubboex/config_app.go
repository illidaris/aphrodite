package dubboex

import (
	"fmt"
	"os"
	"path"

	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/config_center"
	log "dubbo.apache.org/dubbo-go/v3/logger"
	"dubbo.apache.org/dubbo-go/v3/metrics"
	"dubbo.apache.org/dubbo-go/v3/otel/trace"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"
	"github.com/illidaris/aphrodite/pkg/backup"
	"github.com/illidaris/aphrodite/pkg/netex"
	"github.com/spf13/cast"
)

type AppConfig struct {
	Mode     string      `yaml:"mode"`
	Nacos    NacosConfig `yaml:"nacos"`
	Isremote bool        `yaml:"isremote"`
}

type NacosConfig struct {
	Clientconfig  Clientconfig    `yaml:"clientconfig"`
	Serverconfigs []Serverconfigs `yaml:"serverconfigs"`
	Service       Service         `yaml:"service"`
}

type Clientconfig struct {
	Appname     string `yaml:"appname"`
	Namespaceid string `yaml:"namespaceid"`
	Password    string `yaml:"password"`
	Username    string `yaml:"username"`
}

type Serverconfigs struct {
	Grpcport int    `yaml:"grpcport"`
	Ipaddr   string `yaml:"ipaddr"`
	Port     int    `yaml:"port"`
}

type Service struct {
	Clustername string   `yaml:"clustername"`
	Ephemeral   bool     `yaml:"ephemeral"`
	Groupname   string   `yaml:"groupname"`
	Servicename string   `yaml:"servicename"`
	Weight      int      `yaml:"weight"`
	Others      []Others `yaml:"others"`
	Enable      bool     `yaml:"enable"`
	Healthy     bool     `yaml:"healthy"`
	Ip          string   `yaml:"ip"`
	Port        int      `yaml:"port"`
	MetricPort  int      `yaml:"metricport"`
}

type Others struct {
	Groupname   string `yaml:"groupname"`
	Servicename string `yaml:"servicename"`
}

func (app AppConfig) NewDubboInstance(opts ...dubbo.InstanceOption) (*dubbo.Instance, error) {
	allOpts := app.DefaultInstanceOptions()
	allOpts = append(allOpts, opts...)
	ins, err := dubbo.NewInstance(
		allOpts...,
	)
	if err != nil {
		return nil, err
	}
	return ins, err
}

func (app AppConfig) DefaultInstanceOptions() []dubbo.InstanceOption {
	opts := []dubbo.InstanceOption{
		dubbo.WithName(app.Nacos.Service.Servicename),
		dubbo.WithTracing(
			trace.WithEnabled(), // enable tracing feature
			trace.WithStdoutExporter(),
		),
		dubbo.WithLogger(
			log.WithLevel("info"),
			log.WithZap(),
		),
		dubbo.WithShutdown(),
	}

	if len(app.Nacos.Serverconfigs) == 0 {
		return opts
	}
	nacosSrv := app.Nacos.Serverconfigs[0]
	opts = append(opts,
		dubbo.WithConfigCenter(
			config_center.WithNacos(),
			config_center.WithGroup(app.Nacos.Service.Groupname),
			config_center.WithNamespace(app.Nacos.Clientconfig.Namespaceid),
			config_center.WithDataID(app.Nacos.Service.Servicename),
			config_center.WithUsername(app.Nacos.Clientconfig.Username),
			config_center.WithPassword(app.Nacos.Clientconfig.Password),
			config_center.WithAddress(fmt.Sprintf("%s:%d", nacosSrv.Ipaddr, nacosSrv.Port)),
		),
		dubbo.WithRegistry(
			registry.WithNacos(),
			registry.WithGroup(app.Nacos.Service.Groupname),
			registry.WithNamespace(app.Nacos.Clientconfig.Namespaceid),
			registry.WithUsername(app.Nacos.Clientconfig.Username),
			registry.WithPassword(app.Nacos.Clientconfig.Password),
			registry.WithAddress(fmt.Sprintf("%s:%d", nacosSrv.Ipaddr, nacosSrv.Port)),
		),
	)

	protocolOpts := []protocol.Option{
		protocol.WithTriple(),
	}
	if app.Nacos.Service.Ip != "" { // 指定IP
		protocolOpts = append(protocolOpts, protocol.WithIp(app.Nacos.Service.Ip))
	}
	protocolOpts = append(protocolOpts, protocol.WithPort(app.GetPort()))
	opts = append(opts, dubbo.WithProtocol(protocolOpts...))

	if app.Nacos.Service.MetricPort > 0 {
		opts = append(opts, dubbo.WithMetrics(
			metrics.WithEnabled(),                          // default false
			metrics.WithPrometheus(),                       // set prometheus metric, default prometheus
			metrics.WithPrometheusExporterEnabled(),        // enable prometheus exporter default false
			metrics.WithPort(app.Nacos.Service.MetricPort), // prometheus http exporter listen at 9099,default 9090
			metrics.WithMetadataEnabled(),                  // enable metadata center metrics, default true
			metrics.WithRegistryEnabled(),                  // enable registry metrics, default true
			metrics.WithConfigCenterEnabled(),              // enable config center metrics, default true)
		))
	}
	return opts
}

// GetPort 获取端口 默认从配置文件中获取，其次从缓存文件中获取，最后随机端口
func (app AppConfig) GetPort() int {
	if app.Nacos.Service.Port > 0 {
		return app.Nacos.Service.Port
	}
	rootDir, _ := os.Getwd()
	tmpPortFile := path.Join(rootDir, TMP_PORT_FILE)
	portStr := backup.ReadFrmDisk(tmpPortFile)
	port := cast.ToInt(portStr)
	if port > 0 {
		return port
	}
	p, _ := netex.GetFreePort()
	if p > 0 {
		backup.WriteToDisk(tmpPortFile, cast.ToString(p))
		return p
	}
	return p
}
