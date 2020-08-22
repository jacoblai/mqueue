package main

import (
	"context"
	"cors"
	"dbEngine"
	"flag"
	"fmt"
	"github.com/jacoblai/httprouter"
	"limit"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	var (
		host = flag.String("l", ":8088", "host")
	)
	flag.Parse()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println(err)
		return
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	db := dbEngine.NewDbEngin(dir)

	router := httprouter.New()
	router.POST("/api/db", db.EnDb)
	router.GET("/api/db/:id", db.PeekDb)
	router.DELETE("/api/db/:id", db.DelDb)

	srv := &http.Server{Handler: limit.Limit(cors.CORS(router)), ErrorLog: nil}
	srv.Addr = *host
	go func() {
		if err = srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	log.Println("db server on port", *host)

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	cleanup := make(chan bool)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range signalChan {
			ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
			go func() {
				_ = srv.Shutdown(ctx)
				cleanup <- true
			}()
			<-cleanup
			db.Close()
			fmt.Println("safe exit")
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
