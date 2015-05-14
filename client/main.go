package main

import (
	"flag"
	"log"
	"zombie"
)

func main() {
	host := flag.String("host", "127.0.0.1", "host IP Address")
	client, err := zombie.NewClient(*host + ":1234")
	if err != nil {
		log.Fatal(err.Error())
	}
	client.Run()
}
