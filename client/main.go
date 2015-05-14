package main

import (
	"flag"
	"zombie"
)

func main() {
	host := flag.String("host", "127.0.0.1", "host IP Address")
	client := zombie.NewClient(*host + ":1234")
	client.Run()
}
