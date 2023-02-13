package io

import (
	"os"
)

func check(e error) {
	//@title check
	//@param
	//Return
	if e != nil {
		panic(e)
	}
}
func ReadFileContent(fname string) (data []byte) {
	//@title ReadFileContent:
	//@param
	//Return
	data, err := os.ReadFile(fname)
	check(err)

	return

}
