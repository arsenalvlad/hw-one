package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type TelnetClienter struct {
	Address string
	Timeout time.Duration
	In      io.ReadCloser
	Out     io.Writer
	Conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TelnetClienter{
		Address: address,
		Timeout: timeout,
		In:      in,
		Out:     out,
		Conn:    nil,
	}
}

func (t *TelnetClienter) Connect() error {
	conn, err := net.Dial("tcp", t.Address)
	if err != nil {
		return fmt.Errorf("can not connect: %w", err)
	}

	t.Conn = conn

	return nil
}

func (t *TelnetClienter) Send() error {
	sms, err := io.ReadAll(t.In)
	if err != nil {
		return fmt.Errorf("can not read all from stdin: %w", err)
	}
	_, err = t.Conn.Write(sms)
	if err != nil {
		return fmt.Errorf("can not write message to host: %w", err)
	}

	return nil
}

func (t *TelnetClienter) Receive() error {
	_, err := io.Copy(t.Out, t.Conn)
	if err != nil {
		return fmt.Errorf("can not receive message to stdout: %w", err)
	}

	return nil
}

func (t *TelnetClienter) Close() error {
	err := t.Conn.Close()
	if err != nil {
		return fmt.Errorf("can not close connect: %w", err)
	}

	return nil
}
