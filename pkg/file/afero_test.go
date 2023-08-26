package file

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAferoFileGetPut(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name string
		fs   FileAPI
	}{
		{"MemoryFS", NewMemoryFile()},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new aferoFile instance

			defer func() {
				// Delete the file
				err := tc.fs.Delete(ctx, "testfile.txt")
				assert.NoError(t, err, "Failed to delete the file")
			}()

			// Write data to the file
			data := []byte("Test data")
			err := tc.fs.Put(ctx, "testfile.txt", data)
			assert.NoError(t, err, "Failed to write data to the file")

			// Read data from the file
			readData, err := tc.fs.Get(ctx, "testfile.txt")
			assert.NoError(t, err, "Failed to read data from the file")

			// Check if the read data is the same as the written data
			assert.Equal(t, data, readData, "Read data does not match the written data")
		})
	}
}

func TestAferoFileAppendString(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name string
		fs   FileAPI
	}{
		{"MemoryFS", NewMemoryFile()},
		{"OSFS", NewOSFile()},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new aferoFile instance

			defer func() {
				// Delete the file
				err := tc.fs.Delete(ctx, "append_testfile.txt")
				assert.NoError(t, err, "Failed to delete the file")
			}()

			// Append data to the file
			data := "Test data 1\n"
			err := tc.fs.AppendString(ctx, "append_testfile.txt", data)
			assert.NoError(t, err, "Failed to append data to the file")

			// Append more data to the file
			data2 := "Test data 2\n"
			err = tc.fs.AppendString(ctx, "append_testfile.txt", data2)
			assert.NoError(t, err, "Failed to append data to the file")

			// Read data from the file
			readData, err := tc.fs.Get(ctx, "append_testfile.txt")
			assert.NoError(t, err, "Failed to read data from the file")

			// Check if the read data is the same as the appended data
			expectedData := data + data2
			assert.Equal(t, expectedData, string(readData), "Read data does not match the appended data")

		})
	}
}
