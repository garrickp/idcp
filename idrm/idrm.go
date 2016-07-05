package main

import (
	"flag"
	"fmt"
	"github.com/garrickp/idtools/idlib"
	"os"
)

func IdempotantRemove(c *idlib.Context, filePath string) {

	// Check that filePath exists
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {

		// Exists, let's remove it.
		// Hardlink to a temporary file
		tmpPath := filePath + ".rmtmp"
		if linkErr := os.Link(filePath, tmpPath); linkErr != nil {
			c.SetError(fmt.Sprintf("error creating link from '%s' to '%s'", filePath, tmpPath), linkErr)
		}

		// Unlink original file
		if removeErr := os.Remove(filePath); removeErr != nil {
			c.SetError(fmt.Sprintf("unable to remove file path '%#v'", filePath), removeErr)
		}

		// Unlink the temporary file
		if removeErr := os.Remove(tmpPath); removeErr != nil {
			c.SetError(fmt.Sprintf("unable to remove temporary file path '%#v'", tmpPath), removeErr)
		}

		c.AddComment(fmt.Sprintf("removed file '%s'", filePath))
		c.Changed = true

	}

	return
}

func main() {

	var (
		filePath    string
		usrFilePath string
	)

	flag.StringVar(&usrFilePath, "filepath", "", "path of the file to be removed")
	flag.Parse()

	if usrFilePath == "" {
		usrFilePath = flag.Arg(0)
	}

	if usrFilePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	c := idlib.NewContext("remove")
	c.Begin()

	filePath = idlib.CleanupPath(usrFilePath)

	c.AddValue("filepath", filePath)

	IdempotantRemove(c, filePath)

	c.Finish()
}
