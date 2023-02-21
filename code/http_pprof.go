package code

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func InitHttpPprof() {
	arr := make([]string, 10000000)
	arr[1] = "a"

	go func() {
		for {
			time.Sleep(200 * time.Millisecond)
			log.Println("test")
		}
	}()
	http.ListenAndServe(":9012", nil)
}
