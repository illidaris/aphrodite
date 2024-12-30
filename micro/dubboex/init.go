package dubboex

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/common/config"
	"dubbo.apache.org/dubbo-go/v3/config_center"
)

// NewInstance 初始化Dubbo实例
func NewInstance(opts ...dubbo.InstanceOption) error {
	instance, err := core.NewDubboInstance(opts...)
	if err != nil {
		return err
	}
	ins = instance
	return nil
}

func InitFrmDubboNacos() error {
	d := config.GetEnvInstance().GetDynamicConfiguration()
	listeners := []*ViperListener{
		NewViperListener(core.Nacos.Service.Servicename, core.Nacos.Service.Groupname),
	}
	for _, v := range core.Nacos.Service.Others {
		listeners = append(listeners, NewViperListener(v.Servicename, v.Groupname))
	}
	for _, v := range listeners {
		if content, err := d.GetRule(v.Name, config_center.WithGroup(v.Group)); err == nil {
			v.Parse(v.Name, content)
		}
		d.AddListener(v.Name, v, config_center.WithGroup(v.Group))
	}
	return nil
}
