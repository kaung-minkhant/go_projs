package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(http.ResponseWriter, *http.Request)
}

type LoadBalancer interface {
	getNextAvailableServer() Server
	serveProxy(http.ResponseWriter, *http.Request)
}

type SimpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}

type SimpleLoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []Server
}

func NewSimpleServer(addr string) *SimpleServer {
	serverUrl, err := url.Parse(addr)
	handleError(err)

	return &SimpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

func NewSimpleLoadbalancer(port string, servers []Server) *SimpleLoadBalancer {
	return &SimpleLoadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("Error %s\n", err)
	}
}

func (ss *SimpleServer) Address() string {
	return ss.addr
}

func (ss *SimpleServer) IsAlive() bool {
	return true
}

func (ss *SimpleServer) Serve(w http.ResponseWriter, r *http.Request) {
	ss.proxy.ServeHTTP(w, r)
}

func (lb *SimpleLoadBalancer) getNextAvailableServer() Server {
	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	lb.roundRobinCount++
	for !server.IsAlive() {
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	return server

}

func (lb *SimpleLoadBalancer) serveProxy(w http.ResponseWriter, r *http.Request) {
	server := lb.getNextAvailableServer()
	fmt.Printf("forwarding request to address %s\n", server.Address())
	server.Serve(w, r)
}

func main() {
	servers := []Server{
		// NewSimpleServer("https://www.google.com"),
		NewSimpleServer("https://www.github.com"),
		NewSimpleServer("https://www.duckduckgo.com"),
	}

	lb := NewSimpleLoadbalancer(":8080", servers)

	startLoadBalancer(lb, lb.port)
}

func startLoadBalancer(lb LoadBalancer, port string) {
	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		lb.serveProxy(w, r)
	}

	http.HandleFunc("/", handleRedirect)

	fmt.Printf("Starting Loadbalance at port localhost%s\n", port)

	log.Fatal(http.ListenAndServe(port, nil))
}
