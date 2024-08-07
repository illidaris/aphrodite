package file

import (
	"os"
	"path"
	"testing"

	fileex "github.com/illidaris/file/path"
)

// TestToFilePath tests the ToFilePath function for various scenarios.
func TestToFilePath(t *testing.T) {
	goTestDir := path.Join(os.TempDir(), "go_unit_test")
	_ = fileex.MkdirIfNotExist(goTestDir)
	// Define test cases
	tests := []struct {
		name    string
		target  string
		headers [][]string
		rows    [][]string
		wantErr bool
	}{
		{
			name:    "Valid file creation",
			target:  path.Join(goTestDir, "test_export.xlsx"),
			headers: [][]string{{"Header1", "Header2"}},
			rows:    [][]string{{"Row1Data1", "Row1Data2"}, {"Row2Data1", "Row2Data2"}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exr := ExcelExporter(DEF_SHEET_NAME) // Assuming Exporter is a type in your package
			// Call the function under test
			err := exr.ToFilePath(tt.target, tt.headers, tt.rows...)
			if err != nil {
				t.Errorf("ToFilePath() error = %v", err)
			}
			defer os.Remove(tt.target)
			b, err := fileex.ExistOrNot(tt.target)
			if err != nil {
				t.Error(err.Error())
			}
			if !b {
				t.Errorf("%s not exist", tt.target)
			}
			imr := ExcelImporter(DEF_SHEET_NAME)
			rows, err := imr.FromFilePath(tt.target)
			if err != nil {
				t.Error(err.Error())
			}
			datas := append(tt.headers, tt.rows...)
			b = equal(rows, datas)
			if !b {
				t.Errorf("ExcelImporter() got = %v, want %v", rows, datas)
			}
		})
	}
}

// TestToFilePath tests the ToFilePath function for various scenarios.
func TestCSVToFilePath(t *testing.T) {
	goTestDir := path.Join(os.TempDir(), "go_unit_test")
	_ = fileex.MkdirIfNotExist(goTestDir)
	// Define test cases
	tests := []struct {
		name    string
		target  string
		headers [][]string
		rows    [][]string
		wantErr bool
	}{
		{
			name:    "Valid file creation",
			target:  path.Join(goTestDir, "test_export.csv"),
			headers: [][]string{{"Header1", "Header2"}},
			rows:    [][]string{{"Row1Data1", "Row1Data2"}, {"Row2Data1", "Row2Data2"}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exr := CsvExporter() // Assuming Exporter is a type in your package
			// Call the function under test
			err := exr.ToFilePath(tt.target, tt.headers, tt.rows...)
			if err != nil {
				t.Errorf("ToFilePath() error = %v", err)
			}
			defer os.Remove(tt.target)
			b, err := fileex.ExistOrNot(tt.target)
			if err != nil {
				t.Error(err.Error())
			}
			if !b {
				t.Errorf("%s not exist", tt.target)
			}
			imr := CsvImporter(nil)
			rows, err := imr.FromFilePath(tt.target)
			if err != nil {
				t.Error(err.Error())
			}
			datas := append(tt.headers, tt.rows...)
			b = equal(rows, datas)
			if !b {
				t.Errorf("CsvImporter() got = %v, want %v", rows, datas)
			}
		})
	}
}
