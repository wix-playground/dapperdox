package navigation

import ()

type NavigationNode struct {
	ChildMap  map[string]*NavigationNode
	Children  []*NavigationNode
	SortOrder string
	Name      string
	Id        string
	Uri       string
}

type ByOrder []*NavigationNode

func (n ByOrder) Len() int {
	return len(n)
}
func (n ByOrder) Less(a, b int) bool {
	return n[a].SortOrder < n[b].SortOrder
}
func (n ByOrder) Swap(a, b int) {
	n[a], n[b] = n[b], n[a]
}
