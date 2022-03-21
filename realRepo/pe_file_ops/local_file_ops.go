package pe_file_ops

import (
	"fmt"
	"io/ioutil"
	"os"
)

func GetFileContentAsByte(filePath string) []byte {
	myFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	myByteVal, _ := ioutil.ReadAll(myFile)
	return myByteVal
}

func GetFileContentAsString(filePath string) string {
	myFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	myByteVal, _ := ioutil.ReadAll(myFile)
	return string(myByteVal)
}
