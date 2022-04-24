// 下载相关的优化案例

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func DownLoad1() {

	f, err := os.OpenFile("./download.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal("create file fail err=", err)
	}

	defer func() {
		err = f.Close()
		log.Fatal("file close fail err=", err)
	}()

	resp, err := http.Get("")
	if err != nil {
		log.Fatal("http request fail err=", err)
	}

	defer func() {
		resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("read data fail err=", err)
	}
	b, err := f.Write(body)
	if err != nil {
		log.Fatal("write file fail err=", err)
	}

	fmt.Println("succ! byte=", b)
}

func DownLoad2() {
	f, err := os.OpenFile("./download.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal("create file fail err=", err)
	}

	defer func() {
		err = f.Close()
		log.Fatal("file close fail err=", err)
	}()

	resp, err := http.Get("")
	if err != nil {
		log.Fatal("http request fail err=", err)
	}

	defer func() {
		resp.Body.Close()
	}()

	b, err := io.Copy(f, resp.Body)
	if err != nil {
		log.Fatal("data write fail err=", err)
	}

	fmt.Println("succ! byte=", b)
}

func main() {

}


