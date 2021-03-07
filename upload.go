package main

import (
	"bytes"
  "net/http"
  "io/ioutil"
	"github.com/takoyaki-3/filesync/pkg"
)

func main(){
	Upload("volume/a.txt",pkg.ReadBytes("./upload.go"))
}

func Upload(path string,rawData []byte)[]byte{
	conf := pkg.LoadConfig()
	
	url := pkg.APIEndpoint(conf)+"upload?sign="+pkg.Sign()+"&path="+path
	resp, _ := http.Post(url,"application/json",bytes.NewBuffer(rawData))
  defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	return byteArray
}
