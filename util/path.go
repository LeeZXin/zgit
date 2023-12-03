package util

import (
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
	"zgit/setting"
)

const windowsSharingViolationError syscall.Errno = 32

func IsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil || os.IsExist(err) {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func JoinRepoPath(v ...string) string {
	return filepath.Join(setting.RepoDir(), filepath.Join(v...)+".git")
}

func JoinWikiPath(v ...string) string {
	return filepath.Join(setting.RepoDir(), filepath.Join(v...)+".wiki")
}

// RemoveAll removes the named file or (empty) directory with at most 5 attempts.
func RemoveAll(name string) error {
	var err error
	for i := 0; i < 5; i++ {
		err = os.RemoveAll(name)
		if err == nil {
			break
		}
		unwrapped := err.(*os.PathError).Err
		if unwrapped == syscall.EBUSY || unwrapped == syscall.ENOTEMPTY || unwrapped == syscall.EPERM || unwrapped == syscall.EMFILE || unwrapped == syscall.ENFILE {
			// try again
			<-time.After(100 * time.Millisecond)
			continue
		}
		if unwrapped == windowsSharingViolationError && runtime.GOOS == "windows" {
			// try again
			<-time.After(100 * time.Millisecond)
			continue
		}
		if unwrapped == syscall.ENOENT {
			// it's already gone
			return nil
		}
	}
	return err
}
