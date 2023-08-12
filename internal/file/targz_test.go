package file

import (
	"context"
	"os"
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
