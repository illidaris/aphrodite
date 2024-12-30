package convert

import (
	"encoding/json"
	"testing"
)

func TestSnake(t *testing.T) {
	type Demo struct {
		MonsterId   uint32 `json:"monster_id,omitempty"`
		UserId      uint64 `json:"user_id,omitempty"`
		SpaceId     uint64 `json:"space_id,omitempty"`
		NextStartId string `json:"next_start_id,omitempty"`
		State       int32  `json:"state,omitempty"`
		PageId      uint32 `json:"page_id,omitempty"`
		PageNum     uint32 `json:"page_num,omitempty"`
	}
	demo := &Demo{
		MonsterId:   1,
		NextStartId: "xxx",
	}
	bs, err := json.Marshal(JsonCamelCase{Value: demo, IsUpper: false})
	if err != nil {
		println(err.Error)
	}
	println(string(bs))
}
