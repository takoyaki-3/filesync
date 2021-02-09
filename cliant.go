package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/takoyaki-3/filesync/pkg"
)

type FileInfos struct {
	List []FileInfo `json:"list"`
}

type FileInfo struct {
	FileName  string `json:"filename"`
	Path      string `json:"path"`
	Directory string `json:"directory"`
}

func main() {
	conf := pkg.LoadConfig()

	// ファイルリストを取得
	infos := GetList("./data")

	// ファイルを順繰りに取得
	for _, v := range infos.List {
		fmt.Println(v.Path)
		os.MkdirAll(v.Directory, 0777)
		raw := GetHTTP(pkg.APIEndpoint(conf) + "download?sign=" + pkg.Sign() + "&path=" + v.Path)
		pkg.WriteByte(v.Path, raw)
		GetHTTP(pkg.APIEndpoint(conf) + "remove?sign=" + pkg.Sign() + "&path=" + v.Path)
	}
}

func GetHTTP(url string) []byte {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	return byteArray
}

func GetList(path string) FileInfos {
	conf := pkg.LoadConfig()
	var infos FileInfos
	raw := GetHTTP(pkg.APIEndpoint(conf) + "getlist?sign=" + pkg.Sign() + "&path=" + path)
	fmt.Println(pkg.APIEndpoint(conf) + "getlist?sign=" + pkg.Sign() + "&path=" + path)
	if err := json.Unmarshal(raw, &infos); err != nil {
		log.Fatal(err)
	}
	return infos
}
