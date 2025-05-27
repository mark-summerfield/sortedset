// Copyright © 2024-25 Mark Summerfield. All rights reserved.

// ([TOC]) A sorted set based on a red-black tree.
//
// See also https://pkg.go.dev/github.com/mark-summerfield/set
//
// [TOC]: file:///home/mark/app/golib/doc/index.html
package sortedset

import (
	"fmt"
	"iter"
	"strings"

	"github.com/mark-summerfield/unum"
)

type Comparable = unum.Comparable

// SortedSet zero value is usable. Create with statements like these:
//
//	var set SortedSet[string]
//	set := SortedSet[int]{}
//
// or use [New]:
//
//	set := New(1, 2, 4)
type SortedSet[E Comparable] struct {
	root *node[E]
	size int
}

// New returns a new SortedSet containing the given elements (if any).
// If no elements are given, the type must be specified since it can't be
// inferred.
func New[E Comparable](elements ...E) SortedSet[E] {
	sset := SortedSet[E]{}
	for _, element := range elements {
		sset.Add(element)
	}
	return sset
}

type node[E Comparable] struct {
	element     E
	red         bool
	left, right *node[E]
}

// Add adds a new element into the SortedSet and returns true; or does
// nothing and returns false if the element is already present.
// For example:
//
//	ok := sset.Add(element).
func (me *SortedSet[E]) Add(element E) bool {
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
	if root == nil { // If element was in the SortedSet it would go here
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

// Len returns the number of items in the SortedSet.
func (me *SortedSet[E]) Len() int { return me.size }

// All returns a for .. range iterable of the SortedSet's elements, e.g.,
// for element := range sset.All()
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

// AllX returns an iterator, e.g.,
// for count, element := range sset.AllX(1) ...
func (me *SortedSet[E]) AllX(start ...int) iter.Seq2[int, E] {
	return func(yield func(int, E) bool) {
		i := 0
		if len(start) > 0 {
			i = start[0]
		}
		for key := range me.All() {
			if !yield(i, key) {
				return
			}
			i++
		}
	}
}

// Contains returns true if the element is in the SortedSet; otherwise
// false. For example:
//
//	ok := sset.Contains(element).
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

// Delete deletes the given element from the SortedSet and returns true, or
// does nothing and returns false if the element is not in the SortedSet.
// For example:
//
//	deleted := sset.Delete(element).
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

// Clear deletes all the elements in the SortedSet.
// See also [Delete].
func (me *SortedSet[E]) Clear() {
	me.root = nil
	me.size = 0
}

// IsEmpty returns true if there are no elements in the set; otherwise
// returns false.
func (me *SortedSet[E]) IsEmpty() bool { return me.size == 0 }

// Difference returns a new SortedSet that contains the elements which are
// in this SortedSet that are not in the other SortedSet.
func (me *SortedSet[E]) Difference(other SortedSet[E]) SortedSet[E] {
	diff := New[E]()
	for element := range me.All() {
		if !other.Contains(element) {
			diff.Add(element)
		}
	}
	return diff
}

// SymmetricDifference returns a new SortedSet that contains the elements
// which are in this SortedSet or the other SortedSet—but not in both
// SortedSets.
func (me *SortedSet[E]) SymmetricDifference(other SortedSet[E]) SortedSet[E] {
	diff := New[E]()
	for element := range me.All() {
		if !other.Contains(element) {
			diff.Add(element)
		}
	}
	for element := range other.All() {
		if !me.Contains(element) {
			diff.Add(element)
		}
	}
	return diff
}

// Intersection returns a new SortedSet that contains the elements this
// SortedSet has in common with the other SortedSet.
func (me *SortedSet[E]) Intersection(other SortedSet[E]) SortedSet[E] {
	intersection := New[E]()
	for element := range me.All() {
		if other.Contains(element) {
			intersection.Add(element)
		}
	}
	return intersection
}

// Union returns a new SortedSet that contains the elements from this
// SortedSet and from the other SortedSet (with no duplicates of course).
// See also [SortedSet.Unite].
func (me *SortedSet[E]) Union(other SortedSet[E]) SortedSet[E] {
	union := me.Clone()
	union.Unite(other)
	return union
}

// Unite adds all the elements from other that aren't already in this
// SortedSet to this SortedSet.
// See also [SortedSet.Union].
func (me *SortedSet[E]) Unite(other SortedSet[E]) {
	for element := range other.All() {
		me.Add(element)
	}
}

// IsDisjoint returns true if this SortedSet has no elements in common with
// the other SortedSet; otherwise returns false.
func (me *SortedSet[E]) IsDisjoint(other SortedSet[E]) bool {
	for element := range me.All() {
		if other.Contains(element) {
			return false
		}
	}
	return true
}

// IsSubsetOf returns true if this SortedSet is a subset of the other
// SortedSet, i.e., if every member of this SortedSet is in the other
// SortedSet; otherwise returns false.
func (me *SortedSet[E]) IsSubsetOf(other SortedSet[E]) bool {
	for element := range me.All() {
		if !other.Contains(element) {
			return false
		}
	}
	return true
}

// IsSupersetOf returns true if this SortedSet is a superset of the other
// SortedSet, i.e., if every member of the other SortedSet is in this
// SortedSet; otherwise returns false.
func (me SortedSet[E]) IsSupersetOf(other SortedSet[E]) bool {
	return other.IsSubsetOf(me)
}

// Equal returns true if this SortedSet has the same elements as the other
// SortedSet; otherwise returns false.
func (me *SortedSet[E]) Equal(other SortedSet[E]) bool {
	if me.Len() != other.Len() {
		return false
	}
	for element := range me.All() {
		if !other.Contains(element) {
			return false
		}
	}
	return true
}

// Clone returns a copy of this SortedSet.
func (me *SortedSet[E]) Clone() SortedSet[E] {
	clone := SortedSet[E]{}
	for element := range me.All() {
		clone.Add(element)
	}
	return clone
}

// ToSlice returns this SortedSet's elements as a sorted slice.
// For iteration either use this, or if you only need one value at a time,
// use [All] or [AllX].
func (me *SortedSet[E]) ToSlice() []E {
	slice := make([]E, 0, me.Len())
	for element := range me.All() {
		slice = append(slice, element)
	}
	return slice
}

// String returns a human readable string representation of the SortedSet.
func (me *SortedSet[E]) String() string {
	format := "%s%v"
	if me.hasStringElements() {
		format = "%s%q"
	}
	var out strings.Builder
	out.WriteByte('{')
	sep := ""
	for _, element := range me.ToSlice() {
		fmt.Fprintf(&out, format, sep, element)
		sep = " "
	}
	out.WriteByte('}')
	return out.String()
}

func (me *SortedSet[E]) hasStringElements() bool {
	for element := range me.All() {
		_, ok := any(element).(string)
		return ok
	}
	return false
}
