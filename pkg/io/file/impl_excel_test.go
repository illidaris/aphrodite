package file

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Define a mock writer to be used in tests
type mockWriter struct {
	bytes.Buffer
}

func (m *mockWriter) Write(p []byte) (int, error) {
	return m.Buffer.Write(p)
}

func TestExcelExporter(t *testing.T) {
	// Create a mock writer
	mockTar := &mockWriter{}

	// Define headers and rows for the Excel sheet
	headers := [][]string{{"Header1", "Header2"}}
	rows := [][]string{{"Value1", "Value2"}}

	// Call the ExcelExporter function with the sheet name and options
	exporter := ExcelExporter(DEF_SHEET_NAME)

	// Call the returned Exporter function with the mock writer, headers, and rows
	err := exporter(mockTar, headers, rows...)

	// Assert that there were no errors during export
	assert.NoError(t, err)

	// Convert the mock writer's content to a byte slice
	content := mockTar.Bytes()

	// Assert that the content is not empty
	assert.NotEmpty(t, content)
}

func TestExcelImporter(t *testing.T) {
	// Create a mock reader with some test data
	exporter := ExcelExporter(DEF_SHEET_NAME)
	want := &mockWriter{}
	_ = exporter(want, [][]string{{"a", "b"}}, []string{"1", "2"})
	// Call the ExcelImporter function with the sheet name and options
	importer := ExcelImporter(DEF_SHEET_NAME)
	rows, err := importer(want)
	// Assert that there were no errors during import
	assert.NoError(t, err)
	// Assert that the imported rows match the expected data
	expectedRows := [][]string{{"a", "b"}, {"1", "2"}}
	assert.Equal(t, expectedRows, rows)
}

func TestGetSheetOrDefault(t *testing.T) {
	// Test with an empty sheet name
	sheet := ""
	result := GetSheetOrDefault(sheet)
	expected := DEF_SHEET_NAME
	assert.Equal(t, expected, result, "Expected to get the default sheet name")

	// Test with a non-empty sheet name
	sheet = "CustomSheet"
	result = GetSheetOrDefault(sheet)
	expected = "CustomSheet"
	assert.Equal(t, expected, result, "Expected to get the custom sheet name")
}
