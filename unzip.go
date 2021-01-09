package main

import (
	"os"
  "fmt"
	"log"
	"time"
  "net/http"
  "io/ioutil"
	"crypto/sha256"
)

const APIEndpoint = "http://localhost:11182/"

func main() {
	res := GetHTTP(APIEndpoint+"unzip?sign="+Sign()+"&path=volume/TY.zip&dist=volume/TY")
	fmt.Println(res)
}

func GetHTTP(url string)[]byte{
  resp, _ := http.Get(url)
  defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	return byteArray
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

	file.Read(buf)
	return buf
}

func Sign()string{
	now := time.Now().UTC().Format("2006/01/02 15:04:05")

	key := readFileAsBytes("./key")

	hash := sha256.Sum256(append([]byte(now),key...))
	return fmt.Sprintf("%x",hash)
}