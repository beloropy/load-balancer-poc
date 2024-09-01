package node

import (
	"fmt"
	"load-balancer/client"
	"net/http"
	"sync"
	"time"
)

type Node interface {
	TryServeHTTP(w http.ResponseWriter, r *http.Request) bool
}

func NewNode(nodeId int, client client.Client, bpmLimit int64, rpmLimit int32) Node {
	n := &node{
		nodeId,
		client,
		bpmLimit,
		0,
		rpmLimit,
		0,
		sync.RWMutex{},
		time.NewTicker(time.Minute),
	}

	go func() {
		for range n.resetTicker.C {
			n.resetRateLimits()
		}
	}()

	return n
}

type node struct {
	nodeId int
	client client.Client

	// Each node can have a different rate limit.
	// Rate limits are measured in two ways: BPMLimit (http body Bytes Per Minute), RPMLimit (Requests Per Minute)
	BPMLimit   int64
	currentBPM int64

	RPMLimit   int32
	currentRPM int32

	mu sync.RWMutex

	resetTicker *time.Ticker
}

func (n *node) TryServeHTTP(w http.ResponseWriter, r *http.Request) bool {
	bodyBytes := r.ContentLength
	if bodyBytes == -1 {
		// TODO: Handle errors.
		return false
	}

	if !n.checkAndUpdateRateLimits(bodyBytes) {
		fmt.Printf("node id %d (%s): 429 Too Many Requests\n", n.nodeId, n.client.GetEndpoint())
		return false
	}

	n.client.ServeHTTP(w, r)
	fmt.Printf("node id %d (%s) serves a http request\n", n.nodeId, n.client.GetEndpoint())

	return true
}

func (n *node) checkAndUpdateRateLimits(bodyBytes int64) bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Rate limits can be hit across any of the options depending on what occurs first.
	if n.currentBPM+bodyBytes > n.BPMLimit {
		return false
	}

	if n.currentRPM+1 > n.RPMLimit {
		return false
	}

	n.currentBPM += bodyBytes
	n.currentRPM++

	return true
}

func (n *node) resetRateLimits() {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.currentBPM, n.currentRPM = 0, 0
	fmt.Printf("Rate limits of node id %d (%s) has been reset, current time: %v\n", n.nodeId, n.client.GetEndpoint(), time.Now())
}
