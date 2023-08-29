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
		err := client.Receive()
		if err != nil {
			fmt.Println(err)
		}
	}()

	wg.Wait()
}

func stdinSend(ctx context.Context, wg *sync.WaitGroup, in *bytes.Buffer, t TelnetClient) {
	go func() {
		<-ctx.Done()
		wg.Done()
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		resp, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		in.WriteString(resp)
		err = t.Send()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
