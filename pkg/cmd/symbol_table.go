package cmd

import (
	"bytes"
	"fmt"
)

// A basic implementation of a hash map with linked lists as buckets

// MapSize - number of buckets in the map, representing unique hash values or
// in other words - heads of the linked lists buckets.
const MapSize = 5

// Node - a unit in the linked list containing a link to a next unit
type Node struct {
	Key   string
	Value string
	Next  *Node
}

type SymbolTable struct {
	Data []*Node
}

// NewSymbolTable - constructor for SymbolTable
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{Data: make([]*Node, MapSize)}
}

// String prints a single node
func (n *Node) String() string {
	return fmt.Sprintf("[%s, %s]\n", n.Key, n.Value)
}

// String prints a whole symbol table content as well as nulls
func (h *SymbolTable) String() string {
	var output bytes.Buffer
	fmt.Fprintln(&output, "{")
	for _, n := range h.Data {
		if n != nil {
			fmt.Fprintf(&output, "\t[%s: %s]\n", n.Key, n.Value)
			for node := n.Next; node != nil; node = node.Next {
				fmt.Fprintf(&output, "\t[%s: %s]\n", node.Key, node.Value)
			}
		}
	}

	fmt.Fprintln(&output, "}")

	return output.String()
}

// Insert inserts a variable by a given variable name (key) and its value
func (h *SymbolTable) Insert(key string, value string) {
	index := getIndex(key)

	if h.Data[index] == nil {
		// index is empty, go ahead and insert
		h.Data[index] = &Node{Key: key, Value: value}
	} else {
		// there is a collision, get into linked-list mode
		startingNode := h.Data[index]
		for ; startingNode.Next != nil; startingNode = startingNode.Next {
			if startingNode.Key == key {
				// the key exists, it's a modifying operation
				startingNode.Value = value
				return
			}
		}
		startingNode.Next = &Node{Key: key, Value: value}
	}
}

// Get returns the value by a given key if it exists
func (h *SymbolTable) Get(key string) (string, bool) {
	index := getIndex(key)
	if h.Data[index] != nil {
		startingNode := h.Data[index]
		for ; ; startingNode = startingNode.Next {
			if startingNode.Key == key {
				// key matched
				return startingNode.Value, true
			}

			if startingNode.Next == nil {
				break
			}
		}
	}

	return "", false
}

// hash calculates hash of the key by Jenkins' hash function algorithm
func hash(key string) uint8 {
	var h uint8
	for _, ch := range key {
		h += uint8(ch)
		h += h << 10
		h ^= h >> 6
	}

	h += h << 3
	h ^= h >> 11
	h += h << 15

	return h
}

// getIndex() - returns index of a bucket modding the hash of the key by the number of buckets
func getIndex(key string) (index int) {
	return int(hash(key)) % MapSize
}
