package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/takoyaki-3/filesync/pkg"
)

func main() {
	conf := pkg.LoadConfig()

	for {
		time.Sleep(time.Second)
		raw := GetHTTP(pkg.APIEndpoint(conf) + "auth?sign=" + pkg.Sign())
		if len(raw) == 0 {
			fmt.Println("error,", time.Now())
			continue
		}
		fmt.Println("ok,", time.Now())
	}
}

func GetHTTP(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}
	}
	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	return byteArray
}
