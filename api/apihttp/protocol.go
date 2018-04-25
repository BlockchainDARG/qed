package apihttp

import (
	"fmt"
	"verifiabledata/balloon"
	"verifiabledata/balloon/history"
)

type Snapshot struct {
	HyperDigest   string
	HistoryDigest string
	Version       uint
	Event         string
	//TODO: implement this
	// EventDigest   string
}

type HistoryNode struct {
	Digest       string
	Index, Layer uint
}

type Proofs struct {
	HyperAuditPath   []string
	HistoryAuditPath []HistoryNode
}

type MembershipProof struct {
	Key                         string
	KeyDigest                   string
	IsMember                    bool
	Proofs                      *Proofs
	QueryVersion, ActualVersion uint
}

func assemblySnapshot(commitment *balloon.Commitment, event string) *Snapshot {
	return &Snapshot{
		fmt.Sprintf("%064x", commitment.HyperDigest),
		fmt.Sprintf("%064x", commitment.HistoryDigest),
		commitment.Version,
		event,
	}
}

func assemblyHyperAuditPath(path [][]byte) []string {
	result := make([]string, 0)
	for _, elem := range path {
		result = append(result, fmt.Sprintf("%064x", elem))
	}
	return result
}

func assemblyHistoryAuditPath(path []history.Node) []HistoryNode {
	result := make([]HistoryNode, 0)
	for _, elem := range path {
		result = append(result, assemblyHistoryNode(elem))
	}
	return result
}

func assemblyHistoryNode(node history.Node) HistoryNode {
	return HistoryNode{
		fmt.Sprintf("%064x", node.Digest),
		node.Index,
		node.Layer,
	}
}

func assemblyMembershipProof(event string, proof *balloon.MembershipProof) *MembershipProof {
	return &MembershipProof{
		event,
		fmt.Sprintf("%064x", proof.KeyDigest),
		proof.Exists,
		&Proofs{
			assemblyHyperAuditPath(proof.HyperProof),
			assemblyHistoryAuditPath(proof.HistoryProof),
		},
		proof.QueryVersion,
		proof.ActualVersion,
	}
}