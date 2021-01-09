package main

import (
	"os"
  "fmt"
	"log"
	"time"
	"bytes"
  "net/http"
  "io/ioutil"
	"crypto/sha256"
)

const APIEndpoint = "http://localhost:11182/"

func main(){
	Upload("volume/a.txt",readFileAsBytes("./upload.go"))
}

func Upload(path string,rawData []byte)[]byte{
	url := APIEndpoint+"upload?sign="+Sign()+"&path="+path
	resp, _ := http.Post(url,"application/json",bytes.NewBuffer(rawData))
  defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	return byteArray
}

func Sign()string{
	now := time.Now().UTC().Format("2006/01/02 15:04:05")

	key := readFileAsBytes("./key")

	hash := sha256.Sum256(append([]byte(now),key...))
	return fmt.Sprintf("%x",hash)
}

func readFileAsBytes(path string)[]byte{
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal("error occured 'os.Open()'")
		panic(err)
	}

	info,_:=file.Stat()
	buf := make([]byte,info.Size())
	fmt.Println(info.Size())

	file.Read(buf)
	return buf
}