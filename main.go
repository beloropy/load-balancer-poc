package main

import (
	"load-balancer/client"
	"load-balancer/load_balancer"
	"load-balancer/node"
	"net/http"
)

func main() {
	nodes := []node.Node{
		node.NewNode(
			1,
			client.NewClient("httpbin.org"),
			10,
			1,
		),
		node.NewNode(
			2,
			client.NewClient("httpbin.org"),
			20,
			2,
		),
	}

	if err := http.ListenAndServe(":8080", load_balancer.NewLoadBalancer(nodes)); err != nil {
		// TODO: Handle errors.
		return
	}
}
