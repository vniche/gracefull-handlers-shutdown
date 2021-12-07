package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/vniche/gracefull-handlers-shutdown/handlers"
)

func main() {
	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, os.Interrupt)

	// var waitGroup int = 0
	ctx := context.Background()

	// let's say this is a gRPC server
	handlers.New(ctx, func(ctx context.Context, shutdown chan struct{}, handler *handlers.Handler) {
		defer func() {
			start := time.Now()
			fmt.Printf("stopping gRPC server\n")
			r := rand.Intn(10)
			time.Sleep(time.Duration(r) * time.Second)
			fmt.Printf("stopped gRPC server in %f.2 seconds\n", time.Since(start).Seconds())
			handler.Done()
		}()

		start := time.Now()
		fmt.Printf("starting gRPC server\n")
		r := rand.Intn(10)
		time.Sleep(time.Duration(r) * time.Second)
		fmt.Printf("started gRPC server in %f.2 seconds\n", time.Since(start).Seconds())

		<-shutdown
	})

	// let's say this is a events handler
	handlers.New(ctx, func(ctx context.Context, shutdown chan struct{}, handler *handlers.Handler) {
		defer func() {
			start := time.Now()
			fmt.Printf("stopping events handler\n")
			r := rand.Intn(10)
			time.Sleep(time.Duration(r) * time.Second)
			fmt.Printf("stopped events handler in %f.2 seconds\n", time.Since(start).Seconds())
			handler.Done()
		}()

		start := time.Now()
		fmt.Printf("starting events handler\n")
		r := rand.Intn(10)
		time.Sleep(time.Duration(r) * time.Second)
		fmt.Printf("started events handler in %f.2 seconds\n", time.Since(start).Seconds())

		<-shutdown
	})

	<-interruptSignal
	handlers.GracefullyShutdown()

	fmt.Printf("all stopped")
}
