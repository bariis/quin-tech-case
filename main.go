package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bariis/quin-tech-case/server"
	"github.com/bariis/quin-tech-case/task"
)

func main() {

	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt)

	srv := server.NewServer()

	log.Printf("server has started running: %v", srv.Server.Addr)

	go func() {
		if err := srv.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	taskService := task.NewTaskService()

	srv.TaskService = *taskService

	<-done

	log.Println("server is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := srv.Server.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}

	log.Println("server is down")
}
