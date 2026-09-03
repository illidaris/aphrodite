package canal

import (
	"context"
	"reflect"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	pbe "github.com/withlin/canal-go/protocol/entry"
	"google.golang.org/protobuf/proto"
)

func TestConvertCursor(t *testing.T) {
	tests := []struct {
		name string
		raw  any
		tp   int8
		want any
	}{
		{name: "integer", raw: "42", tp: 0, want: 42},
		{name: "string", raw: 42, tp: 1, want: "42"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertCursor(tt.raw, tt.tp); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ConvertCursor() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestSyncConnectorConvert(t *testing.T) {
	change := &pbe.RowChange{EventTypePresent: &pbe.RowChange_EventType{EventType: pbe.EventType_UPDATE}, RowDatas: []*pbe.RowData{{AfterColumns: []*pbe.Column{{Name: "id", Value: "7", IsKey: true}, {Name: "name", Value: "alice"}}}}}
	data, err := proto.Marshal(change)
	if err != nil {
		t.Fatal(err)
	}
	connector := &SyncConnector{SyncConnectorOption: SyncConnectorOption{BaseConnectorOption: BaseConnectorOption{Index: "default", TableMap: map[string]string{"users": "people"}, Log: DefaultLogger{}}}}
	entries, err := connector.Convert(context.Background(), pbe.Entry{EntryTypePresent: &pbe.Entry_EntryType{EntryType: pbe.EntryType_ROWDATA}, StoreValue: data, Header: &pbe.Header{TableName: "users"}})
	if err != nil {
		t.Fatal(err)
	}
	want := []Entry{{Id: "7", Act: ActionIndex, Index: "people", Doc: map[string]any{"id": "7", "name": "alice"}}}
	if !reflect.DeepEqual(entries, want) {
		t.Fatalf("Convert() = %#v, want %#v", entries, want)
	}
}

func TestQuery2Map(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	mock.ExpectPrepare(regexp.QuoteMeta("SELECT id, name FROM users WHERE id > ?")).ExpectQuery().WithArgs(1).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name"}).AddRow(2, []byte("alice")))
	got, err := Query2Map(context.Background(), db, "SELECT id, name FROM users WHERE id > ?", 1)
	if err != nil {
		t.Fatal(err)
	}
	want := []map[string]any{{"id": int64(2), "name": "alice"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Query2Map() = %#v, want %#v", got, want)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
