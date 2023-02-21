package io

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func ReadFromDir(dirname string) map[string]string {
	// @title ReadFromDir
	// @param dirname string "input the dirname your want to get files/md5
	// return fileMap map "map with filename and md5"
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}
	fileMap := make(map[string]string)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePath := dirname + "/" + file.Name()
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			continue
		}
		md5Sum := fmt.Sprintf("%x", md5.Sum(data))
		fileMap[file.Name()] = md5Sum
	}

	return fileMap
}

func ReadFileContent(fname string) (data []string) {
	//@title ReadFileContent:
	//@param
	//Return
	file, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return

}
