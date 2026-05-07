package os

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFileExists verifies FileExists returns the correct value for existing
// and non-existing paths.
func TestFileExists(t *testing.T) {
	tmp, err := os.CreateTemp("", "fileexists_test")
	require.NoError(t, err)
	tmp.Close()
	defer os.Remove(tmp.Name())

	assert.True(t, FileExists(tmp.Name()))
	assert.False(t, FileExists(tmp.Name()+".nonexistent"))
}

// TestReadFile verifies that ReadFile can read an existing file and returns an
// error for a non-existent one.
func TestReadFile(t *testing.T) {
	content := []byte("hello readfile")
	tmp, err := os.CreateTemp("", "readfile_test")
	require.NoError(t, err)
	_, err = tmp.Write(content)
	require.NoError(t, err)
	tmp.Close()
	defer os.Remove(tmp.Name())

	data, err := ReadFile(tmp.Name())
	require.NoError(t, err)
	assert.Equal(t, content, data)

	_, err = ReadFile(tmp.Name() + ".nonexistent")
	require.Error(t, err)
}

// TestWriteFile verifies that WriteFile creates a file with the given content
// and that the permissions are set correctly.
func TestWriteFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "writefile_test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "testfile")
	content := []byte("write file content")
	require.NoError(t, WriteFile(path, content, 0600))

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, content, data)
}

// TestMustWriteFile verifies that MustWriteFile successfully writes to a valid
// path.
func TestMustWriteFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "mustwritefile_test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "mustwrite")
	content := []byte("must write content")

	// Should not panic on a writable path.
	require.NotPanics(t, func() {
		MustWriteFile(path, content, 0600)
	})

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, content, data)
}

func TestCopyFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	content := []byte("hello world")
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}

	copyfile := fmt.Sprintf("%s.copy", tmpfile.Name())
	if err := CopyFile(tmpfile.Name(), copyfile); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(copyfile); os.IsNotExist(err) {
		t.Fatal("copy should exist")
	}
	data, err := os.ReadFile(copyfile)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, content) {
		t.Fatalf("copy file content differs: expected %v, got %v", content, data)
	}
	os.Remove(copyfile)
}

func TestEnsureDir(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ensure-dir")
	require.NoError(t, err)
	defer os.RemoveAll(tmp)

	// Should be possible to create a new directory.
	err = EnsureDir(filepath.Join(tmp, "dir"), 0755)
	require.NoError(t, err)
	require.DirExists(t, filepath.Join(tmp, "dir"))

	// Should succeed on existing directory.
	err = EnsureDir(filepath.Join(tmp, "dir"), 0755)
	require.NoError(t, err)

	// Should fail on file.
	err = os.WriteFile(filepath.Join(tmp, "file"), []byte{}, 0644)
	require.NoError(t, err)
	err = EnsureDir(filepath.Join(tmp, "file"), 0755)
	require.Error(t, err)

	// Should allow symlink to dir.
	err = os.Symlink(filepath.Join(tmp, "dir"), filepath.Join(tmp, "linkdir"))
	require.NoError(t, err)
	err = EnsureDir(filepath.Join(tmp, "linkdir"), 0755)
	require.NoError(t, err)

	// Should error on symlink to file.
	err = os.Symlink(filepath.Join(tmp, "file"), filepath.Join(tmp, "linkfile"))
	require.NoError(t, err)
	err = EnsureDir(filepath.Join(tmp, "linkfile"), 0755)
	require.Error(t, err)
}

// Ensure that using CopyFile does not truncate the destination file before
// the origin is positively a non-directory and that it is ready for copying.
// See https://github.com/tendermint/tendermint/issues/6427
func TestTrickedTruncation(t *testing.T) {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "pwn_truncate")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpDir)

	originalWALPath := filepath.Join(tmpDir, "wal")
	originalWALContent := []byte("I AM BECOME DEATH, DESTROYER OF ALL WORLDS!")
	if err := os.WriteFile(originalWALPath, originalWALContent, 0755); err != nil {
		t.Fatal(err)
	}

	// 1. Sanity check.
	readWAL, err := os.ReadFile(originalWALPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(readWAL, originalWALContent) {
		t.Fatalf("Cannot proceed as the content does not match\nGot:  %q\nWant: %q", readWAL, originalWALContent)
	}

	// 2. Now cause the truncation of the original file.
	// It is absolutely legal to invoke os.Open on a directory.
	if err := CopyFile(tmpDir, originalWALPath); err == nil {
		t.Fatal("Expected an error")
	}

	// 3. Check the WAL's content
	reReadWAL, err := os.ReadFile(originalWALPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(reReadWAL, originalWALContent) {
		t.Fatalf("Oops, the WAL's content was changed :(\nGot:  %q\nWant: %q", reReadWAL, originalWALContent)
	}
}

