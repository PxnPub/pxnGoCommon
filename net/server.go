package net;

import(
	Fmt    "fmt"
	Net    "net"
	Errors "errors"
	PxnSan "github.com/PoiXson/pxnGoCommon/utils/san"
);



func NewServerSocket(bind string) (Net.Listener, error) {
	if bind == "" { return nil, Errors.New("bind address required"); }
	protocol, addr, port := SplitProtocolAddressPort(bind);
	if protocol == "" { return nil, Errors.New("protocol is required"); }
	switch protocol {
	case "unix":
		if len(addr) < 5 { return nil, Fmt.Errorf("Invalid unix socket: %s", addr); }
		if err := RemoveOldUnixSocket(addr); err != nil { return nil, err; }
		resolved, err := Net.ResolveUnixAddr(protocol, addr);
		if err != nil { return nil, err; }
		listen, err := Net.ListenUnix(protocol, resolved);
		if err != nil { return nil, err; }
		return listen, nil;
	case "tcp", "tcp4", "tcp6":
		if !PxnSan.IsSafeDomain(addr) { return nil, Fmt.Errorf("Invalid address: %s", addr); }
		if port == 0                  { return nil, Fmt.Errorf("Invalid port: %d"); }
		addrport := Fmt.Sprintf("%s:%d", addr, port);
		resolved, err := Net.ResolveTCPAddr(protocol, addrport);
		if err != nil { return nil, err; }
		listen, err := Net.ListenTCP(protocol, resolved);
		if err != nil { return nil, err; }
		return listen, nil;
	default: break;
	}
	return nil, Fmt.Errorf("Unknown protocol: %s", protocol);
}



func NewServerUDP(bind string) (*Net.UDPConn, error) {
	if bind == "" { return nil, Errors.New("bind address required"); }
	protocol, addr, port := SplitProtocolAddressPort(bind);
	if protocol == "" { return nil, Errors.New("protocol is required"); }
	switch protocol {
	case "udp", "udp4", "udp6":
		if !PxnSan.IsSafeDomain(addr) { return nil, Fmt.Errorf("Invalid address: %s", addr); }
		if port == 0                  { return nil, Fmt.Errorf("Invalid port: %d"); }
		addrport := Fmt.Sprintf("%s:%d", addr, port);
		resolved, err := Net.ResolveUDPAddr(protocol, addrport);
		if err != nil { return nil, err; }
		listen, err := Net.ListenUDP(protocol, resolved);
		if err != nil { return nil, err; }
		return listen, nil;
	default: break;
	}
	return nil, Fmt.Errorf("Unknown protocol: %s", protocol);
}
