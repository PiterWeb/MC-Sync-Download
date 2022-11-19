package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func downloadWorldFromUrl(link string) (File, error) {

	resp, err := http.Get(link)

	if err != nil {
		return File{}, err
	}

	defer resp.Body.Close()

	fileData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return File{}, err
	}

	fileName := resp.Header.Get("Name")

	folderName := strings.Split(fileName, "-")[0]

	file := File{
		Data:       fileData,
		FolderName: folderName,
		Name:       fileName,
	}

	return file, nil

}
