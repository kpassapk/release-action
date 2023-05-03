package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var errFileNotFound = errors.New("file not found")

// getValuesYaml searches for the values.yaml file in the given directory and its subdirectories.
// which contains repo and stack in its path. It returns the path to the first file which meets
// these criteria, or an error otherwise.
func getValuesYaml(dir string, repo string, stack string) (string, error) {
	var valuesFilePath string

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}

		fmt.Println(d.Name())
		if d.Name() != "values.yaml" {
			return nil
		}

		s := string(filepath.Separator)
		if !strings.Contains(path, s+repo+s) || !strings.Contains(path, s+stack+s) {
			return nil
		}

		valuesFilePath = path
		return nil
	})

	if err != nil {
		return "", err
	}

	if valuesFilePath == "" {
		return "", errFileNotFound
	}

	return valuesFilePath, nil
}

var errTagNotFound = errors.New("tag not found")

// withUpdatedReleaseTag looks for a `tag: ***` line in the input byte array s. If it finds the tag, it then replaces
// all occurrences of that tag in the byte array.
func withUpdatedReleaseTag(s []byte, tag string) ([]byte, error) {
	re := regexp.MustCompile(`tag: "?([.a-zA-Z0-9_-]+)"?`)

	og := re.FindSubmatch(s)
	if len(og) != 2 {
		return nil, errTagNotFound
	}

	res := bytes.ReplaceAll(s, og[1], []byte(tag))
	return res, nil
}

//goaction:required
//goaction:description root directory of services repo containing all the values.yml
//goaction:default .
var path = os.Getenv("PROJECT_ROOT")

//goaction:required
//goaction:description name of repository to deploy, e.g. commerce-integrations-transformers
var repo = os.Getenv("REPOSITORY")

//goaction:required
//goaction:description name of the stack to deploy, e.g. qa, develop, staging
var stack = os.Getenv("STACK")

//goaction:required
//goaction:description new tag to place in the values.yml, e.g. v1.2.3
var tag = os.Getenv("TAG")

func main() {
	valuesYaml, err := getValuesYaml(path, repo, stack)
	if err != nil {
		panic(err)
	}
	inf, err := os.ReadFile(valuesYaml)
	if err != nil {
		panic(err)
	}
	out, err := withUpdatedReleaseTag(inf, tag)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(valuesYaml, out, 0666); err != nil {
		panic(err)
	}
	fmt.Printf("Wrote %s\n", valuesYaml)
}
