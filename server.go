package main

import (
	"io"
	"os"
	"log"
	"fmt"
	"time"
	"bytes"
	"strings"
	"net/http"
	"crypto/sha256"
	"encoding/json"
	"path/filepath"
	"io/ioutil"
	"archive/zip"
	"archive/tar"
	"compress/gzip"
)

const port = ":11182"

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
	mux . HandleFunc("/remover", remover);
	mux . HandleFunc("/getlist", getlist);
	mux . HandleFunc("/chagekey",chagekey);
	mux . HandleFunc("/unzip",unzip);
	mux . HandleFunc("/untargz",untargz);

	// http.Serverのオブジェクトを確保
	// &をつけること構造体ではなくポインタを返却
	server := &http.Server{}; // or new (http.Server);
	server.Addr = port;
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

func unzip(w http.ResponseWriter, r *http.Request) {
	if !Authentication(w,r){
		return
	}
	queryparm := r.URL.Query()
	if src,ok:=queryparm["path"];ok{
		if dist,ok:=queryparm["dist"];ok{
			os.MkdirAll(dist[0],0777)
			fmt.Println(src[0],dist[0])
			Unzip(src[0],dist[0])
			fmt.Fprintln(w, "ok");
			return
		}
	}
	fmt.Fprintln(w, "fail");
}

func untargz(w http.ResponseWriter, r *http.Request){
	if !Authentication(w,r){
		return
	}
	queryparm := r.URL.Query()
	if src,ok:=queryparm["path"];ok{
		file, _ := os.Open(src[0])
		defer file.Close()

		// gzipの展開
		gzipReader, _ := gzip.NewReader(file)
		defer gzipReader.Close()

		// tarの展開
		tarReader := tar.NewReader(gzipReader)

		for {
			tarHeader, err := tarReader.Next()
			if err == io.EOF {
				break
			}

			bufsize := 1024
			fullBuf := []byte{}
			for {
				buf := make([]byte,bufsize)
				n,_ := tarReader.Read(buf)
				fullBuf = append(fullBuf, buf[:n]...)
				if n != bufsize {
					break
				}
			}
			fpath := strings.Replace(tarHeader.Name,"\\","/",-1)
			os.MkdirAll(filepath.Dir(fpath), 0666)
			WriteByte(fpath,fullBuf)
		}
	}
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

func Exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func remove(w http.ResponseWriter, r *http.Request) {
	if !Authentication(w,r){
		return
	}
	queryparm := r.URL.Query()

	if v,ok:=queryparm["path"];ok{
		if Exists(v[0]) {
			if err := os.Remove(v[0]); err != nil {
				fmt.Println(err)
			} else {
				fmt.Fprintln(w, "success.");
			}	
		} else {
			fmt.Fprintln(w, "no such file.");
		}
	}
}

func remover(w http.ResponseWriter, r *http.Request) {
	if !Authentication(w,r){
		return
	}
	queryparm := r.URL.Query()

	if v,ok:=queryparm["path"];ok{
		if Exists(v[0]) {
			paths,_,_ := dirwalk(v[0])
			for _,fpath:=range paths{
				if err := os.Remove(fpath); err != nil {
					fmt.Println(err)
				} else {
				}
			}
			fmt.Fprintln(w, "success.");
		} else {
			fmt.Fprintln(w, "no such file.");
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
	wf, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
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

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
			return err
	}
	defer r.Close()

	for _, f := range r.File {
			rc, err := f.Open()
			if err != nil {
					return err
			}
			defer rc.Close()

			if f.FileInfo().IsDir() {
					path := filepath.Join(dest, f.Name)
					os.MkdirAll(path, f.Mode())
			} else {
					buf := make([]byte, f.UncompressedSize)
					_, err = io.ReadFull(rc, buf)
					if err != nil {
							return err
					}

					path := filepath.Join(dest, f.Name)
					if err = ioutil.WriteFile(path, buf, f.Mode()); err != nil {
							return err
					}
			}
	}

	return nil
}
