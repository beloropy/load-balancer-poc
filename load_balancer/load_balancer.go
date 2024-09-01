package load_balancer

import (
	"load-balancer/node"
	"net/http"
)

type LoadBalancer interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

var _ LoadBalancer = (*defaultLoadBalancer)(nil)

func NewLoadBalancer(nodes []node.Node) LoadBalancer {
	return &defaultLoadBalancer{
		nodes: nodes,
	}
}

type defaultLoadBalancer struct {
	nodes []node.Node
}

func (l *defaultLoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var tryServeHTTP bool
	for _, n := range l.nodes {
		if tryServeHTTP = n.TryServeHTTP(w, r); tryServeHTTP {
			break
		}
	}

	if !tryServeHTTP {
		http.Error(w, "429 Too Many Requests", http.StatusTooManyRequests)
		return
	}

	return
}
