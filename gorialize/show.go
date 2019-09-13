package gorialize

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/drosseau/degob"
)

func ShowOne(dirPath string, filename string) error {
	passphrase := os.Getenv("GORIALIZE_PASS")
	var dir *Directory
	if passphrase == "" {
		dir = NewDirectory("", false)
	} else {
		dir = NewEncryptedDirectory("", false, passphrase)
	}
	id, err := strconv.Atoi(filename)
	if err != nil {
		return errors.New("Resource ID parameter must be a number")
	}
	q := dir.newQueryWithID("show", nil, id)
	q.DirPath = dirPath
	q.ThwartIOBasePathEscape()
	q.ExitIfDirNotExist()
	q.BuildResourcePath()
	q.ReadGobFromDisk()
	q.DecryptGobBuffer()
	if q.FatalError != nil {
		if q.FatalError.Error()[:6] == "cipher" {
			fmt.Println("Failed to decrypt with GORIALIZE_PASS environment variable.")
		}
		return q.FatalError
	}
	reader := bytes.NewReader(q.GobBuffer)
	dec := degob.NewDecoder(reader)
	gobs, err := dec.Decode()
	if err != nil {
		fmt.Println("Failed to decode gob. If directory is encrypted set GORIALIZE_PASS environment variable.")
		return err
	}
	for _, g := range gobs {
		err = g.WriteValue(os.Stdout, degob.SingleLine)
		if err != nil {
			return err
		}
	}
	return nil
}

func ShowAll(dirPath string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	passphrase := os.Getenv("GORIALIZE_PASS")
	var dir *Directory
	if passphrase == "" {
		dir = NewDirectory("", false)
	} else {
		dir = NewEncryptedDirectory("", false, passphrase)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		id, err := strconv.Atoi(f.Name())
		if err != nil {
			continue
		}
		q := dir.newQueryWithID("show", nil, id)
		q.DirPath = dirPath
		q.ThwartIOBasePathEscape()
		q.ExitIfDirNotExist()
		q.BuildResourcePath()
		q.ReadGobFromDisk()
		q.DecryptGobBuffer()
		if q.FatalError != nil {
			if q.FatalError.Error()[:6] == "cipher" {
				fmt.Println("Failed to decrypt with GORIALIZE_PASS environment variable.")
			}
			return q.FatalError
		}
		reader := bytes.NewReader(q.GobBuffer)
		dec := degob.NewDecoder(reader)
		gobs, err := dec.Decode()
		if err != nil {
			fmt.Println("Failed to decode gob. If directory is encrypted set GORIALIZE_PASS environment variable.")
			return err
		}
		for _, g := range gobs {
			err = g.WriteValue(os.Stdout, degob.SingleLine)
			if err != nil {
				return err
			}
		}
	}
	return nil
}