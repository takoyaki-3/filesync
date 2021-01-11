package pkg

import (
	"os"
  "fmt"
	"log"
	"time"
  // "net/http"
  // "io/ioutil"
	"crypto/sha256"
)

type Config struct {
	Port     int    `json:"port"`
	Hostname string `json:"hostname"`
}

// func 

func ReadBytes(path string)[]byte{
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

func Sign()string{
	now := time.Now().UTC().Format("2006/01/02 15:04:05")

	key := ReadBytes("./key")

	hash := sha256.Sum256(append([]byte(now),key...))
	return fmt.Sprintf("%x",hash)
}