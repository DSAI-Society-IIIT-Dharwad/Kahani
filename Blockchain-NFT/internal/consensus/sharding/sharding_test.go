package sharding

import "testing"

func TestSelectNodeDeterministic(t *testing.T) {
	nodes := []string{"node-a", "node-b", "node-c"}

	cases := []string{"user-1", "user-2", "user-3", "", "another-user"}

	for _, key := range cases {
		expected := SelectNode(nodes, key)
		if expected == "" {
			t.Fatalf("expected non-empty selection for key %s", key)
		}

		for i := 0; i < 5; i++ {
			if got := SelectNode(nodes, key); got != expected {
				t.Fatalf("expected deterministic selection for key %s, iteration %d: %s vs %s", key, i, expected, got)
			}
		}
	}
}

func TestSelectNodeHandlesEdgeCases(t *testing.T) {
	if got := SelectNode(nil, "anything"); got != "" {
		t.Fatalf("expected empty string when nodes missing, got %s", got)
	}

	single := []string{"node-1"}
	if got := SelectNode(single, "user-1"); got != "node-1" {
		t.Fatalf("expected node-1 for single node, got %s", got)
	}
}
