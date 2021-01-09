package main

import (
	"bytes"
  "net/http"
  "io/ioutil"
	"github.com/takoyaki-3/filesync/pkg"
)

const APIEndpoint = "http://localhost:11182/"

func main(){
	Upload("volume/a.txt",pkg.ReadBytes("./upload.go"))
}

func Upload(path string,rawData []byte)[]byte{
	url := APIEndpoint+"upload?sign="+pkg.Sign()+"&path="+path
	resp, _ := http.Post(url,"application/json",bytes.NewBuffer(rawData))
  defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	return byteArray
}
