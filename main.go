package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/Konstantin8105/CalculixRPCserver/serverCalculix"
)

// example of calculix RPC server
func main() {
	calculix := new(serverCalculix.Calculix)
	err := rpc.Register(calculix)
	if err != nil {
		fmt.Println("Cannot register the calculix")
		return
	}
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	err = http.Serve(l, nil)
	if err != nil {
		fmt.Println("Cannot serve the calculix")
		return
	}
}
