// Copyright © 2024 Mark Summerfield. All rights reserved.
package sortedset

import "iter"

// Comparable allows only string or integer elements.
type Comparable interface {
	~string | ~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// SortedSet zero value is usable. Create with statements like these:
// var set SortedSet[string]
// set := SortedSet[int]{}
type SortedSet[E Comparable] struct {
	root *node[E]
	size int
}

type node[E Comparable] struct {
	element     E
	red         bool
	left, right *node[E]
}

// Insert inserts a new element into the Tree and returns true; or does
// nothing and returns false if the element is already present.
// For example:
//
//	ok := tree.Insert(element).
func (me *SortedSet[E]) Insert(element E) bool {
	inserted := false
	me.root, inserted = me.insert(me.root, element)
	me.root.red = false
	if inserted {
		me.size++
	}
	return inserted
}

func (me *SortedSet[E]) insert(root *node[E], element E) (*node[E], bool) {
	inserted := false
	if root == nil { // If element was in the tree it would go here
		return &node[E]{element: element, red: true}, true
	}
	if isRed(root.left) && isRed(root.right) {
		colorFlip(root)
	}
	if element < root.element {
		root.left, inserted = me.insert(root.left, element)
	} else if root.element < element {
		root.right, inserted = me.insert(root.right, element)
	}
	root = insertRotation(root)
	return root, inserted
}

func isRed[E Comparable](root *node[E]) bool {
	return root != nil && root.red
}

func colorFlip[E Comparable](root *node[E]) {
	root.red = !root.red
	if root.left != nil {
		root.left.red = !root.left.red
	}
	if root.right != nil {
		root.right.red = !root.right.red
	}
}

func insertRotation[E Comparable](root *node[E]) *node[E] {
	if isRed(root.right) && !isRed(root.left) {
		root = rotateLeft(root)
	}
	if isRed(root.left) && isRed(root.left.left) {
		root = rotateRight(root)
	}
	return root
}

func rotateLeft[E Comparable](root *node[E]) *node[E] {
	x := root.right
	root.right = x.left
	x.left = root
	x.red = root.red
	root.red = true
	return x
}

func rotateRight[E Comparable](root *node[E]) *node[E] {
	x := root.left
	root.left = x.right
	x.right = root
	x.red = root.red
	root.red = true
	return x
}

// Len returns the number of items in the tree.
func (me *SortedSet[E]) Len() int { return me.size }

// All returns a for .. range iterable of the set's elements, e.g.,
// for element := range tree.All()
func (me *SortedSet[E]) All() iter.Seq[E] {
	return func(yield func(E) bool) {
		all(me.root, yield)
	}
}

func all[E Comparable](root *node[E], yield func(E) bool) bool {
	if root != nil {
		return all(root.left, yield) &&
			yield(root.element) &&
			all(root.right, yield)
	}
	return true
}

// Contains returns true if the element is in the tree; otherwise false.
// For example:
//
//	ok := set.Contains(element).
func (me *SortedSet[E]) Contains(element E) bool {
	root := me.root
	for root != nil {
		if element < root.element {
			root = root.left
		} else if root.element < element {
			root = root.right
		} else {
			return true
		}
	}
	return false
}

// Delete deletes the element-value with the given element from the
// set and returns true, or does nothing and returns false if
// there is no element-value with the given element. For example:
//
//	deleted := set.Delete(element).
//
// See also [Clear]
func (me *SortedSet[E]) Delete(element E) bool {
	deleted := false
	if me.root != nil {
		if me.root, deleted = delete_(me.root, element); me.root != nil {
			me.root.red = false
		}
	}
	if deleted {
		me.size--
	}
	return deleted
}

func delete_[E Comparable](root *node[E], element E) (*node[E], bool) {
	deleted := false
	if element < root.element {
		if root.left != nil {
			if !isRed(root.left) && !isRed(root.left.left) {
				root = moveRedLeft(root)
			}
			root.left, deleted = delete_(root.left, element)
		}
	} else {
		if isRed(root.left) {
			root = rotateRight(root)
		}
		if element == root.element && root.right == nil {
			return nil, true
		}
		if root.right != nil {
			root, deleted = deleteRight(root, element)
		}
	}
	return fixUp(root), deleted
}

func moveRedLeft[E Comparable](root *node[E]) *node[E] {
	colorFlip(root)
	if root.right != nil && isRed(root.right.left) {
		root.right = rotateRight(root.right)
		root = rotateLeft(root)
		colorFlip(root)
	}
	return root
}

func deleteRight[E Comparable](root *node[E], element E) (*node[E], bool) {
	deleted := false
	if !isRed(root.right) && !isRed(root.right.left) {
		root = moveRedRight(root)
	}
	if element == root.element {
		smallest := first(root.right)
		root.element = smallest.element
		root.right = deleteMinimum(root.right)
		deleted = true
	} else {
		root.right, deleted = delete_(root.right, element)
	}
	return root, deleted
}

func moveRedRight[E Comparable](root *node[E]) *node[E] {
	colorFlip(root)
	if root.left != nil && isRed(root.left.left) {
		root = rotateRight(root)
		colorFlip(root)
	}
	return root
}

// We do not provide an exported First() method because this
// is an implementation detail.
func first[E Comparable](root *node[E]) *node[E] {
	for root.left != nil {
		root = root.left
	}
	return root
}

func deleteMinimum[E Comparable](root *node[E]) *node[E] {
	if root.left == nil {
		return nil
	}
	if !isRed(root.left) && !isRed(root.left.left) {
		root = moveRedLeft(root)
	}
	root.left = deleteMinimum(root.left)
	return fixUp(root)
}

func fixUp[E Comparable](root *node[E]) *node[E] {
	if isRed(root.right) {
		root = rotateLeft(root)
	}
	if isRed(root.left) && isRed(root.left.left) {
		root = rotateRight(root)
	}
	if isRed(root.left) && isRed(root.right) {
		colorFlip(root)
	}
	return root
}

// Clear deletes the entire tree.
// See also [Delete]
func (me *SortedSet[E]) Clear() {
	me.root = nil
	me.size = 0
}
