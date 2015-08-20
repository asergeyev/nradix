// Copyright (C) 2015 Alex Sergeyev
// This project is licensed under the terms of the MIT license.
// Read LICENSE file for information for all notices and permissions.

package nradix

import "errors"

type Node struct {
	left, right, parent *Node
	value               interface{}
}

type Tree struct {
	root   *Node
	free   *Node
	has128 bool
}

const startbit = uint32(0x80000000)

var (
	ErrNodeBusy = errors.New("NodeBusy")
	ErrNotFound = errors.New("NoSuchNode")
	ErrBadNode  = errors.New("BadNode")
)

func NewTree(preallocate int) *Tree {
	tree := new(Tree)
	tree.root = tree.newnode()
	if preallocate == 0 {
		return tree
	}

	// Simplification, static preallocate max 6 bits
	if preallocate > 6 || preallocate < 0 {
		preallocate = 6
	}

	var key, mask uint32

	for inc := startbit; preallocate > 0; inc, preallocate = inc>>1, preallocate-1 {
		key = 0
		mask >>= 1
		mask |= startbit

		for {
			tree.insert32(key, mask, nil)
			key += inc
			if key == 0 { // magic bits collide
				break
			}
		}
	}

	return tree
}

func (tree *Tree) insert32(key, mask uint32, value interface{}) error {
	bit := startbit
	node := tree.root
	next := tree.root
	for bit&mask != 0 {
		if key&bit != 0 {
			next = node.right
		} else {
			next = node.left
		}
		if next == nil {
			break
		}
		bit >>= 1
		node = next
	}
	if next != nil {
		if node.value != nil {
			return ErrNodeBusy
		}
		node.value = value
		return nil
	}
	for bit&mask != 0 {
		next = tree.newnode()
		next.parent = node
		if key&bit != 0 {
			node.right = next
		} else {
			node.left = next
		}
		bit >>= 1
		node = next
	}
	node.value = value

	return nil
}

func (tree *Tree) delete32(key, mask uint32) error {
	bit := startbit
	node := tree.root
	for node != nil && bit&mask != 0 {
		if key&bit != 0 {
			node = node.right
		} else {
			node = node.left
		}
		bit >>= 1
	}
	if node == nil {
		return ErrNotFound
	}

	if node.right != nil && node.left != nil {
		// keep it just trim value
		if node.value != nil {
			node.value = nil
			return nil
		}
		return ErrNotFound
	}

	// need to trim leaf
	for {
		if node.parent.right == node {
			node.parent.right = nil
		} else {
			node.parent.left = nil
		}
		// reserve this node for future use
		node.right = tree.free
		tree.free = node
		// move to parent, check if it's free of value and children
		node = node.parent
		if node.right != nil || node.left != nil || node.value != nil {
			break
		}
		// do not delete root node
		if node.parent == nil {
			break
		}
	}

	return nil
}

func (tree *Tree) find32(key, mask uint32) (value interface{}) {
	bit := startbit
	node := tree.root
	for node != nil {
		if node.value != nil {
			value = node.value
		}
		if key&bit != 0 {
			node = node.right
		} else {
			node = node.left
		}
		bit >>= 1
	}
	return value
}

func (tree *Tree) newnode() (p *Node) {
	if tree.free != nil {
		p = tree.free
		tree.free = tree.free.right
		return p
	}

	// ideally should be aligned in array but for now just let Go decide:
	return new(Node)
}
