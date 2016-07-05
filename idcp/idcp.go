package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"github.com/garrickp/idtools/idlib"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
)

func IsSameFile(c *idlib.Context, source, dest string) bool {
	var sourceInfo os.FileInfo
	var destInfo os.FileInfo
	var err error

	// Check that dest exists
	destInfo, err = os.Stat(dest)
	if os.IsNotExist(err) {
		return false
	}

	// Check file sizes
	sourceInfo, err = os.Stat(source)
	if err != nil {
		log.Fatal("error running stat on", source, ":", err.Error())
	}

	if sourceInfo.Size() != destInfo.Size() {
		return false
	}

	// Check first byte
	sourceFile, err := os.Open(source)
	if err != nil {
		c.SetError("unable to open source for reading", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Open(dest)
	if err != nil {
		c.SetError("unable to open destination for reading", err)
	}
	defer destFile.Close()

	sourceBuff := make([]byte, 1)
	destBuff := make([]byte, 1)
	if _, err = sourceFile.Read(sourceBuff); err != nil {
		log.Fatal("error reading from sourcefile:", err.Error())
	}

	if _, err = destFile.Read(destBuff); err != nil {
		log.Fatal("error reading from destfile:", err.Error())
	}

	if !reflect.DeepEqual(sourceBuff, destBuff) {
		return false
	}

	// Check last byte
	if _, err = sourceFile.ReadAt(sourceBuff, sourceInfo.Size()-1); err != nil {
		log.Fatal("error reading from sourcefile:", err.Error())
	}
	if _, err = destFile.ReadAt(destBuff, destInfo.Size()-1); err != nil {
		log.Fatal("error reading from destfile:", err.Error())
	}

	if !reflect.DeepEqual(sourceBuff, destBuff) {
		return false
	}

	// Check hash
	sourceFile.Seek(0, 0)
	destFile.Seek(0, 0)

	sourceHashSum := make([]byte, 0)
	destHashSum := make([]byte, 0)
	sourceHash := sha256.New()
	destHash := sha256.New()

	if _, err = io.Copy(sourceHash, sourceFile); err != nil {
		log.Fatal("error obtaining hash from source file:", err.Error())
	}
	if _, err = io.Copy(destHash, destFile); err != nil {
		log.Fatal("error obtaining hash from dest file:", err.Error())
	}

	sourceHashSum = sourceHash.Sum(sourceHashSum)
	destHashSum = destHash.Sum(destHashSum)

	if !reflect.DeepEqual(sourceHashSum, destHashSum) {
		return false
	}

	return true
}

/*
	This operation is NOT idempotant; Go does not offer a friendly way to obtain
	the uid, gid, or mode of a file.
*/
func SetPerms(c *idlib.Context, path, owner, group string, mode os.FileMode) {
	if !idlib.OSSupportChown() {
		return
	}

	// XXX Explore FileInfo.Sys() and see if we can get the Linux file mode out of there.
	uid := os.Getuid()
	gid := os.Getgid()

	if owner != "" {
		uid, gid = idlib.LookupUser(owner)
	}
	if group != "" {
		gid = idlib.LookupGroup(group)
	}

	if err := os.Chown(path, uid, gid); err != nil {
		c.SetError(fmt.Sprintf("unable to change ownership of '%s'", path), err)
	}

	if err := os.Chmod(path, mode); err != nil {
		c.SetError(fmt.Sprintf("unable to change mode of '%s'", path), err)
	}

}

func fileCopy(c *idlib.Context, source, dest string) {
	sourceFile, err := os.Open(source)
	if err != nil {
		c.SetError("unable to open source file for reading", err)
	}
	defer sourceFile.Close()

	destFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		c.SetError("unable to open destination file for writing", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)

	if err != nil {
		c.SetError(fmt.Sprintf("error encountered copying contents of '%s' to '%s'", source, dest), err)
	}

	return
}

func IdempotantCopy(c *idlib.Context, source, dest, owner, group string, mode os.FileMode) {

	if !IsSameFile(c, source, dest) {

		// Write contents of source to tmp file alongside dest
		tmpDest := dest + ".tmp"

		fileCopy(c, source, tmpDest)

		// Set permissions on temporary file so they carry over upon renaming
		SetPerms(c, tmpDest, owner, group, mode)

		// Move tmp file to dest
		err := os.Rename(tmpDest, dest)
		if err != nil {
			c.SetError(fmt.Sprintf("error moving file '%s' to '%s'", tmpDest, dest), err)
		}

		c.AddComment(fmt.Sprintf("copied file from '%s' to '%s'", source, dest))
		c.Changed = true

	} else {

		SetPerms(c, dest, owner, group, mode)

	}

	return
}

func main() {

	var (
		usrSource string
		usrDest   string
		source    string
		dest      string
		usrMode   string
		owner     string
		group     string
		mode      os.FileMode
	)

	flag.StringVar(&usrSource, "source", "", "source of the file to be copied")
	flag.StringVar(&usrDest, "dest", "", "destination that the file should be copied to")
	flag.StringVar(&owner, "owner", "", "owner of the destination file")
	flag.StringVar(&group, "group", "", "group of the destination file")
	flag.StringVar(&usrMode, "mode", "0644", "file mode of the destination file")
	flag.Parse()

	if usrSource == "" {
		usrSource = flag.Arg(0)
	}
	if usrDest == "" {
		usrDest = flag.Arg(1)
	}

	if usrDest == "" || usrSource == "" {
		flag.Usage()
		os.Exit(1)
	}

	c := idlib.NewContext("copy")
	c.Begin()

	source = idlib.CleanupPath(usrSource)
	dest = idlib.CleanupPath(usrDest)

	parsedMode, parseErr := strconv.ParseInt(usrMode, 8, 32)
	if parseErr != nil {
		c.SetError("invalid mode provided", parseErr)
	}
	mode = os.FileMode(parsedMode)

	c.AddValue("source", source)
	c.AddValue("destination", dest)

	IdempotantCopy(c, source, dest, owner, group, mode)

	c.Finish()
}
