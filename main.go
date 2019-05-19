package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func userpusheventHandle(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("hello"))
}
func main() {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGHUP)
	http.HandleFunc("/gitee/pushevent", userpusheventHandle)
	go func() {
		err := http.ListenAndServe(":8090", nil)
		if err != nil {
			log.Fatalln("server 启动错误", err)
			panic(err)
		}
	}()
	fmt.Println("server 启动成功")
	<-sigs
}
