package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
	"github.com/takoyaki-3/filesync/pkg"
)

const pkg.APIEndpoint(conf) = "http://localhost:11182/"

func main() {
	conf := pkg.LoadConfig()
	res := GetHTTP(pkg.APIEndpoint(conf)+"unzip?sign="+pkg.Sign()+"&path=volume/TY.zip&dist=volume/TY")
	fmt.Println(res)
}

func GetHTTP(url string)[]byte{
  resp, _ := http.Get(url)
  defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	return byteArray
}
