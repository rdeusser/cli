// +build tools

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func ensureModPath() error {
	_, err := ioutil.ReadFile("go.mod")
	if err != nil {
		return err
	}
	return nil
}

func findTestData() ([]string, error) {
	var paths []string

	if err := ensureModPath(); err != nil {
		return paths, err
	}

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == "tools" {
			log.Println("skipping tools directory")
			return filepath.SkipDir
		}

		if filepath.Ext(info.Name()) != ".go" {
			return nil
		}

		file, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		if bytes.Contains(file, []byte("github.com/hexops/autogold")) {
			absPath, err := filepath.Abs(filepath.Dir(path))
			if err != nil {
				return err
			}
			paths = append(paths, absPath)
		}

		return nil
	})
	if err != nil {
		return paths, err
	}

	return paths, nil
}

func updateTestData() error {
	cmd := exec.Command("go", "test", "-v", "-timeout", "2m", "-update", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	var hadError bool

	if err := ensureModPath(); err != nil {
		log.Fatal(err)
	}

	paths, err := findTestData()
	if err != nil {
		log.Fatal(err)
	}

	for _, path := range paths {
		if err := os.Chdir(path); err != nil {
			log.Fatal(err)
		}

		log.Printf("updating testdata for %s", path)

		if err := updateTestData(); err != nil {
			hadError = true
			fmt.Println(err)
		}
	}

	if hadError {
		log.Println("some error(s) occurred in some of the tests")
	} else {
		log.Println("successfully updated testdata!")
	}
}
