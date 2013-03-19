// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package indexer

import (
	"fmt"
	"github.com/andreaskoch/docs/util"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Addresser interface {
	GetAbsolutePath() string
	GetRelativePath(basePath string) string
}

func GetIndex(repositoryPath string) Index {

	// check if the supplied repository path is set
	if strings.Trim(repositoryPath, " ") == "" {
		panic("The repository path cannot be null or empty.")
	}

	// check if the supplied repository path exists
	if !util.FileExists(repositoryPath) {
		panic(fmt.Sprintf("The supplied repository path `%v` does not exist.", repositoryPath))
	}

	// get all repository items in the supplied repository path
	items := findAllItems(repositoryPath)
	index := NewIndex(repositoryPath, items)

	return index
}

func findAllItems(repositoryPath string) []Item {

	items := make([]Item, 0, 100)

	directoryEntries, err := ioutil.ReadDir(repositoryPath)
	if err != nil {
		fmt.Printf("An error occured while indexing the repository path `%v`. Error: %v\n", repositoryPath, err)
		return nil
	}

	// item search
	directoryContainsItem := false
	for _, element := range directoryEntries {

		itemPath := filepath.Join(repositoryPath, element.Name())

		// check if the file a markdown file
		isMarkdown := isMarkdownFile(itemPath)
		if !isMarkdown {
			continue
		}

		// search for child items
		childs := getChildItems(repositoryPath)

		// create item and append to list
		item, err := NewItem(itemPath, childs)
		if err != nil {
			fmt.Printf("Skipping item: %s\n", err)
			continue
		}

		items = append(items, item)

		// item has been found
		directoryContainsItem = true
		break
	}

	// search in sub directories if there is no item in the current folder
	if !directoryContainsItem {
		items = append(items, getChildItems(repositoryPath)...)
	}

	return items
}

func isMarkdownFile(absoluteFilePath string) bool {
	fileExtension := strings.ToLower(filepath.Ext(absoluteFilePath))
	return fileExtension == ".md"
}

func getChildItems(itemPath string) []Item {

	childItems := make([]Item, 0, 5)

	files, _ := ioutil.ReadDir(itemPath)
	for _, element := range files {

		if element.IsDir() {
			path := filepath.Join(itemPath, element.Name())
			childsInPath := findAllItems(path)
			childItems = append(childItems, childsInPath...)
		}

	}

	return childItems
}