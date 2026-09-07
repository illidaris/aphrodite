package canal

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMappingInfo(t *testing.T) {
	raw := `{"settings":{"number_of_shards":3,"number_of_replicas":2},"mappings":{"properties":{"title":{"type":"keyword"},"userId":{"type":"long"},"actual":{"type":"long"},"discount":{"type":"long"},"amount":{"type":"long"},"tradeFee":{"type":"long"},"thdFee":{"type":"long"},"cancelExpireBase":{"type":"long"},"cancelExpire":{"type":"long"},"cancelReason":{"type":"integer"},"cancelAt":{"type":"long"},"payingAt":{"type":"long"},"paidAt":{"type":"long"},"closeAt":{"type":"long"},"closeError":{"type":"keyword"},"payWay":{"type":"integer"},"payTer":{"type":"integer"},"ztOrderId":{"type":"keyword"},"ztOrderType":{"type":"integer"},"ztOrderExpire":{"type":"long"},"ztOrderLink":{"type":"keyword"},"thdOrderId":{"type":"keyword"},"ipInt":{"type":"long"},"num":{"type":"long"},"bizId":{"type":"long"},"status":{"type":"integer"},"createBy":{"type":"long"},"createAt":{"type":"long"},"updateBy":{"type":"long"},"updateAt":{"type":"long"},"describe":{"type":"keyword"},"payWayStr":{"type":"keyword"},"orderType":{"type":"integer"},"productId":{"type":"keyword"},"sellerId":{"type":"long"},"releaseExpire":{"type":"long"},"releaseAt":{"type":"long"},"isRelease":{"type":"byte"}}}}`
	m := &MappingInfo{}
	err := json.Unmarshal([]byte(raw), m)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%v", m.ToFieldMap())
}
