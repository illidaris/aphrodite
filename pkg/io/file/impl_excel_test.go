package file

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
)

type mockExcel struct {
	*excelize.File
}

func (e *mockExcel) SetSheetCellValue(sheet, cell string, value interface{}) error {
	return e.File.SetCellValue(sheet, cell, value)
}
func (e *mockExcel) Write(w io.Writer) error {
	return e.File.Write(w)
}
func (e *mockExcel) GetSheetRows(sheet string) ([][]string, error) {
	return e.File.GetRows(sheet)
}

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

	f := excelize.NewFile()
	// Call the ExcelExporter function with the sheet name and options
	exporter := ExcelExporter(&mockExcel{File: f}, excelize.CoordinatesToCellName, DEF_SHEET_NAME)

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
	f1 := excelize.NewFile()
	mExcel1 := &mockExcel{File: f1}
	exporter := ExcelExporter(mExcel1, excelize.CoordinatesToCellName, DEF_SHEET_NAME)
	want := &mockWriter{}
	_ = exporter(want, [][]string{{"a", "b"}}, []string{"1", "2"})
	// Call the ExcelImporter function with the sheet name and options
	mExcel2 := &mockExcel{File: f1}
	importer := ExcelImporter(mExcel2, DEF_SHEET_NAME)
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
