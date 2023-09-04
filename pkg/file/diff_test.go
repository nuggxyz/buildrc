package file

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestReadAndCompareFiles(t *testing.T) {
	fs := afero.NewMemMapFs()
	fs2 := afero.NewMemMapFs()

	cases := []struct {
		name       string
		content1   string
		content2   string
		expectSame bool
	}{
		{"Same content", "Hello", "Hello", true},
		{"Different content", "Hello", "World", false},
		{"Empty content", "", "", true},
		{"Empty vs non-empty", "", "World", false},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if err := afero.WriteFile(fs, "file1.txt", []byte(tt.content1), 0644); err != nil {
				t.Fatal(err)
			}
			if err := afero.WriteFile(fs2, "file1.txt", []byte(tt.content2), 0644); err != nil {
				t.Fatal(err)
			}
			same := readAndCompareFiles(fs, fs2, "file1.txt")
			assert.Equal(t, tt.expectSame, same)
		})
	}
}

func TestSliceDiff(t *testing.T) {
	cases := []struct {
		name     string
		slice1   []string
		slice2   []string
		expected []string
	}{
		{"Same slices", []string{"a", "b"}, []string{"a", "b"}, []string{}},
		{"Extra in slice2", []string{"a"}, []string{"a", "b"}, []string{"- b"}},
		{"Missing in slice2", []string{"a", "b"}, []string{"a"}, []string{"+ b"}},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			diff := sliceDiff(tt.slice1, tt.slice2)
			assert.Equal(t, tt.expected, diff)
		})
	}
}

func TestDiff(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name     string
		f1       map[string]string
		f2       map[string]string
		glob     []string
		expected []string
	}{
		{
			name:     "Detect missing file",
			f1:       map[string]string{"a.txt": "Hello"},
			f2:       map[string]string{"a.txt": "Hello", "b.txt": "World"},
			glob:     []string{"*.txt"},
			expected: []string{"- b.txt"},
		},
		{
			name:     "muli-level glob",
			f1:       map[string]string{"a.txt": "Hello", "c.tar.gz": "World"},
			f2:       map[string]string{"a.txt": "Hello", "c.tar.gz": "Tree"},
			glob:     []string{"*"},
			expected: []string{"~ c.tar.gz"},
		},
		{
			name:     "muli file glob",
			f1:       map[string]string{"a.txt": "Hello", "c.tar.gz": "World"},
			f2:       map[string]string{"a.txt": "Hello", "c.tar.gz": "Tree"},
			glob:     []string{"*.{txt,tar.gz}"},
			expected: []string{"~ c.tar.gz"},
		},
		{
			name:     "muli-level glob",
			f1:       map[string]string{"a.txt": "Hello", "a/b/c/c.tar.gz": "World"},
			f2:       map[string]string{"a.txt": "Hello", "a/b/c/c.tar.gz": "Tree"},
			glob:     []string{"**/*.{txt,tar.gz}"},
			expected: []string{"~ a/b/c/c.tar.gz"},
		},

		{
			name: "multi file glob",
			f1: map[string]string{
				"a.txt":  "Hello",
				"b.txt":  "World",
				"c.txt":  "Tree",
				"go.mod": `module github.com/walteh/buildrc`,
			},
			f2: map[string]string{
				"a.txt":  "Hello",
				"b.txt":  "World",
				"c.txt":  "Tree",
				"go.mod": `module oops`,
			},
			glob:     []string{"*.txt", "go.mod"},
			expected: []string{"~ go.mod"},
		},
		{
			name: "multi file glob with missing file",
			f1: map[string]string{
				"md/abc.md": `# abc`,
				"md/def.md": `# def`,
				"md/ghi.md": `# ghi`,
			},
			f2: map[string]string{
				"md/abc.md": `# abc`,
				"md/def.md": `# def`,
			},
			glob:     []string{"**/*.md"},
			expected: []string{"+ md/ghi.md"},
		},
		{
			name: "multi double glob",
			f1: map[string]string{
				"md/a/d/d/abc.md": `# abc`,
				"md/c/d/e/def.md": `# def`,
				"md/ghi.md":       `# ghi`,
				"abc.txt":         "Hello",
			},
			f2: map[string]string{
				"md/a/d/d/abc.md": `# abc`,
				"md/c/d/e/def.md": `# def`,
				"md/ghi.md":       `# ghi`,
			},
			glob:     []string{"md/**", "*.txt"},
			expected: []string{"+ abc.txt"},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			fs1 := afero.NewMemMapFs()
			for name, content := range tt.f1 {
				if err := afero.WriteFile(fs1, filepath.Join("a", name), []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
			}
			for name, content := range tt.f2 {
				if err := afero.WriteFile(fs1, filepath.Join("b", name), []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
			}
			diff, err := Diff(ctx, fs1, "a", "b", tt.glob)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, diff)
		})
	}
}
