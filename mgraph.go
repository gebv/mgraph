package mgraph

import (
	"errors"
)

var (
	ErrGraphExists  = errors.New("graph exists")
	ErrNodeNotFound = errors.New("node not found")
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
