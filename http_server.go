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
	"github.com/takoyaki-3/goc"
)

type FileList struct{
	Files FileInfo `json:"file_infos"`
}

type FileInfo struct{
	name string

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
	server.Addr = ":11180";
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
	fmt.Fprintln(w, "hello, world.");
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
		paths,filenames:=goc.Dirwalk(v[0])
		fmt.Fprintln(w, paths,filenames);
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