package tests

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gebv/mgraph"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	DB  *sql.DB
	Ctx context.Context
	MGP *mgraph.MGraphPostgres
)

func newTestManipulator(t *testing.T) *testManipulator {
	graphName := fmt.Sprint("ts-", time.Now().UnixNano())
	graphID, rootID, err := MGP.CreateGraph(Ctx, graphName)
	assert.NoError(t, err, "create graph")
	assert.NotEmpty(t, rootID)
	t.Logf("create new graph=%q root_id=%d\n", graphName, rootID)

	return &testManipulator{
		graphID: graphID,
		rootID:  rootID,
		t:       t,
		s:       MGP,
		names: map[string]int64{
			"root": rootID,
		},
	}
}

type testManipulator struct {
	rootID  int64
	graphID int64
	t       *testing.T
	s       *mgraph.MGraphPostgres
	names   map[string]int64
}

func (g *testManipulator) AddNodeToROOT(name string) {
	nodeID, err := g.s.Add(context.Background(), g.rootID)
	require.NoError(g.t, err, "failed add new node")
	g.names[name] = nodeID
}

func (g *testManipulator) AddNode(name, parentName string) {
	parentID := g.names[parentName]
	require.NotEmpty(g.t, parentID, "failed find parent node ID by name")
	nodeID, err := g.s.Add(context.Background(), parentID)
	require.NoError(g.t, err, "failed add new node")
	require.NotEmpty(g.t, parentID, "empty new node ID")
	g.names[name] = nodeID
}

func (g *testManipulator) MoveNode(name, newParentName string) error {
	parentID := g.names[newParentName]
	require.NotEmpty(g.t, parentID, "failed find parent node ID by name")
	nodeID := g.names[name]
	require.NotEmpty(g.t, nodeID, "failed find node ID by name")
	g.t.Logf("MOVE(%q, %q), MOVE(%d, %d)", name, newParentName, nodeID, parentID)
	if err := g.s.Move(context.Background(), nodeID, parentID); err != nil {
		return err
	}
	return nil
}

func (g *testManipulator) RemoveNode(name string) error {
	nodeID := g.names[name]
	require.NotEmpty(g.t, nodeID, "failed find node ID by name")
	return g.s.Remove(context.Background(), nodeID)
}

type nodeWithAlias struct {
	Name string
	Path string
}

func (g *testManipulator) findNameByID(nodeID int64) string {
	for name, id := range g.names {
		if id == nodeID {
			return name
		}
	}
	return ""
}

func (g *testManipulator) replacePath(path string) string {
	res := []string{}
	for _, nodeStrID := range strings.Split(path, ".") {
		nodeID, _ := strconv.ParseInt(nodeStrID, 10, 64)
		nodeName := g.findNameByID(nodeID)
		if nodeName == "" {
			g.t.Errorf("not found node name by ID=%d", nodeID)
		}
		res = append(res, nodeName)
	}
	return strings.Join(res, ".")
}

func (g *testManipulator) GraphNodes() []nodeWithAlias {
	nodes := g.s.GraphNodes(context.Background(), g.graphID)
	res := []nodeWithAlias{}
	for _, node := range nodes {
		res = append(res, nodeWithAlias{
			Name: g.findNameByID(node.NodeID),
			Path: g.replacePath(strOr(node.Path)),
		})
	}
	return res
}

func (g *testManipulator) AssertGraphNodes(want []nodeWithAlias) bool {
	wantMap := map[string]nodeWithAlias{}
	for _, item := range want {
		wantMap[item.Name] = item
	}

	gotMap := map[string]nodeWithAlias{}
	for _, item := range g.GraphNodes() {
		gotMap[item.Name] = item
	}
	return assert.EqualValues(g.t, wantMap, gotMap)
}

func strOr(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

// func clearAll(t *testing.T, db *sql.DB) {
// 	_, err := db.Exec(`DELETE FROM mgraph.graph;`)
// 	if err != nil && err != sql.ErrNoRows {
// 		t.Fatal(err, "failed clear TABLE mgraph.graph")
// 	}
// 	_, err = db.Exec(`DELETE FROM mgraph.graph_nodes;`)
// 	if err != nil && err != sql.ErrNoRows {
// 		t.Fatal(err, "failed clear TABLE mgraph.graph_nodes")
// 	}
// }
