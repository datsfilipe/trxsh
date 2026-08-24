package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMoveFileDirectoryWithTrailingSeparator(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "result")
	dst := filepath.Join(root, "trashed-result")
	if err := os.Mkdir(src, 0755); err != nil {
		t.Fatal(err)
	}

	if err := moveFile(src+string(os.PathSeparator), dst); err != nil {
		t.Fatalf("moveFile() error = %v", err)
	}
	if info, err := os.Stat(dst); err != nil || !info.IsDir() {
		t.Fatalf("destination is not a directory: info=%v err=%v", info, err)
	}
}

func TestMoveFileSymlinkToDirectoryWithTrailingSeparator(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "target")
	src := filepath.Join(root, "result")
	dst := filepath.Join(root, "trashed-result")
	if err := os.Mkdir(target, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, src); err != nil {
		t.Fatal(err)
	}

	if err := moveFile(src+string(os.PathSeparator), dst); err != nil {
		t.Fatalf("moveFile() error = %v", err)
	}
	if info, err := os.Lstat(dst); err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("destination is not a symlink: info=%v err=%v", info, err)
	}
	if _, err := os.Stat(target); err != nil {
		t.Fatalf("symlink target was affected: %v", err)
	}
}
