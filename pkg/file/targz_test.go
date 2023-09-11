package file

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
)

func TestTargzAndUntargz(t *testing.T) {
	fs := afero.NewMemMapFs()

	ctx := context.Background()

	ctx = zerolog.New(zerolog.NewConsoleWriter()).With().Logger().WithContext(ctx)

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

	ctx = zerolog.New(zerolog.NewConsoleWriter()).With().Logger().WithContext(ctx)

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

	if err = afero.Walk(fs, testDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Fatalf("Error walking directory: %v", err)
		}

		fmt.Println(path)
		return nil
	}); err != nil {
		t.Fatalf("Error walking directory: %v", err)
	}

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

	if err = afero.Walk(fs, testDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Fatalf("Error walking directory: %v", err)
		}

		fmt.Println(path)
		return nil
	}); err != nil {
		t.Fatalf("Error walking directory: %v", err)
	}

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

//go:embed testdata/files
var testdata embed.FS

func TestUntarResources(t *testing.T) {
	ctx := context.Background()

	ctx = zerolog.New(zerolog.NewConsoleWriter()).With().Logger().Level(zerolog.TraceLevel).WithContext(ctx)

	src, err := NewEmbedFs(ctx, testdata, "testdata/files")
	if err != nil {
		t.Fatalf("Error creating embed fs: %v", err)
	}

	tests := []struct {
		file     string
		contents []string
		err      error
	}{
		{
			file: "gotestsum_1.10.1_darwin_arm64.tar.gz",
			contents: []string{
				"gotestsum_1.10.1_darwin_arm64/gotestsum",
				"gotestsum_1.10.1_darwin_arm64/LICENSE",
				"gotestsum_1.10.1_darwin_arm64/README.md",
				"gotestsum_1.10.1_darwin_arm64/LICENSE.md",
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {

			dst := afero.NewMemMapFs()

			fle, err := src.Open(tt.file)
			if err != nil {
				t.Fatalf("Error opening file: %v", err)
			}

			err = afero.WriteReader(dst, tt.file, fle)
			if err != nil {
				t.Fatalf("Error writing file: %v", err)
			}

			fle, err = Untargz(ctx, dst, tt.file)
			if err != nil {
				t.Fatalf("Untargz() error = %v", err)
			}

			dirs, err := fle.Readdir(-1)
			if err != nil {
				t.Fatalf("Error reading directory: %v", err)
			}

			if len(dirs) != len(tt.contents) {
				t.Fatalf("Expected %d directory, got %d", len(tt.contents), len(dirs))
			}

			for _, c := range dirs {
				if slices.Contains(tt.contents, c.Name()) {
					t.Fatalf("Expected %s to be in directory", c)
				}
			}
		})
	}
}
