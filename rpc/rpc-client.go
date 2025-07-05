package rpc;

import(
	Log     "log"
	Fmt     "fmt"
	Time    "time"
	Sync    "sync"
	Context "context"
	Errors  "errors"
	GRPC    "google.golang.org/grpc"
	GInsec  "google.golang.org/grpc/credentials/insecure"
	GConty  "google.golang.org/grpc/connectivity"
//	GZIP    "google.golang.org/grpc/encoding/gzip"
	PxnUtil "github.com/PoiXson/pxnGoCommon/utils"
	PxnSan  "github.com/PoiXson/pxnGoCommon/utils/san"
	PxnNet  "github.com/PoiXson/pxnGoCommon/net"
	PxnServ "github.com/PoiXson/pxnGoCommon/service"
);



type ClientRPC struct {
	mut_state    Sync.Mutex
	service     *PxnServ.Service
	// transport
	remote      string
	use_tls     bool
	grpc_client *GRPC.ClientConn
}



func NewClientRPC(service *PxnServ.Service, remote string) *ClientRPC {
	return &ClientRPC{
		service: service,
		remote:  remote,
	};
}



func (rpc *ClientRPC) Start() error {
	rpc.mut_state.Lock();
	defer rpc.mut_state.Unlock();
	if rpc.grpc_client != nil { return Errors.New("RPC client already started, in ClientRPC->Start()"); }
	if rpc.remote      == ""  { return Errors.New("RPC address is required, in ClientRPC->Start()"   ); }
	protocol, address, port := PxnNet.SplitProtocolAddressPort(rpc.remote);
	if protocol        == ""  { return Errors.New("protocol is required, in ClientRPC->Start()"      ); }
	switch protocol {
	case "unix":
		rpc.use_tls = false;
//TODO
panic("UNFINISHED UNIX RPC CLIENT");
		break;
	case "tcp", "tcp4", "tcp6":
		if rpc.use_tls { Log.Printf("%sTLS Enabled",  LogPrefix);
		} else {         Log.Printf("%sTLS Disabled", LogPrefix); }
		if !PxnSan.IsSafeDomain(address) {
			return Fmt.Errorf("Invalid address: %s, in ClientRPC->Start()", address); }
		if port == 0 { return Fmt.Errorf("Invalid port: %d, in ClientRPC->Start()"); }
		addrport := Fmt.Sprintf("%s:%d", address, port);
		backoff_maxdelay, err := Time.ParseDuration(DefaultBackoffMaxDelay);
		if err != nil { Log.Panicf("Invalid backoff max delay, in ClientRPC->Start()", err); }
		grpc_client, err := GRPC.NewClient(
			addrport,
			GRPC.WithTransportCredentials(GInsec.NewCredentials()),
			GRPC.WithBackoffMaxDelay(backoff_maxdelay),
		);
		if err != nil { return Fmt.Errorf("%s, failed to connect, in ClientRPC->Start()", err); }
		rpc.grpc_client = grpc_client;
		break;
	default: return Fmt.Errorf("Unknown protocol: %s, in ClientRPC->Start()", protocol);
	}
	go rpc.Serve();
	PxnUtil.SleepC();
	return nil;
}

func (rpc *ClientRPC) Serve() {
	rpc.service.WaitGroup.Add(1);
	defer func() {
		rpc.Close();
		rpc.service.WaitGroup.Done();
	}();
//TODO: remove this
//	rpc.Service.AddCloseE(rpc);
	Log.Printf("Connecting RPC.. %s", rpc.remote);
	rpc.grpc_client.WaitForStateChange(Context.Background(), GConty.Connecting);
	state := rpc.grpc_client.GetState();
	switch state {
	case GConty.Idle, GConty.Ready: break;
	default: Log.Panic("Connect state failure %s, in ClientRPC->Serve()", state);
	}
//TODO
//TODO
//TODO
//TODO
//TODO: replace this with a health listener
//https://github.com/grpc/grpc-go/tree/v1.73.0/health
	last_state := state;
	LOOP_STATE:
	for {
		rpc.grpc_client.WaitForStateChange(Context.Background(), last_state);
		state := rpc.grpc_client.GetState();
		switch state {
		case GConty.Idle:       Log.Printf("%sIdle.. %s",       LogPrefix, rpc.remote);
		case GConty.Connecting: Log.Printf("%sConnecting.. %s", LogPrefix, rpc.remote);
		case GConty.Ready:      Log.Printf("%sReady. %s",       LogPrefix, rpc.remote);
		case GConty.TransientFailure:
			Log.Printf("%sReconnecting.. %s", LogPrefix, rpc.remote);
		case GConty.Shutdown: break LOOP_STATE;
		}
		last_state = state;
	}
//TODO
//TODO
//TODO
//TODO
//TODO
}



func (rpc *ClientRPC) Close() {
	rpc.service.WaitGroup.Add(1);
	rpc.mut_state.Lock();
	defer func() {
		rpc.mut_state.Unlock();
		rpc.service.WaitGroup.Done();
	}();
	if rpc.grpc_client != nil {
		if err := rpc.grpc_client.Close(); err != nil {
			Log.Printf("%s, in ClientRPC->Close()", err); }
		rpc.grpc_client = nil;
	}
}

func (rpc *ClientRPC) IsStopping() bool {
	return (rpc.grpc_client.GetState() == GConty.Shutdown);
}



func (rpc *ClientRPC) GetClientGRPC() *GRPC.ClientConn {
	return rpc.grpc_client;
}

func (rpc *ClientRPC) SetClientGRPC(grpc_client *GRPC.ClientConn) *GRPC.ClientConn {
	rpc.mut_state.Lock();
	defer rpc.mut_state.Unlock();
	previous := rpc.grpc_client;
	rpc.grpc_client = grpc_client;
	return previous;
}
