package sharding

import "hash/fnv"

// SelectNode deterministically chooses a node identifier for the provided key.
func SelectNode(nodes []string, key string) string {
	if len(nodes) == 0 {
		return ""
	}

	if len(nodes) == 1 || key == "" {
		return nodes[0]
	}

	sum := fnv.New32a()
	_, _ = sum.Write([]byte(key))
	idx := int(sum.Sum32()) % len(nodes)
	if idx < 0 {
		idx += len(nodes)
	}

	return nodes[idx]
}
