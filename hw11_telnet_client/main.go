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
	out := os.Stdout

	ctx, cancel := context.WithCancel(context.Background())

	ctx2, stop2 := signal.NotifyContext(ctx, os.Interrupt)
	defer stop2()

	ctx3, stop3 := signal.NotifyContext(ctx, syscall.SIGQUIT)
	defer stop3()

	client := NewTelnetClient(os.Args[len(os.Args)-2]+":"+os.Args[len(os.Args)-1], timeout, io.NopCloser(in), out)

	err := client.Connect()
	if err != nil {
		fmt.Printf("could not telnet client connect : %s\n", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go stdinSend(ctx, &wg, in, client)

	go func() {
		for {
			select {
			case <-ctx2.Done():
				fmt.Println("I\nwill be\nback!")
				cancel()
				return
			case <-ctx3.Done():
				fmt.Println("I'm telnet client\nBye-bye")
				cancel()
				return
			}
		}
	}()

	wg.Wait()
}

//func recieveWatch(wg *sync.WaitGroup, in bytes.Buffer, t TelnetClient) {
//	err = t.Receive()
//	if err != nil {
//		fmt.Printf("could not receive message: %s\n", err)
//		wg.Done()
//		return
//	}
//}

func stdinSend(ctx context.Context, wg *sync.WaitGroup, in *bytes.Buffer, t TelnetClient) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		default:
			fmt.Println(scanner.Text(), "one")
			in.WriteString(scanner.Text())
			err := t.Send()
			if err != nil {
				fmt.Printf("could not send message: %s\n", err)
				wg.Done()
				return
			}
			if err := scanner.Err(); err != nil {
				fmt.Printf("scanner Err: %s\n", err)
			}
		}

	}

}
