package main

import (
	"flag"
	"os/exec"
	"encoding/json"
	"strings"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"bytes"
	"io/ioutil"
	"io"
	"time"
	"fmt"
	"runtime"
)
var (
	// ShellPath="D:\\go_project\\src\\webhook\\gitpull.sh"
	password="1q2w3e4r5t6y7u"
	Queue =make(chan struct{},100)
	QuitChan=make(chan struct{})
)
const(

)

func userpusheventHandle(w http.ResponseWriter, req *http.Request) {
	if !strings.EqualFold("post",req.Method){
		w.Write([]byte("仅支持post提交"))
		return
	}
	body,err:=ioutil.ReadAll(req.Body)
	if err!=nil{
		log.Println("读取请求流错误，err=",err)
	}
	decoder:=json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	var dest interface{}
	if err:=decoder.Decode(&dest);err!=nil{
		if err!=io.EOF{
			log.Println(err)
		}
	}
	bodyMap:=dest.(map[string]interface{})
	log.Println("请求body",bodyMap["password"])
	if bodyMap["password"] == password{
		Queue <- struct{}{}
		w.Write([]byte("ok"))
		return
	}
	w.Write([]byte("密码错误，无效请求"))
}

func exeshell(path string){
	for {
		select {
		case e:=<- Queue:
			log.Println(e)
			log.Println("执行脚本")
			var out []byte
			var err error
			if runtime.GOOS == "windows"{
				cmd:=exec.Command("sh","-c",path)
				out,err=cmd.CombinedOutput()
			}else{
				cmd:=exec.Command("/bin/bash","-c",path)
				out,err=cmd.CombinedOutput()
			}
			if err!=nil{
				log.Printf("cmd.Run() faild with %s .\n",err)
			}
			log.Println(string(out))
		case <- time.After(2*time.Second):
			// log.Println("循环执行")
		case <-QuitChan:
			log.Println("退出")
			break
		}
	}
}

func main() {
	ShellPath:=flag.String("path", "./gitpull.sh", "shell file path")
	port:=flag.Int("port",10066,"http linsten port.")
	flag.Parse()
	// log.Println(*ShellPath,*port)
	flag.Usage()

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	http.HandleFunc("/gitee/pushevent", userpusheventHandle)
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d",*port), nil)
		if err != nil {
			log.Fatalln("server 启动错误", err)
			panic(err)
		}
	}()
	go exeshell(*ShellPath)
	log.Println("server 启动成功,监听端口：",*port)
	<-sigs
	QuitChan<-struct{}{}
	log.Println("server 退出成功")
}
