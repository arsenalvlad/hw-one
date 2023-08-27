package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 20, "timeout duration")
}

func main() {
	flag.Parse()

	in := &bytes.Buffer{}

	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel,
		syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer stop()

	client := NewTelnetClient(os.Args[len(os.Args)-2]+":"+os.Args[len(os.Args)-1], timeout, io.NopCloser(in), os.Stdout)
	defer client.Close()

	err := client.Connect()
	if err != nil {
		fmt.Printf("could not telnet client connect : %s\n", err)
		return
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go stdinSend(ctx, &wg, in, client)

	go func() {
		<-ctx.Done()
		sig := <-signalChanel

		switch sig {
		case syscall.SIGTERM, syscall.SIGINT:
			fmt.Println("I'm telnet client\nBye-bye")
		case os.Interrupt:
			fmt.Println("I\nwill be\nback!")
		}
	}()

	wg.Wait()
}

func stdinSend(ctx context.Context, wg *sync.WaitGroup, in *bytes.Buffer, t TelnetClient) {
	scanner := bufio.NewScanner(os.Stdin)

	go func() {
		<-ctx.Done()
		wg.Done()
	}()

	for scanner.Scan() {
		err := t.Receive()
		if err != nil {
			return
		}

		in.WriteString(scanner.Text())
		err = t.Send()
		if err != nil {
			return
		}

		err = t.Receive()
		if err != nil {
			return
		}

		if err = scanner.Err(); err != nil {
			fmt.Printf("scanner Err: %s\n", err)
		}
	}
}
