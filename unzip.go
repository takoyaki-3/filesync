package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
	"github.com/takoyaki-3/filesync/pkg"
)

const APIEndpoint = "http://localhost:11182/"

func main() {
	res := GetHTTP(APIEndpoint+"unzip?sign="+pkg.Sign()+"&path=volume/TY.zip&dist=volume/TY")
	fmt.Println(res)
}

func GetHTTP(url string)[]byte{
  resp, _ := http.Get(url)
  defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	return byteArray
}
