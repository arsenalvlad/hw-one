package main

import (
	"flag"
	"fmt"
	"strings"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	if strings.Contains(to, "/root") {
		fmt.Println("error file doesn't write to /root")
		return
	}

	err := Copy(from, to, offset, limit)
	if err != nil {
		fmt.Println(fmt.Errorf("error copy from to file: %w", err))
		return
	}
}
