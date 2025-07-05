package web;

import(
	Log     "log"
	Fmt     "fmt"
	Net     "net"
	HTTP    "net/http"
	Sync    "sync"
	Time    "time"
	Strings "strings"
	Context "context"
	Errors  "errors"
	Gorilla "github.com/gorilla/mux"
	GHands  "github.com/gorilla/handlers"
	PxnUtil "github.com/PoiXson/pxnGoCommon/utils"
	PxnNet  "github.com/PoiXson/pxnGoCommon/net"
	PxnServ "github.com/PoiXson/pxnGoCommon/service"
);



type WebServer struct {
	mut_state Sync.Mutex
	service   *PxnServ.Service
	// transport
	bind      string
	use_tls   bool
	proxied   bool
	listen    Net.Listener
	Router    HTTP.Handler
	server    *HTTP.Server
}



func NewWebServer(service *PxnServ.Service, bind string, proxied bool) *WebServer {
	web := WebServer{
		service: service,
		bind:    bind,
		proxied: proxied,
		server:  &HTTP.Server{},
	};
	return &web;
}



func (web *WebServer) Start() error {
	web.mut_state.Lock();
	defer web.mut_state.Unlock();
	if web.bind == "" { web.bind = DefaultBindWeb; }
	if web.bind == "" { return Errors.New("Bind address is required"); }
	listen, err := PxnNet.NewServerSocket(web.bind);
	if err != nil { return Fmt.Errorf(
		"%s for NewServerSocket() in WebServer->Start()", err); }
	web.listen = listen;
	go web.Serve();
	PxnUtil.SleepC();
	return nil;
}

func (web *WebServer) Close() {
	web.service.WaitGroup.Add(1);
	web.mut_state.Lock();
	defer func() {
		web.mut_state.Unlock();
		web.service.WaitGroup.Done();
	}();
	if web.listen != nil {
		if err := web.listen.Close(); err != nil {
			Log.Printf("%v, in WebServer->Close()", err); }
		if err := web.server.Shutdown(Context.Background()); err != nil {
			Log.Printf("%v, in WebServer->Close()", err); }
		web.listen = nil;
	}
}



func (web *WebServer) Serve() {
	web.service.WaitGroup.Add(1);
	defer func() {
		web.Close();
		web.service.WaitGroup.Done();
	}();
	web.service.AddClose(web);
	Log.Printf("Starting Web Server.. %s", web.bind);
	web.server.Handler = web.Router;
	if err := web.server.Serve(web.listen); err != nil {
		if !Strings.HasSuffix(err.Error(), "use of closed network connection") {
			Log.Printf("%v, in WebServer->Serve()", err); }}
}



// handler
func (web *WebServer) WithHandler(router HTTP.Handler) *WebServer {
	web.Router = router;
	return web;
}

// gorilla router
func (web *WebServer) WithGorilla() *Gorilla.Router {
	router := Gorilla.NewRouter();
	router.NotFoundHandler = HTTP.HandlerFunc(web.PageNotFound);
	if web.proxied { router.Use(GHands.ProxyHeaders); }
	router.Use(web.MiddlewareStats);
	web.Router = router;
	return router;
}

//TODO
/*
func (web *WebServer) WithRobots() *WebServer {
//User-agent: *
//Disallow: /build/
//Disallow: /task/
	return web;
}
*/



func (web *WebServer) MiddlewareStats(next HTTP.Handler) HTTP.Handler {
	return HTTP.HandlerFunc(func(w HTTP.ResponseWriter, r *HTTP.Request) {
		start := Time.Now();
		next.ServeHTTP(w, r);
		Log.Printf("%s%s %s in %v", LogPrefixWeb, r.Method, r.URL.Path, Time.Since(start));
	});
}

func (web *WebServer) PageNotFound(w HTTP.ResponseWriter, r *HTTP.Request) {
	HTTP.Error(w, "404 Not Found", HTTP.StatusNotFound);
	Log.Printf("%s404 %s %s", LogPrefixWeb, r.Method, r.URL.Path);
}



func AddStaticRoute(router HTTP.Handler) {
	fs := HTTP.FileServer(HTTP.Dir("./static"));
	// gorilla mux
	if mux, ok := router.(*Gorilla.Router); ok {
		mux.PathPrefix("/static/").Handler(HTTP.StripPrefix("/static/", fs));
	} else
	// std http mux
	if mux, ok := router.(*HTTP.ServeMux); ok {
		mux.Handle("/static/", HTTP.StripPrefix("/static/", fs));
	// unknown mux
	} else {
		Log.Panicf("Unsupported mux type: %T", router);
	}
}



func NewRedirect(target string) HTTP.HandlerFunc  {
	return func(out HTTP.ResponseWriter, in *HTTP.Request) {
		HTTP.Redirect(out, in, target, HTTP.StatusFound);
	};
}



func (web *WebServer) AddIconFile(file string) {
//TODO
//web.Router
}



//TODO
/*
func (web *WebServer) GetStats() *Stats {
	type StatsRPC struct {
		CountConns uint64
		CountReqs  uint64
	}
}
*/
