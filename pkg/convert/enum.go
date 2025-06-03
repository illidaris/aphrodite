package convert

import "reflect"

type Item struct {
	Code  int64  `json:"code"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

// Map2Item map类型转成结构
func Map2Item(m interface{}) []*Item {
	iter := reflect.ValueOf(m).MapRange()
	items := make([]*Item, 0)
	for iter.Next() {
		item := &Item{
			Code: iter.Key().Int(),
			Name: iter.Value().String(),
		}
		f := iter.Key().MethodByName("String")
		if f.IsValid() {
			alias := f.Call([]reflect.Value{})
			item.Alias = alias[0].String()
		}
		items = append(items, item)
	}
	return items
}
