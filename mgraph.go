package mgraph

import (
	"errors"
)

var (
	ErrGraphExists                 = errors.New("mgraph: graph exists")
	ErrNodeNotFound                = errors.New("mgraph: node not found")
	ErrNotAllowedMoveInOwnSubthred = errors.New("mgraph: not allowed to move in own subthread")
)

type Graph struct {
	GraphID     int64
	RootNodeID  int64
	HashedExtID string
}

type Node struct {
	NodeID   int64
	ParentID *int64
	Path     *string

	// virtual fields
	Level int
}
