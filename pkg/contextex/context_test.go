package contextex

import "testing"

func TestContextKeyID(t *testing.T) {
	tests := []struct {
		name string
		key  ContextKey
		id   string
		want ContextKey
	}{
		{"elastic", ElasticID, "abc", "_aphrodite_es_abc"},
		{"db transaction", DbTxID, "tx-1", "_aphrodite_dbtx_tx-1"},
		{"mongo", MongoID, "", "_aphrodite_mongo_"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.key.ID(tt.id); got != tt.want {
				t.Errorf("ID() = %q, want %q", got, tt.want)
			}
		})
	}
}
