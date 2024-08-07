package plugin

import "gorm.io/gorm"

type MoreFieldPlugin struct {
}

func (p *MoreFieldPlugin) Name() string {
	return "plugin_more_field"
}
func (p *MoreFieldPlugin) Initialize(db *gorm.DB) error {
	return nil
}
