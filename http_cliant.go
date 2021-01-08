package main

import (
	"os"
  "fmt"
	"log"
	"time"
  "net/http"
  "io/ioutil"
	"crypto/sha256"
	"encoding/json"
)

const APIEndpoint = "http://c3.d.takoyaki3.com:11182/"

type FileInfos struct{
	List []FileInfo `json:"list"`
}

type FileInfo struct{
	FileName string `json:"filename"`
	Path string	`json:"path"`
	Directory string `json:"directory"`
}

func main() {

	// ファイルリストを取得
	infos := GetList("./data")

	// ファイルを順繰りに取得
	for _,v := range infos.List{
		fmt.Println(v.Path)
		os.MkdirAll(v.Directory,0777)
		raw := GetHTTP(APIEndpoint+"download?sign="+Sign()+"&path="+v.Path)
		WriteByte(v.Path,raw)
		GetHTTP(APIEndpoint+"remove?sign="+Sign()+"&path="+v.Path)
	}
}

func GetHTTP(url string)[]byte{
  resp, _ := http.Get(url)
  defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	return byteArray
}

func GetList(path string)FileInfos{
	var infos FileInfos
	raw := GetHTTP(APIEndpoint+"getlist?sign="+Sign()+"&path="+path)
	fmt.Println(APIEndpoint+"getlist?sign="+Sign()+"&path="+path)
	if err := json.Unmarshal(raw, &infos); err != nil {
		log.Fatal(err)
	}
	return infos
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

	file.Read(buf)
	return buf
}

func Sign()string{
	now := time.Now().UTC().Format("2006/01/02 15:04:05")

	key := readFileAsBytes("./key")

	hash := sha256.Sum256(append([]byte(now),key...))
	return fmt.Sprintf("%x",hash)
}