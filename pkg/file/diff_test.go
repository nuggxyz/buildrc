package file

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func FuzzDiff(f *testing.F) {

	// Create a context and a filesystem
	ctx := context.Background()

	// Generate random maps for f1 and f2
	f.Fuzz(func(t *testing.T, a, b, c, d string) {

		fs1 := afero.NewMemMapFs()

		for k, v := range map[string]string{
			"x": a,
			"y": b,
		} {
			if err := afero.WriteFile(fs1, "a/"+k, []byte(v), 0644); err != nil {
				t.Errorf("Error writing file: %v", err)
			}
		}

		for k, v := range map[string]string{
			"x": c,
			"y": d,
		} {
			if err := afero.WriteFile(fs1, "b/"+k, []byte(v), 0644); err != nil {
				t.Errorf("Error writing file: %v", err)
			}
		}

		checked := []string{}

		if a != c {
			checked = append(checked, "~ x")
		}

		if b != d {
			checked = append(checked, "~ y")
		}

		// Call the Diff function, skip the fuzz iteration on error
		diffs, err := Diff(ctx, fs1, "a", "b", []string{"**/*"})
		if err != nil {
			t.Errorf("Error diffing: %v", err)
		}

		slices.Sort(checked)
		slices.Sort(diffs)

		if !assert.Equal(t, checked, diffs) {

			t.Fail()
		}

	})

	if f.Failed() {
		fmt.Println("Failed")
	}

}

func TestReadAndCompareFiles(t *testing.T) {
	ctx := context.Background()

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
			same := readAndCompareFiles(ctx, fs, fs2, "file1.txt")
			assert.Equal(t, tt.expectSame, same)
		})
	}
}

func TestSliceDiff(t *testing.T) {
	ctx := context.Background()

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
			diff := sliceDiff(ctx, tt.slice1, tt.slice2)
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
			name:     "Detect missing file ",
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
		{
			name: "multi double glob a",
			f1: map[string]string{
				"md/a/d/d/abc.md": `# abcd`,
				"md/c/d/e/def.md": `# def`,
				"md/ghi.md":       `# ghi`,
				"abc.txt":         "Hello",
			},
			f2: map[string]string{
				"md/a/d/d/abc.md": `# abc`,
				"md/c/d/e/def.md": `# def`,
				"md/ghi.md":       `# ghi`,
				"abc.txt":         "Hello",
			},
			glob:     []string{"**/*"},
			expected: []string{"~ md/a/d/d/abc.md"},
		},
		{
			name: "extra nested file",
			f1: map[string]string{
				"md/a/d/d/abc.md": `# abcd`,
				"md/c/d/e/def.md": `# def`,
			},
			f2: map[string]string{
				"md/c/d/e/def.md": `# def`,
			},
			glob:     []string{"**/*.md"},
			expected: []string{"+ md/a/d/d/abc.md"},
		},
		{
			name: "missing nested file",
			f1: map[string]string{
				"md/c/d/e/def.md": `# def`,
			},
			f2: map[string]string{
				"md/a/d/d/abc.md": `# abcd`,
				"md/c/d/e/def.md": `# def`,
			},
			glob:     []string{"**/*.md"},
			expected: []string{"- md/a/d/d/abc.md"},
		},
		{
			name: "file missing not matching glob",
			f1: map[string]string{
				"md/c/d/e/def.md": `# def`,
			},
			f2: map[string]string{
				"md/a/d/d/abc.txt": `# abcd`,
				"md/c/d/e/def.md":  `# def`,
			},
			glob:     []string{"**/*.md"},
			expected: []string{},
		},
		{
			name: "file missing not matching glob and dir",
			f1: map[string]string{
				"md/c/d/e/def.md": `# def a`,
			},
			f2: map[string]string{
				"md/a/d/d/abc.txt": `# abcd`,
				"md/c/d/e/def.md":  `# def`,
			},
			glob:     []string{"md/c/d/e/*.md"},
			expected: []string{"~ md/c/d/e/def.md"},
		},
		{
			name: "with dir as file no diff",
			f1: map[string]string{
				"md/c/d/e": `DIR`,
			},
			f2: map[string]string{
				"md/c/d/e": `DIR`,
			},
			glob:     []string{"**/*"},
			expected: []string{},
		},
		{
			name: "with dir as file yes diff",
			f1: map[string]string{
				"md/c/d/e": `DIR`,
			},
			f2: map[string]string{
				"md/c/d/g": `DIR`,
			},
			glob: []string{"**/*"},
			expected: []string{
				"+ md/c/d/e",
				"- md/c/d/g",
			},
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
			slices.Sort(tt.expected)
			slices.Sort(diff)
			assert.Equal(t, tt.expected, diff)
		})
	}
}

func TestDiffWithGitIgnore(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name                   string
		f1                     map[string]string
		f2                     map[string]string
		gitignore              []string
		glob                   []string
		expected               []string
		expectedAfterGitIgnore []string
	}{
		{
			name:                   "Detect missing file",
			f1:                     map[string]string{"a.txt": "Hello"},
			f2:                     map[string]string{"a.txt": "Hello", "b.txt": "World"},
			gitignore:              []string{"*.txt"},
			glob:                   []string{"*.txt"},
			expected:               []string{"- b.txt"},
			expectedAfterGitIgnore: []string{},
		},
		{
			name:                   "muli-level glob",
			f1:                     map[string]string{"a.txt": "Hello", "c.tar.gz": "World"},
			f2:                     map[string]string{"a.txt": "Hello", "c.tar.gz": "Tree"},
			gitignore:              []string{"*.tar.gz"},
			glob:                   []string{"*"},
			expected:               []string{"~ c.tar.gz"},
			expectedAfterGitIgnore: []string{},
		},
		{
			name:                   "muli file glob",
			f1:                     map[string]string{"a.txt": "Hello", "c.tar.gz": "World"},
			f2:                     map[string]string{"a.txt": "Hello", "c.tar.gz": "Tree"},
			gitignore:              []string{"*.tar.gz"},
			glob:                   []string{"*.{txt,tar.gz}"},
			expected:               []string{"~ c.tar.gz"},
			expectedAfterGitIgnore: []string{},
		},
		{
			name:                   "muli-level glob",
			f1:                     map[string]string{"a.txt": "Hello", "a/b/c/c.tar.gz": "World"},
			f2:                     map[string]string{"a.txt": "Hello", "a/b/c/c.tar.gz": "Tree"},
			gitignore:              []string{"*.tar.gz"},
			glob:                   []string{"**/*.{txt,tar.gz}"},
			expected:               []string{"~ a/b/c/c.tar.gz"},
			expectedAfterGitIgnore: []string{},
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
			gitignore:              []string{"*.txt"},
			glob:                   []string{"*.txt", "go.mod"},
			expected:               []string{"~ go.mod"},
			expectedAfterGitIgnore: []string{"~ go.mod"},
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
			gitignore:              []string{"*.md"},
			glob:                   []string{"**/*.md"},
			expected:               []string{"+ md/ghi.md"},
			expectedAfterGitIgnore: []string{},
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
			gitignore:              []string{"*.txt"},
			glob:                   []string{"md/**", "*.txt"},
			expected:               []string{"+ abc.txt"},
			expectedAfterGitIgnore: []string{},
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

			// Write gitignore
			if err := afero.WriteFile(fs1, ".gitignore", []byte(strings.Join(tt.gitignore, "\n")), 0644); err != nil {
				t.Fatal(err)
			}

			// Filter out gitignored files
			ignoreFile, err := FilterGitIgnored(ctx, fs1, diff)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedAfterGitIgnore, ignoreFile)
		})
	}
}
