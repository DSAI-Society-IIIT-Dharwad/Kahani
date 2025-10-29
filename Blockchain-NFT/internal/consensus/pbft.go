package consensus

import (
	"errors"
	"fmt"
	"sync"

	"storytelling-blockchain/internal/types"
)

// Network abstracts the underlying transport used by the PBFT node.
type Network interface {
	Broadcast(sender string, msg Message) error
	Send(sender, recipient string, msg Message) error
}

// Signer signs and verifies consensus messages.
type Signer interface {
	Sign(data []byte) (string, error)
	Verify(sender string, data []byte, signature string) bool
}

// Finalizer is invoked once a block has been committed by PBFT.
type Finalizer func(block types.Block)

// BlockBuilder constructs a block from the provided transactions.
type BlockBuilder interface {
	BuildBlock(transactions []types.Transaction) (types.Block, error)
}

// Config encapsulates the dependencies required by a PBFT node.
type Config struct {
	ID             string
	Peers          []string
	FaultTolerance int
	Network        Network
	Signer         Signer
	Builder        BlockBuilder
	Finalize       Finalizer
}

// PBFTNode represents a single validator participating in PBFT consensus.
type PBFTNode struct {
	id             string
	peers          []string
	faultTolerance int
	network        Network
	signer         Signer
	builder        BlockBuilder
	finalize       Finalizer

	mu        sync.Mutex
	sequence  int
	instances map[int]*instance
}

type instance struct {
	block       types.Block
	prePrepare  *Message
	prepares    map[string]struct{}
	commits     map[string]struct{}
	commitFired bool
}

var (
	errMissingID        = errors.New("pbft: node id required")
	errMissingNetwork   = errors.New("pbft: network is required")
	errMissingBuilder   = errors.New("pbft: block builder is required")
	errMissingFinalize  = errors.New("pbft: finalize callback is required")
	errInvalidTolerance = errors.New("pbft: fault tolerance must be >= 0")
)

// NewPBFTNode constructs a PBFT node using the provided configuration.
func NewPBFTNode(cfg Config) (*PBFTNode, error) {
	if cfg.ID == "" {
		return nil, errMissingID
	}
	if cfg.Network == nil {
		return nil, errMissingNetwork
	}
	if cfg.Builder == nil {
		return nil, errMissingBuilder
	}
	if cfg.Finalize == nil {
		return nil, errMissingFinalize
	}
	if cfg.FaultTolerance < 0 {
		return nil, errInvalidTolerance
	}

	return &PBFTNode{
		id:             cfg.ID,
		peers:          append([]string(nil), cfg.Peers...),
		faultTolerance: cfg.FaultTolerance,
		network:        cfg.Network,
		signer:         cfg.Signer,
		builder:        cfg.Builder,
		finalize:       cfg.Finalize,
		instances:      make(map[int]*instance),
	}, nil
}

// ProposeBlock begins consensus for the provided transactions.
func (n *PBFTNode) ProposeBlock(transactions []types.Transaction) error {
	block, err := n.builder.BuildBlock(transactions)
	if err != nil {
		return err
	}

	n.mu.Lock()
	n.sequence++
	seq := n.sequence
	inst := &instance{
		block:    block,
		prepares: make(map[string]struct{}),
		commits:  make(map[string]struct{}),
	}
	n.instances[seq] = inst
	n.mu.Unlock()

	msg := Message{
		Type:     MessagePrePrepare,
		View:     0,
		Sequence: seq,
		Block:    block,
		SenderID: n.id,
	}

	if err := n.signMessage(&msg); err != nil {
		return err
	}

	if err := n.broadcast(msg); err != nil {
		return err
	}

	// allow node to process its own message to progress the round
	n.handlePrePrepare(msg)

	return nil
}

// HandleMessage routes an inbound consensus message to the appropriate handler.
func (n *PBFTNode) HandleMessage(msg Message) error {
	switch msg.Type {
	case MessagePrePrepare:
		n.handlePrePrepare(msg)
	case MessagePrepare:
		n.handlePrepare(msg)
	case MessageCommit:
		n.handleCommit(msg)
	default:
		return fmt.Errorf("pbft: unsupported message type %s", msg.Type)
	}
	return nil
}

func (n *PBFTNode) handlePrePrepare(msg Message) {
	if !n.verifyMessage(msg) {
		return
	}

	n.mu.Lock()
	inst, ok := n.instances[msg.Sequence]
	if !ok {
		inst = &instance{
			block:    msg.Block,
			prepares: make(map[string]struct{}),
			commits:  make(map[string]struct{}),
		}
		n.instances[msg.Sequence] = inst
	}

	if inst.prePrepare != nil {
		n.mu.Unlock()
		return
	}

	inst.prePrepare = &msg
	n.mu.Unlock()

	prepare := Message{
		Type:     MessagePrepare,
		View:     msg.View,
		Sequence: msg.Sequence,
		Block:    msg.Block,
		SenderID: n.id,
	}

	if err := n.signMessage(&prepare); err != nil {
		return
	}

	_ = n.broadcast(prepare)
	n.handlePrepare(prepare)
}

func (n *PBFTNode) handlePrepare(msg Message) {
	if !n.verifyMessage(msg) {
		return
	}

	n.mu.Lock()
	inst, ok := n.instances[msg.Sequence]
	if !ok || inst.prePrepare == nil {
		n.mu.Unlock()
		return
	}

	if inst.block.Hash != msg.Block.Hash {
		n.mu.Unlock()
		return
	}

	if _, exists := inst.prepares[msg.SenderID]; exists {
		n.mu.Unlock()
		return
	}

	inst.prepares[msg.SenderID] = struct{}{}
	count := len(inst.prepares) + 1 // include self pre-prepare
	threshold := n.quorumSize()
	emitCommit := count >= threshold
	n.mu.Unlock()

	if emitCommit {
		commit := Message{
			Type:     MessageCommit,
			View:     msg.View,
			Sequence: msg.Sequence,
			Block:    msg.Block,
			SenderID: n.id,
		}

		if err := n.signMessage(&commit); err == nil {
			_ = n.broadcast(commit)
		}
		n.handleCommit(commit)
	}
}

func (n *PBFTNode) handleCommit(msg Message) {
	if !n.verifyMessage(msg) {
		return
	}

	n.mu.Lock()
	inst, ok := n.instances[msg.Sequence]
	if !ok || inst.prePrepare == nil {
		n.mu.Unlock()
		return
	}

	if inst.block.Hash != msg.Block.Hash {
		n.mu.Unlock()
		return
	}

	if _, exists := inst.commits[msg.SenderID]; exists {
		n.mu.Unlock()
		return
	}

	inst.commits[msg.SenderID] = struct{}{}

	commitCount := len(inst.commits)
	threshold := n.quorumSize()

	if commitCount >= threshold && !inst.commitFired {
		inst.commitFired = true
		block := inst.block
		n.mu.Unlock()
		n.finalize(block)
		return
	}

	n.mu.Unlock()
}

func (n *PBFTNode) quorumSize() int {
	return 2*n.faultTolerance + 1
}

func (n *PBFTNode) signMessage(msg *Message) error {
	if n.signer == nil {
		return nil
	}

	digest, err := msg.Digest()
	if err != nil {
		return err
	}

	signature, err := n.signer.Sign(digest)
	if err != nil {
		return err
	}

	msg.Signature = signature
	return nil
}

func (n *PBFTNode) verifyMessage(msg Message) bool {
	if n.signer == nil || msg.Signature == "" {
		return true
	}

	digest, err := msg.Digest()
	if err != nil {
		return false
	}

	return n.signer.Verify(msg.SenderID, digest, msg.Signature)
}

func (n *PBFTNode) broadcast(msg Message) error {
	if n.network == nil {
		return errors.New("pbft: network not configured")
	}

	if err := n.network.Broadcast(n.id, msg); err != nil {
		return err
	}

	return nil
}
