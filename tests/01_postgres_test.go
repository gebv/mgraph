package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Simple1(t *testing.T) {
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

	// TODO: moving in the same branch

}
