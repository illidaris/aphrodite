package dubboex

import (
	"dubbo.apache.org/dubbo-go/v3/config_center"
	"github.com/spf13/cast"
)

var _ = config_center.ConfigurationListener(ViperListener{})

func NewViperListener(name string, group string) *ViperListener {
	return &ViperListener{Name: name, Group: group}
}

type ViperListener struct {
	Name  string
	Group string
}

func (l ViperListener) Process(event *config_center.ConfigChangeEvent) {
	l.Parse(event.Key, event.Value)
}

func (l ViperListener) Parse(key string, val interface{}) {
	Store(l.Name, cast.ToString(val))
}
