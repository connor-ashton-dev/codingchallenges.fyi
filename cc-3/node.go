package main

type Node struct {
	char  rune
	freq  int
	left  *Node
	right *Node
}

type NodeQueue []*Node

func (nq NodeQueue) Len() int {
	return len(nq)
}

func (nq NodeQueue) Less(i, j int) bool {
	if nq[i].freq == nq[j].freq {
		return nq[i].char < nq[j].char
	}
	return nq[i].freq < nq[j].freq
}

func (nq NodeQueue) Swap(i, j int) {
	nq[i], nq[j] = nq[j], nq[i]
}

func (nq *NodeQueue) Push(x interface{}) {
	item := x.(*Node)
	*nq = append(*nq, item)
}

func (nq *NodeQueue) Pop() interface{} {
	old := *nq
	n := len(old)
	item := old[n-1]
	*nq = old[0 : n-1]
	return item
}
