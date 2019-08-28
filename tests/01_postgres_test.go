package tests

import (
	"testing"

	"github.com/gebv/mgraph"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test01_01MoveCases(t *testing.T) {
	t.Run("Case1", func(t *testing.T) {
		m := newTestManipulator(t)
		/*
			a1
				a1b1
				a1b2
			a2
			a3
				a3b1
				a3b2
		*/

		m.AddNodeToROOT("a1")

		m.AddNode("a1b1", "a1")
		m.AddNode("a1b2", "a1")

		m.AddNodeToROOT("a2")

		m.AddNodeToROOT("a3")
		m.AddNode("a3b1", "a3")
		m.AddNode("a3b2", "a3")

		want := []nodeWithAlias{
			{"a1", "root"},
			{"a1b1", "root.a1"},
			{"a1b2", "root.a1"},
			{"a2", "root"},
			{"a3", "root"},
			{"a3b1", "root.a3"},
			{"a3b2", "root.a3"},
		}
		m.AssertGraphNodes(want)

		err := m.MoveNode("a3", "a2")
		require.NoError(t, err, "move a3 to a2")

		/*
			a1
				a1b1
				a1b2
			a2
				a3
					a3b1
					a3b2
		*/

		want = []nodeWithAlias{
			{"a1", "root"},
			{"a1b1", "root.a1"},
			{"a1b2", "root.a1"},
			{"a2", "root"},
			{"a3", "root.a2"},
			{"a3b1", "root.a2.a3"},
			{"a3b2", "root.a2.a3"},
		}
		m.AssertGraphNodes(want)
	})

	t.Run("Case2", func(t *testing.T) {
		m := newTestManipulator(t)
		/*
			a1
				a1b1
					a1b1c1
						a1b1c1d1
			a2
		*/

		m.AddNodeToROOT("a1")
		m.AddNodeToROOT("a2")

		m.AddNode("a1b1", "a1")
		m.AddNode("a1b1c1", "a1b1")
		m.AddNode("a1b1c1d1", "a1b1c1")

		want := []nodeWithAlias{
			{"a1", "root"},
			{"a1b1", "root.a1"},
			{"a1b1c1", "root.a1.a1b1"},
			{"a1b1c1d1", "root.a1.a1b1.a1b1c1"},
			{"a2", "root"},
		}
		m.AssertGraphNodes(want)

		err := m.MoveNode("a1b1", "a2")
		require.NoError(t, err, "move a1b1 to a2")

		err = m.MoveNode("a2", "a1")
		require.NoError(t, err, "move a1b1 to a2")

		/*
			a1
				a2
					a1b1
						a1b1c1
							a1b1c1d1
		*/

		want = []nodeWithAlias{
			{"a1", "root"},
			{"a2", "root.a1"},
			{"a1b1", "root.a1.a2"},
			{"a1b1c1", "root.a1.a2.a1b1"},
			{"a1b1c1d1", "root.a1.a2.a1b1.a1b1c1"},
		}
		m.AssertGraphNodes(want)
	})

	t.Run("Case3", func(t *testing.T) {
		m := newTestManipulator(t)
		/*
			1
				2
					3
						4
							5
		*/

		m.AddNodeToROOT("1")

		m.AddNode("2", "1")
		m.AddNode("3", "2")
		m.AddNode("4", "3")
		m.AddNode("5", "4")

		want := []nodeWithAlias{
			{"1", "root"},
			{"2", "root.1"},
			{"3", "root.1.2"},
			{"4", "root.1.2.3"},
			{"5", "root.1.2.3.4"},
		}
		m.AssertGraphNodes(want)

		err := m.MoveNode("3", "1")
		require.NoError(t, err, "move 3 to 1")

		/*
			1
				2
				3
					4
						5
		*/

		want = []nodeWithAlias{
			{"1", "root"},
			{"2", "root.1"},
			{"3", "root.1"},
			{"4", "root.1.3"},
			{"5", "root.1.3.4"},
		}
		m.AssertGraphNodes(want)

		err = m.MoveNode("2", "5")
		require.NoError(t, err, "move a1b1 to a2")

		/*
			1
				3
					4
						5
							2
		*/

		want = []nodeWithAlias{
			{"1", "root"},
			{"3", "root.1"},
			{"4", "root.1.3"},
			{"5", "root.1.3.4"},
			{"2", "root.1.3.4.5"},
		}
		m.AssertGraphNodes(want)
	})

	t.Run("Case4", func(t *testing.T) {
		m := newTestManipulator(t)
		/*
			1
				2
					3
						4
							5
		*/

		m.AddNodeToROOT("1")

		m.AddNode("2", "1")
		m.AddNode("3", "2")
		m.AddNode("4", "3")
		m.AddNode("5", "4")

		want := []nodeWithAlias{
			{"1", "root"},
			{"2", "root.1"},
			{"3", "root.1.2"},
			{"4", "root.1.2.3"},
			{"5", "root.1.2.3.4"},
		}
		m.AssertGraphNodes(want)

		err := m.MoveNode("2", "5")
		require.Error(t, err, mgraph.ErrNotAllowedMoveInOwnSubthred, "move 2 to 5")

		/*
			1
				2
					3
						4
							5
		*/

		want = []nodeWithAlias{
			{"1", "root"},
			{"2", "root.1"},
			{"3", "root.1.2"},
			{"4", "root.1.2.3"},
			{"5", "root.1.2.3.4"},
		}
		m.AssertGraphNodes(want)
	})
}

func Test01_02RemoveCases(t *testing.T) {
	t.Run("Case1", func(t *testing.T) {
		m := newTestManipulator(t)
		/*
			1
				2
					3
						4
							5
		*/

		m.AddNodeToROOT("1")

		m.AddNode("2", "1")
		m.AddNode("3", "2")
		m.AddNode("4", "3")
		m.AddNode("5", "4")

		want := []nodeWithAlias{
			{"1", "root"},
			{"2", "root.1"},
			{"3", "root.1.2"},
			{"4", "root.1.2.3"},
			{"5", "root.1.2.3.4"},
		}
		m.AssertGraphNodes(want)

		err := m.RemoveNode("3")
		require.NoError(t, err, "remove 3")

		/*
			1
				2
		*/

		want = []nodeWithAlias{
			{"1", "root"},
			{"2", "root.1"},
		}
		m.AssertGraphNodes(want)

		err = m.RemoveNode("2")
		require.NoError(t, err, "remove 2")

		/*
			1
		*/

		want = []nodeWithAlias{
			{"1", "root"},
		}
		m.AssertGraphNodes(want)
	})
}

func Test01_03ExceptionalSituations(t *testing.T) {
	t.Run("addToNonExistingParent", func(t *testing.T) {
		m := newTestManipulator(t)

		m.AddNodeToROOT("1")
		_, err := m.s.Add(Ctx, 123123123)
		assert.Error(t, err, mgraph.ErrNodeNotFound)

		err = m.s.Remove(Ctx, 123123123)
		assert.Error(t, err, mgraph.ErrNodeNotFound)
	})

	t.Run("moveAndRemoveToNotExistsParent", func(t *testing.T) {
		m := newTestManipulator(t)

		m.AddNodeToROOT("1")

		m.AddNode("2", "1")
		m.AddNode("3", "2")
		m.AddNode("4", "3")
		m.AddNode("5", "4")

		want := []nodeWithAlias{
			{"1", "root"},
			{"2", "root.1"},
			{"3", "root.1.2"},
			{"4", "root.1.2.3"},
			{"5", "root.1.2.3.4"},
		}
		m.AssertGraphNodes(want)

		err := m.s.Move(Ctx, m.names["2"], 123123)
		require.Error(t, err, mgraph.ErrNodeNotFound)

		want = []nodeWithAlias{
			{"1", "root"},
			{"2", "root.1"},
			{"3", "root.1.2"},
			{"4", "root.1.2.3"},
			{"5", "root.1.2.3.4"},
		}
		m.AssertGraphNodes(want)

		err = m.s.Remove(Ctx, 123123)
		require.Error(t, err, mgraph.ErrNodeNotFound)

		want = []nodeWithAlias{
			{"1", "root"},
			{"2", "root.1"},
			{"3", "root.1.2"},
			{"4", "root.1.2.3"},
			{"5", "root.1.2.3.4"},
		}
		m.AssertGraphNodes(want)
	})
}
