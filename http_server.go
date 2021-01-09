package main

import (
	"io"
	"os"
	"log"
	"fmt"
	"time"
	"bytes"
	"net/http"
	"crypto/sha256"
	"encoding/json"
	"path/filepath"
	"io/ioutil"
)

type FileInfos struct{
	List []FileInfo `json:"list"`
}

type FileInfo struct{
	FileName string `json:"filename"`
	Path string	`json:"path"`
	Directory string `json:"directory"`
}

func main(){
	fmt.Println("start")

	mux := http.NewServeMux();
	mux . HandleFunc("/auth", auth);
	mux . HandleFunc("/upload", upload);
	mux . HandleFunc("/download", download);
	mux . HandleFunc("/remove", remove);
	mux . HandleFunc("/getlist", getlist);
	mux . HandleFunc("/chagekey",chagekey);

	// http.Serverのオブジェクトを確保
	// &をつけること構造体ではなくポインタを返却
	server := &http.Server{}; // or new (http.Server);
	server.Addr = ":11182";
	server.Handler = mux;
	server.ListenAndServe();
}

var key []byte

func auth(w http.ResponseWriter, r *http.Request) {
	if !Authentication(w,r){
		return
	}
	fmt.Fprintln(w, true);
}

func chagekey(w http.ResponseWriter, r *http.Request) {
	if !Authentication(w,r){
		return
	}
	fmt.Fprintln(w, true);
	key = readFileAsBytes("./key")
}

func upload(w http.ResponseWriter, r *http.Request) {
	if !Authentication(w,r){
		return
	}
	// このハンドラ関数へのアクセスはPOSTメソッドのみ認める
	if  (r.Method != "POST") {
		fmt.Fprintln(w, "Please access by POST.");
		return;
	}
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, r.Body); err != nil {
    // return nil, err
	}

	queryparm := r.URL.Query()
	if v,ok:=queryparm["path"];ok{
		WriteByte(v[0],buf.Bytes())
		fmt.Fprintln(w, "uploaded.");
	}
}

func download(w http.ResponseWriter, r *http.Request) {
	if !Authentication(w,r){
		return
	}
	queryparm := r.URL.Query()

	if v,ok:=queryparm["path"];ok{
		fileName := ""
		rawData := readFileAsBytes(v[0])
	
		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	
		writer := bytes.NewBuffer(rawData)
		io.Copy(w, writer)	
	} else {
		fmt.Fprintln(w, "hello, world.");
	}
}

func getlist(w http.ResponseWriter, r *http.Request) {
	if !Authentication(w,r){
		return
	}
	queryparm := r.URL.Query()

	if v,ok:=queryparm["path"];ok{
		paths,filenames,directories:=dirwalk(v[0])

		filelist := []FileInfo{}

		for k,p:=range paths{
			fi := FileInfo{}
			fi.Path = p
			fi.FileName = filenames[k]
			fi.Directory = directories[k]
			filelist = append(filelist, fi)
		}

		resp := FileInfos{}
		resp.List = filelist

		outputJson, err := json.Marshal(&resp)
		if err != nil{
			log.Fatalln(err)
		}
		fmt.Fprintln(w, string(outputJson));
	}
}

func remove(w http.ResponseWriter, r *http.Request) {
	if !Authentication(w,r){
		return
	}
	queryparm := r.URL.Query()

	if v,ok:=queryparm["path"];ok{
		if err := os.Remove(v[0]); err != nil {
			fmt.Println(err)
		} else {
			fmt.Fprintln(w, "hello, world.");
		}
	}
}

func dirwalk(dir string) ([]string,[]string,[]string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths,filenames,directories []string
	for _, file := range files {
		if file.IsDir() {
			p,f,d := dirwalk(filepath.Join(dir, file.Name()))
			paths = append(paths, p...)
			filenames = append(filenames, f...)
			directories = append(directories, d...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
		filenames = append(filenames, file.Name())
		directories = append(directories, dir)
	}

	return paths,filenames,directories
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

func Authentication(w http.ResponseWriter, r *http.Request)bool{
	queryparm := r.URL.Query()
	if v,ok:=queryparm["sign"];ok{
		sign := v[0]
	
		now := time.Now().UTC()

		times := []string{}
		times = append(times,now.Format("2006/01/02 15:04:05"))
		times = append(times,now.Add(time.Second).Format("2006/01/02 15:04:05"))
		times = append(times,now.Add(-time.Second).Format("2006/01/02 15:04:05"))

		if len(key)==0{
			key = readFileAsBytes("./key")
		}

		for _,v:=range times{
			hash := sha256.Sum256(append([]byte(v),key...))
			if fmt.Sprintf("%x",hash) == sign{
				return true
			}
		}
	}

	fmt.Fprintln(w, "Authentication failure");
	return false
}