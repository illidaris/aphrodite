package convert

import (
	"testing"
)

func TestSetdefaultFieldFilterLevel(t *testing.T) {
	SetdefaultFieldFilterLevel(FieldFilterLevelEncode, "field1", "field2")

	if defaultFieldFilterLevel != FieldFilterLevelEncode {
		t.Errorf("Expected defaultFieldFilterLevel to be %d, but got %d", FieldFilterLevelEncode, defaultFieldFilterLevel)
	}

	_, has1 := innerAllowedFields["field1"]
	_, has2 := innerAllowedFields["field1"]

	if len(innerAllowedFields) < 2 || !has1 || !has2 {
		t.Errorf("%v", innerAllowedFields)
	}
}

func BenchmarkDefFieldFilter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DefFieldFilter("<p>test2222</p>")
	}
}

func TestDefFieldFilter(t *testing.T) {
	SetdefaultFieldFilterLevel(FieldFilterLevelEncode, "field1", "field2")

	result, _ := DefFieldFilter("<p>test</p>")

	if result != "" {
		t.Errorf("Expected DefFieldFilter to return '%s', but got '%s'", "&lt;p&gt;test&lt;/p&gt;", result)
	}
}

func TestFieldFilter(t *testing.T) {
	SetdefaultFieldFilterLevel(FieldFilterLevelEncode, "field1", "field2")

	result, _ := FieldFilter("<p>test</p>", FieldFilterLevelEncode)

	if result != "" {
		t.Errorf("Expected FieldFilter to return '%s' with level %d, but got '%s'", "&lt;p&gt;test&lt;/p&gt;", FieldFilterLevelEncode, result)
	}

	result, _ = FieldFilter("<script>alert('xss')</script>", FieldFilterLevelAssign, "script")

	if result != "" {
		t.Errorf("Expected FieldFilter to return an empty string with level %d and fields %v, but got '%s'", FieldFilterLevelAssign, []string{"script"}, result)
	}
}

func TestIsField(t *testing.T) {
	result := IsField("valid_field")

	if !result {
		t.Errorf("Expected IsField to return true for 'valid_field', but got false")
	}

	result = IsField("invalid_field@")

	if result {
		t.Errorf("Expected IsField to return false for 'invalid_field@', but got true")
	}
}

func TestMatchString(t *testing.T) {
	result, err := MatchString("test", "^[a-zA-Z]+$")

	if err != nil {
		t.Errorf("Expected MatchString to return no error, but got '%v'", err)
	}

	if !result {
		t.Errorf("Expected MatchString to return true for 'test' and '^([a-zA-Z]+)$', but got false")
	}

	result, err = MatchString("test123", "^[a-zA-Z]+$")

	if err != nil {
		t.Errorf("Expected MatchString to return no error, but got '%v'", err)
	}

	if result {
		t.Errorf("Expected MatchString to return false for 'test123' and '^([a-zA-Z]+)$', but got true")
	}
}

func TestAddAllowFields(t *testing.T) {
	AddAllowFields("field1", "field2")

	_, has1 := innerAllowedFields["field1"]
	_, has2 := innerAllowedFields["field1"]

	if len(innerAllowedFields) < 2 || !has1 || !has2 {
		t.Errorf("%v", innerAllowedFields)
	}
}

func TestIsAllowFields(t *testing.T) {
	AddAllowFields("field1", "field2")

	result := IsAllowFields("field1", innerAllowedFields)

	if !result {
		t.Errorf("Expected IsAllowFields to return true for 'field1', but got false")
	}

	result = IsAllowFields("field3", innerAllowedFields)

	if result {
		t.Errorf("Expected IsAllowFields to return false for 'field3', but got true")
	}
}

func TestAddAllowSortField(t *testing.T) {
	AddAllowSortField("field1", "field2")

	_, has1 := innerAllowedFields["field1"]
	_, has2 := innerAllowedFields["field1"]

	if len(innerAllowedFields) < 2 || !has1 || !has2 {
		t.Errorf("%v", innerAllowedFields)
	}
}
