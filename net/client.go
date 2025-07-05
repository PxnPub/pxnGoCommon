package net;

import(
	Fmt    "fmt"
	Net    "net"
	Errors "errors"
	PxnSan "github.com/PoiXson/pxnGoCommon/utils/san"
);



func NewClientSocket(remote string) (Net.Conn, error) {
	if remote == "" { return nil, Errors.New("remote address required"); }
	protocol, addr, port := SplitProtocolAddressPort(remote);
	if protocol == "" { return nil, Errors.New("protocol is required"); }
	switch protocol {
	case "unix":
		if len(addr) < 5 { return nil, Fmt.Errorf("Invalid unix socket: %s", addr ); }
		resolved, err := Net.ResolveUnixAddr(protocol, addr);
		if err != nil { return nil, err; }
		conn, err := Net.DialUnix(protocol, nil, resolved);
		if err != nil { return nil, err; }
		return conn, nil;
	case "tcp", "tcp4", "tcp6":
		if !PxnSan.IsSafeDomain(addr) { return nil, Fmt.Errorf("Invalid address: %s", addr); }
		if port == 0                  { return nil, Fmt.Errorf("Invalid port: %d"); }
		addrport := Fmt.Sprintf("%s:%d", addr, port);
		resolved, err := Net.ResolveTCPAddr(protocol, addrport);
		if err != nil { return nil, err; }
		conn, err := Net.DialTCP(protocol, nil, resolved);
		if err != nil { return nil, err; }
		return conn, nil;
	default: break;
	}
	return nil, Fmt.Errorf("Unknown protocol: %s", protocol);
}
