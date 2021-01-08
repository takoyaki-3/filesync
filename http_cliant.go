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

func main() {
  // url := "http://c3.d.takoyaki3.com:11180/get_file?index=68.9e70bf68a95602c9100347ae67287ebc9f607334108a07123efdf37fc81e8645.0"

	sign := Sign()

	url := "http://localhost:11180/auth?sign="+sign
  resp, _ := http.Get(url)
  defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(byteArray)
	WriteByte("out.mp4",byteArray)
  fmt.Println(string(byteArray)) // htmlをstringで取得
}

func WriteByte(path string,rowData []byte){
	wf, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer wf.Close()

	// データ部分を書き込み
	_, err = wf.Write(rowData)
	if err != nil {
		log.Fatal(err)
	}
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

func Sign()string{
	now := time.Now().UTC().Format("2006/01/02 15:04:05")

	key := readFileAsBytes("./key")

	hash := sha256.Sum256(append([]byte(now),key...))
	return fmt.Sprintf("%x",hash)
}