package network

import (
	"errors"
	"sync"
)

// Message represents a payload sent between peers.
type Message struct {
	From    string
	Payload []byte
}

// Transport routes messages between nodes.
type Transport interface {
	Send(from, to string, payload []byte) error
}

// Node represents a peer in the P2P network.
type Node struct {
	id        string
	transport Transport

	mu    sync.RWMutex
	peers map[string]struct{}

	incoming chan Message
}

// NewNode creates a node with the provided id and transport.
func NewNode(id string, transport Transport) *Node {
	return &Node{
		id:        id,
		transport: transport,
		peers:     make(map[string]struct{}),
		incoming:  make(chan Message, 32),
	}
}

// ID returns the node identifier.
func (n *Node) ID() string {
	return n.id
}

// ConnectToPeer registers the peer with the node.
func (n *Node) ConnectToPeer(peerID string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.peers[peerID] = struct{}{}
}

// Peers returns the list of connected peer identifiers.
func (n *Node) Peers() []string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	ids := make([]string, 0, len(n.peers))
	for id := range n.peers {
		ids = append(ids, id)
	}
	return ids
}

// SendMessage sends a payload to the specified peer via the transport.
func (n *Node) SendMessage(peerID string, payload []byte) error {
	if n.transport == nil {
		return errors.New("network: transport not configured")
	}

	n.mu.RLock()
	_, known := n.peers[peerID]
	n.mu.RUnlock()
	if !known {
		return errors.New("network: peer not connected")
	}

	return n.transport.Send(n.id, peerID, payload)
}

// ReceiveMessages returns a read-only channel for inbound messages.
func (n *Node) ReceiveMessages() <-chan Message {
	return n.incoming
}

// deliver enqueues a message for this node.
func (n *Node) deliver(from string, payload []byte) {
	select {
	case n.incoming <- Message{From: from, Payload: payload}:
	default:
		// drop message if channel full to avoid blocking
	}
}

// InMemoryTransport is a simple transport for tests and local development.
type InMemoryTransport struct {
	mu    sync.RWMutex
	nodes map[string]*Node
}

// NewInMemoryTransport creates an empty transport registry.
func NewInMemoryTransport() *InMemoryTransport {
	return &InMemoryTransport{nodes: make(map[string]*Node)}
}

// Register registers the node with the transport.
func (t *InMemoryTransport) Register(node *Node) {
	t.mu.Lock()
	t.nodes[node.ID()] = node
	t.mu.Unlock()
}

// Send delivers the payload to the intended recipient.
func (t *InMemoryTransport) Send(from, to string, payload []byte) error {
	t.mu.RLock()
	node, ok := t.nodes[to]
	t.mu.RUnlock()
	if !ok {
		return errors.New("network: recipient unknown")
	}

	node.deliver(from, payload)
	return nil
}
