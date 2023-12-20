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

func JoinRelativeRepoPath(corpId, nodeId, repoName string) string {
	return filepath.Join(corpId, nodeId, repoName) + ".git"
}

func JoinRelativeWikiPath(corpId, nodeId, repoName string) string {
	return filepath.Join(corpId, nodeId, repoName) + ".wiki"
}

func JoinAbsRepoPath(corpId, nodeId, repoName string) string {
	return filepath.Join(setting.RepoDir(), JoinRelativeRepoPath(corpId, nodeId, repoName))
}

func JoinAbsWikiPath(corpId, nodeId, wikiName string) string {
	return filepath.Join(setting.RepoDir(), JoinRelativeWikiPath(corpId, nodeId, wikiName))
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

// Rename renames (moves) oldpath to newpath with at most 5 attempts.
func Rename(oldpath, newpath string) error {
	var err error
	for i := 0; i < 5; i++ {
		err = os.Rename(oldpath, newpath)
		if err == nil {
			break
		}
		unwrapped := err.(*os.LinkError).Err
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

		if i == 0 && os.IsNotExist(err) {
			return err
		}

		if unwrapped == syscall.ENOENT {
			// it's already gone
			return nil
		}
	}
	return err
}
