package types

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

//IsDir
//
//dashd in cmsys/file.c
func IsDir(path string) bool {
	theState, err := os.Stat(path)
	if err != nil {
		return false
	}
	return theState.IsDir()
}

func CopyFileToFile(src string, dst string) (err error) {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer destination.Close()

	buf := make([]byte, SYS_BUFFER_SIZE)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}

	return nil
}

func CopyFile(src string, dst string) (err error) {
	if IsDir(dst) {
		statSrc, err := os.Stat(src)
		if err != nil {
			return err
		}

		modeSrc := statSrc.Mode()

		if statSrc.IsDir() {
			return CopyDirToDir(src, dst)
		} else if modeSrc.IsRegular() {
			return CopyFileToDir(src, dst)
		} else {
			return ErrInvalidFile
		}
	} else if IsDir(src) {
		return CopyDirToDir(src, dst)
	} else {
		return CopyFileToFile(src, dst)
	}
}

func Mkdir(path string) error {
	return os.Mkdir(path, DEFAULT_FOLDER_CREATE_PERM)
}

func CopyDirToDir(src string, dst string) (err error) {
	_, err = os.Stat(src)
	if err != nil {
		return err
	}

	_, err = os.Stat(dst)
	if err != nil {
		err = Mkdir(dst)
		if err != nil {
			return err
		}
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		entryName := entry.Name()
		if entryName == "." || entryName == ".." {
			continue
		}

		childSrc := strings.Join([]string{src, entryName}, string(os.PathSeparator))
		childDst := strings.Join([]string{dst, entryName}, string(os.PathSeparator))

		if IsDir(childSrc) {
			Mkdir(childDst)
		}

		CopyFile(childSrc, childDst)
	}

	return nil
}

func CopyFileToDir(src string, dst string) (err error) {
	basename := path.Base(src)
	dstFilename := strings.Join([]string{dst, basename}, string(os.PathSeparator))

	return CopyFileToFile(src, dstFilename)
}

//Rename
//
//Force rename src to dst by recursively-deleting old dst.
func Rename(src string, dst string) (err error) {
	_, err = os.Stat(dst)
	if err == nil {
		os.RemoveAll(dst)
	}

	dirname := path.Dir(dst)
	_, err = os.Stat(dirname)
	if err != nil {
		err = os.MkdirAll(dirname, DEFAULT_FOLDER_CREATE_PERM)
		if err != nil {
			return err
		}
	}

	return os.Rename(src, dst)
}
