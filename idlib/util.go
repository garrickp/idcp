package idlib

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

var osWithChown = []string{"darwin", "linux", "freebsd", "netbsd", "openbsd", "dragonflybsd"}

func OSSupportChown() bool {
	for _, os := range osWithChown {
		if os == runtime.GOOS {
			return true
		}
	}
	return false
}

func CleanupPath(path string) string {
	cleanPath := filepath.Clean(path)
	osPath := filepath.FromSlash(cleanPath)
	absPath, err := filepath.Abs(osPath)
	if err != nil {
		log.Fatal("unable to obtain absolute path to", path, ":", err.Error())
	}
	return absPath
}

func LookupUser(username string) (uid, gid int) {
	if !OSSupportChown() {
		panic("attempting to lookup user on system which does not support chown")
	}

	passwdFile, err := os.Open("/etc/passwd")
	if err != nil {
		return
	}
	defer passwdFile.Close()

	pwFileBuf := bufio.NewReader(passwdFile)

	line, err := pwFileBuf.ReadString('\n')
	for err != io.EOF {
		if len(line) < 1 {
			continue
		}
		lineParts := strings.Split(line, ":")
		if lineParts[0] == username {
			uid, _ = strconv.Atoi(lineParts[2])
			gid, _ = strconv.Atoi(lineParts[3])
			return
		}
		line, err = pwFileBuf.ReadString('\n')
	}
	return
}

func LookupGroup(groupName string) (gid int) {
	if !OSSupportChown() {
		panic("attempting to lookup user on system which does not support chown")
	}

	groupFile, err := os.Open("/etc/group")
	if err != nil {
		return
	}
	defer groupFile.Close()

	groupFileBuf := bufio.NewReader(groupFile)

	line, err := groupFileBuf.ReadString('\n')
	for err != io.EOF {
		if len(line) < 1 {
			continue
		}
		lineParts := strings.Split(line, ":")
		if lineParts[0] == groupName {
			gid, _ = strconv.Atoi(lineParts[2])
			return
		}
		line, err = groupFileBuf.ReadString('\n')
	}
	return
}
