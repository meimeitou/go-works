package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ss(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("timeout")
			return
		default:
			fmt.Println("sleep", ctx.Value("k"))
			time.Sleep(time.Second)
		}
	}
}

func main() {
	ctx := context.WithValue(context.Background(), "k", "test")
	cs, _ := context.WithTimeout(ctx, time.Second*10)
	go ss(cs)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGHUP)
	<-ch
}
