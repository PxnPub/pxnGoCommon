package ratelimit;

import(
	Net    "net"
	PxnNet "github.com/PoiXson/pxnGoCommon/net"
);



type RateLimit interface {
	Start()
	Tick()
	CheckNetAddr(addr Net.Addr) (bool, error)
	CheckStrAddr(addr string  ) (bool, error)
	CheckTupleIP(ip *PxnNet.TupIP) bool
}
