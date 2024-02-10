package convert

import (
	"testing"
)

func TestRandomID(t *testing.T) {
	id := RandomID()
	if len(id) != 36 {
		t.Errorf("Expected RandomID to return a UUID string with length 36, but got %d", len(id))
	}

}

func TestSha1ID(t *testing.T) {
	var right string
	for i := 0; i < 100; i++ {
		data := []byte("test data")
		id := Sha1ID(data)
		if i == 0 {
			right = id
		}
		if id != right {
			t.Error(id, "!= ", right)
		}
	}
}

func TestSha1IDSimple(t *testing.T) {
	data := []byte("example data")
	expectedID := "866384d8-e23b-509b-b217-7e44ba487de1"

	// Generate UUID using the given data
	id := Sha1IDSimple(data)

	// Check if the generated ID matches the expected ID
	if id != expectedID {
		t.Errorf("Expected ID: %s, but got: %s", expectedID, id)
	}
}
