package file

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/nuggxyz/buildrc/internal/logging"
	"github.com/rs/zerolog"
	"github.com/spf13/afero"
)

func TestTargzAndUntargz(t *testing.T) {
	fs := afero.NewMemMapFs()

	ctx := context.Background()

	ctx = logging.NewVerboseLoggerContextWithLevel(ctx, zerolog.TraceLevel)

	tests := []struct {
		name    string
		path    string
		content string
	}{
		{name: "Case 1", content: "This is a test string 1.", path: "test1.txt"},
		{name: "Case 2", content: "This is a test string 2.", path: "abc/test2.txt"},
		{name: "Case 3", content: "This is a test string 3.", path: "abc/123/test3.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create and write the file
			err := afero.WriteFile(fs, tt.path, []byte(tt.content), os.ModePerm)
			if err != nil {
				t.Fatalf("Error writing file: %v", err)
			}

			// Compress the file using Targz
			tar1, err := Targz(ctx, fs, tt.path)
			if err != nil {
				t.Fatalf("Targz() error = %v", err)
			}

			err = fs.Remove(tt.path)
			if err != nil {
				t.Fatalf("Error removing file: %v", err)
			}

			// Decompress the file using Untargz
			_, err = Untargz(ctx, fs, tar1.Name())
			if err != nil {
				t.Fatalf("Untargz() error = %v", err)
			}

			// Read the decompressed content
			decompressedContent, err := afero.ReadFile(fs, tt.path)
			if err != nil {
				t.Fatalf("Error reading decompressed content: %v", err)
			}

			// Read the decompressed content
			compressedContent, err := afero.ReadFile(fs, tt.path+".tar.gz")
			if err != nil {
				t.Fatalf("Error reading decompressed content: %v", err)
			}

			// err = afero.WriteFile(afero.NewOsFs(), "./test1.tar.gz", compressedContent, os.ModePerm)
			// if err != nil {
			// 	t.Fatalf("Error writing file: %v", err)
			// }

			if len(compressedContent) == 0 {
				t.Fatalf("Compressed content is empty")
			}

			// Compare the content
			if string(decompressedContent) != tt.content {
				t.Errorf("Content mismatch: got %s, want %s", string(decompressedContent), tt.content)
			}
		})
	}
}

func TestTargzAndUntargzWithDirChecks(t *testing.T) {
	fs := afero.NewMemMapFs()
	ctx := context.Background()

	ctx = logging.NewVerboseLoggerContextWithLevel(ctx, zerolog.TraceLevel)

	// Path to a directory to test
	testDir := "testDir"

	// Content for testing
	tests := []struct {
		path    string
		content string
	}{
		{"file12.txt", "This is a test string 1."},
		{"subdir/file2.txt", "This is a test string 2."},
		{"subdir/nested/test3.txt", "This is a test string 3."},
	}

	err := fs.Mkdir(testDir, os.ModeDir)
	if err != nil {
		t.Fatalf("Error creating directory: %v", err)
	}

	// Create and write the files
	for _, tt := range tests {
		// dir, _ := filepath.Split(tt.path)
		// if dir != "" {
		// 	if err := fs.MkdirAll(filepath.Join(testDir, dir), os.ModeDir); err != nil {
		// 		t.Fatalf("Error creating directory: %v", err)
		// 	}
		// }
		err := afero.WriteFile(fs, filepath.Join(testDir, tt.path), []byte(tt.content), os.ModePerm)
		if err != nil {
			t.Fatalf("Error writing file: %v", err)
		}
	}

	// Compress the directory using Targz
	tarPath, err := Targz(ctx, fs, testDir)
	if err != nil {
		t.Fatalf("Targz() error = %v", err)
	}
	defer tarPath.Close()

	compressedContent, err := afero.ReadFile(fs, tarPath.Name())
	if err != nil {
		t.Fatalf("Error reading decompressed content: %v", err)
	}

	// err = afero.WriteFile(afero.NewOsFs(), "./test.tar.gz", compressedContent, os.ModePerm)
	// if err != nil {
	// 	t.Fatalf("Error writing file: %v", err)
	// }

	fmt.Println("len compressedContent", len(compressedContent))

	if len(compressedContent) == 0 {
		t.Fatalf("Compressed content is empty")
	}

	afero.Walk(fs, testDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Fatalf("Error walking directory: %v", err)
		}

		fmt.Println(path)
		return nil
	})

	fmt.Println("------------------")

	// Remove the directory
	err = fs.RemoveAll(testDir)
	if err != nil {
		t.Fatalf("Error removing directory: %v", err)
	}

	// Decompress the directory using Untargz
	t2, err := Untargz(ctx, fs, tarPath.Name())
	if err != nil {
		t.Fatalf("Untargz() error = %v", err)
	}
	defer t2.Close()

	if t2.Name() != testDir {
		t.Fatalf("Untargz() error = %v", err)
	}

	afero.Walk(fs, testDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Fatalf("Error walking directory: %v", err)
		}

		fmt.Println(path)
		return nil
	})

	// Check the content of the decompressed files
	for _, tt := range tests {
		decompressedContent, err := afero.ReadFile(fs, filepath.Join(testDir, testDir, tt.path))
		if err != nil {
			t.Fatalf("Error reading decompressed content: %v", err)
		}
		if string(decompressedContent) != tt.content {
			t.Errorf("Content mismatch: got %s, want %s", string(decompressedContent), tt.content)
		}
	}
}
