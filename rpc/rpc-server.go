package rpc;

import(
	Log     "log"
	Fmt     "fmt"
	Net     "net"
	Sync    "sync"
	Errors  "errors"
	GRPC    "google.golang.org/grpc"
	_       "google.golang.org/grpc/encoding/gzip"
	PxnUtil "github.com/PoiXson/pxnGoCommon/utils"
	PxnSan  "github.com/PoiXson/pxnGoCommon/utils/san"
	PxnNet  "github.com/PoiXson/pxnGoCommon/net"
	PxnServ "github.com/PoiXson/pxnGoCommon/service"
);



type ServerRPC struct {
	mut_state   Sync.Mutex
	service     *PxnServ.Service
	// transport
	bind        string
	use_tls     bool
	listen      Net.Listener
	grpc_server *GRPC.Server
}



func NewServerRPC(service *PxnServ.Service, bind string) *ServerRPC {
	return &ServerRPC{
		service: service,
		bind:    bind,
	};
}



func (rpc *ServerRPC) Start() error {
	rpc.mut_state.Lock();
	defer rpc.mut_state.Unlock();
	if rpc.bind == "" { rpc.bind = DefaultBindRPC; }
	if rpc.bind == "" { return Errors.New("Bind address is required, in ServerRPC->Start()"); }
	protocol, address, port := PxnNet.SplitProtocolAddressPort(rpc.bind);
	if protocol == "" { return Errors.New("protocol is required, in ServerRPC->Start()"); }
	Log.Printf("Starting RPC Server.. %s", rpc.bind);
	if rpc.grpc_server == nil { rpc.grpc_server = GRPC.NewServer(); }
	switch protocol {
	case "unix":
		rpc.use_tls = false;
//TODO
panic("UNFINISHED UNIX RPC SERVER");
		break;
	case "tcp", "tcp4", "tcp6":
		if rpc.use_tls { Log.Printf("%sTLS Enabled",  LogPrefix);
		} else {         Log.Printf("%sTLS Disabled", LogPrefix); }
		if !PxnSan.IsSafeDomain(address) {
			return Fmt.Errorf("Invalid address: %s, in ServerRPC->Start()", address); }
		if port == 0 { return Fmt.Errorf("Invalid port: %d, in ServerRPC->Start()"); }
		listen, err := PxnNet.NewServerSocket(rpc.bind);
		if err != nil { return Fmt.Errorf("%s, failed to listen, in ServerRPC->Start()", err); }
		rpc.listen = listen;
		break;
	default: return Fmt.Errorf("Unknown protocol: %s, in ServerRPC->Start()", protocol);
	}
	go rpc.Serve();
	PxnUtil.SleepC();
	return nil;
}

func (rpc *ServerRPC) Serve() {
	rpc.service.WaitGroup.Add(1);
	defer func() {
		rpc.Close();
		rpc.service.WaitGroup.Done();
	}();
	rpc.service.AddClose(rpc);
	if err := rpc.grpc_server.Serve(rpc.listen); err != nil {
		Log.Printf("%s, in ServerRPC->Serve()", err); }
}



func (rpc *ServerRPC) Close() {
	rpc.service.WaitGroup.Add(1);
	defer rpc.service.WaitGroup.Done();
	rpc.mut_state.Lock();
	defer rpc.mut_state.Unlock();
	if rpc.listen != nil {
		if err := rpc.listen.Close(); err != nil {
			Log.Printf("%s, in ServerRPC->Close()", err); }
		rpc.grpc_server.GracefulStop();
		rpc.listen = nil;
	}
}



func (rpc *ServerRPC) GetServerGRPC() *GRPC.Server {
	return rpc.grpc_server;
}

func (rpc *ServerRPC) SetServerGRPC(grpc_server *GRPC.Server) *GRPC.Server {
	rpc.mut_state.Lock();
	defer rpc.mut_state.Unlock();
	previous := rpc.grpc_server;
	rpc.grpc_server = grpc_server;
	return previous;
}
