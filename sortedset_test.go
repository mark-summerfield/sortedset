// Copyright © 2024-25 Mark Summerfield. All rights reserved.
package sortedset

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestAPI(t *testing.T) {
	//                 012345678
	letters := []rune("ZENZEBRAS")
	sset := SortedSet[rune]{}
	for i := 0; i < len(letters); i++ {
		sset.Add(letters[i])
	}
	var out strings.Builder
	for letter := range sset.All() {
		out.WriteString(fmt.Sprintf("%c ", letter))
	}
	text := strings.TrimSpace(out.String())
	expected := "A B E N R S Z"
	if expected != text {
		t.Errorf("expected %q; got %q", expected, text)
	}
	ok := sset.Contains('Y')
	if ok {
		t.Error("expected false; got true")
	}
	ok = sset.Contains('B')
	if !ok {
		t.Errorf("expected true; got %t", ok)
	}
	deleted := sset.Delete('Y')
	if deleted {
		t.Error("expected false; got true")
	}
	deleted = sset.Delete('B')
	if !deleted {
		t.Error("expected true; got false")
	}
	size := sset.Len()
	if size != 6 || sset.IsEmpty() {
		t.Errorf("expected 6; got %d", size)
	}
	sset.Clear()
	size = sset.Len()
	if size != 0 || !sset.IsEmpty() {
		t.Errorf("expected 0; got %d", size)
	}
}

func Test1(t *testing.T) {
	data := []string{"can", "in", "a", "ebony", "go", "be", "dent", "for"}
	expected := []string{
		"a", "be", "can", "dent", "ebony", "for", "go",
		"in",
	}
	var sset SortedSet[string]
	for _, datum := range data {
		sset.Add(datum)
	}
	i := 1
	for word := range sset.All() {
		if word != expected[i-1] {
			t.Errorf("expected %q %d; got %q", expected[i-1], i, word)
		}
		i++
	}
}

func Test2(t *testing.T) {
	data := []int{3, 8, 1, 5, 7, 2, 4, 6, 8, 5, 2, 7}
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8}
	sset := SortedSet[int]{}
	for _, datum := range data {
		sset.Add(datum)
	}
	i := 1
	for element := range sset.All() {
		if element != expected[i-1] {
			t.Errorf("expected %d %d; got %d", expected[i-1], i, element)
		}
		i++
	}
	for element := range sset.All() {
		if element == 1 {
			break
		}
	}
}

func TestStringKeyInsertion(t *testing.T) {
	var wordSet SortedSet[string]
	for _, word := range []string{"one", "Two", "THREE", "four", "Five"} {
		wordSet.Add(strings.ToLower(word))
	}
	var words []string
	for word := range wordSet.All() {
		words = append(words, word)
	}
	actual, expected := strings.Join(words, ""), "fivefouronethreetwo"
	if actual != expected {
		t.Errorf("%q != %q", actual, expected)
	}
}

func TestIntKeyContains(t *testing.T) {
	var intSet SortedSet[int]
	for _, number := range []int{9, 1, 8, 2, 7, 3, 6, 4, 5, 0} {
		intSet.Add(number)
	}
	for _, number := range []int{0, 1, 5, 8, 9} {
		if ok := intSet.Contains(number); !ok {
			t.Errorf("failed to find %d", number)
		}
	}
	for _, number := range []int{-1, -21, 10, 11, 148} {
		if ok := intSet.Contains(number); ok {
			t.Errorf("should not have found %d", number)
		}
	}
}

func TestIntKeyDelete(t *testing.T) {
	var intSet SortedSet[int]
	for _, number := range []int{9, 1, 8, 2, 7, 3, 6, 4, 5, 0} {
		intSet.Add(number)
	}
	if intSet.Len() != 10 {
		t.Errorf("set len %d should be 10", intSet.Len())
	}
	length := 9
	for i, number := range []int{0, 1, 5, 8, 9} {
		if deleted := intSet.Delete(number); !deleted {
			t.Errorf("failed to delete %d", number)
		}
		if intSet.Len() != length-i {
			t.Errorf("map len %d should be %d", intSet.Len(), length-i)
		}
	}
	for _, number := range []int{-1, -21, 10, 11, 148} {
		if deleted := intSet.Delete(number); deleted {
			t.Errorf("should not have deleted nonexistent %d", number)
		}
	}
	if intSet.Len() != 5 {
		t.Errorf("map len %d should be 5", intSet.Len())
	}
}

func TestPassing(t *testing.T) {
	var intSet SortedSet[int]
	intSet.Add(7)
	passTree(&intSet, t)
}

func passTree(sset *SortedSet[int], t *testing.T) {
	for _, number := range []int{9, 3, 6, 4, 5, 0} {
		sset.Add(number)
	}
	if sset.Len() != 7 {
		t.Errorf("should have %d items", 7)
	}
}

// Thanks to Russ Cox for improving these benchmarks
func BenchmarkContainsSuccess(b *testing.B) {
	b.StopTimer() // Don't time creation and population
	var intSet SortedSet[int]
	for i := range 1000000 {
		intSet.Add(i)
	}
	b.StartTimer() // Time the Contains() method succeeding
	for i := range b.N {
		intSet.Contains(i % 1e6)
	}
}

func BenchmarkContainsFailure(b *testing.B) {
	b.StopTimer() // Don't time creation and population
	intSet := SortedSet[int]{}
	for i := range 1000000 {
		intSet.Add(2 * i)
	}
	b.StartTimer() // Time the Contains() method failing
	for i := range b.N {
		intSet.Contains(2*(i%1e6) + 1)
	}
}

func BenchmarkMapInsertion(b *testing.B) {
	m := map[int]int{}
	for i := range 1000000 {
		m[i] = i
	}
}

func BenchmarkMapSortedIteration(b *testing.B) {
	b.StopTimer() // Don't time creation and population
	m := map[int]int{}
	for i := range 1000000 {
		m[i] = i
	}
	b.StartTimer() // Time sort & iterate
	total := 0
	keys := make([]int, 0, len(m))
	for element := range m {
		keys = append(keys, element)
	}
	slices.Sort(keys)
	for _, element := range keys {
		total += element
	}
	b.StopTimer() // Don't time check
	if total != 499999500000 {
		panic(total)
	}
}

func BenchmarkSortedSetInsertion(b *testing.B) {
	var m SortedSet[int]
	for i := range 1000000 {
		m.Add(i)
	}
}

func BenchmarkSortedSetIteration(b *testing.B) {
	b.StopTimer() // Don't time creation and population
	var m SortedSet[int]
	for i := range 1000000 {
		m.Add(i)
	}
	b.StartTimer() // Time sort & iterate
	total := 0
	for element := range m.All() {
		total += element
	}
	b.StopTimer() // Don't time check
	if total != 499999500000 {
		panic(total)
	}
}

func TestDifference(t *testing.T) {
	s := New(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	u := New(2, 4, 6, 8)
	d := s.Difference(u)
	check(d.String(), d.Len(), "{0 1 3 5 7 9}", 6, t)
	d = u.Difference(s)
	check(d.String(), d.Len(), "{}", 0, t)
}

func TestSymmetricDifference(t *testing.T) {
	s := New(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	u := New(2, 4, 6, 8)
	d := s.SymmetricDifference(u)
	e := u.SymmetricDifference(s)
	if !d.Equal(e) {
		t.Errorf("unexpected unequal: d=%v e=%v", d, e)
	}
	check(d.String(), d.Len(), "{0 1 3 5 7 9}", 6, t)
	d = u.SymmetricDifference(s)
	e = u.SymmetricDifference(s)
	if !d.Equal(e) {
		t.Errorf("unexpected unequal: d=%v e=%v", d, e)
	}
	check(d.String(), d.Len(), "{0 1 3 5 7 9}", 6, t)
}

func TestIntersection(t *testing.T) {
	s := New(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	u := New(2, 4, 6, 8)
	x := s.Intersection(u)
	a := u.Intersection(s)
	if !x.Equal(a) {
		t.Errorf("unexpected unequal: %v != %v", x, a)
	}
	check(x.String(), x.Len(), "{2 4 6 8}", 4, t)
	v := New(1, 3, 5)
	y := u.Intersection(v)
	b := v.Intersection(u)
	if !y.Equal(b) {
		t.Errorf("unexpected unequal: %v != %v", y, b)
	}
	check(y.String(), y.Len(), "{}", 0, t)
	z := v.Intersection(u)
	c := u.Intersection(v)
	if !z.Equal(c) {
		t.Errorf("unexpected unequal: %v != %v", z, c)
	}
	check(z.String(), z.Len(), "{}", 0, t)
}

func TestUnion(t *testing.T) {
	s := New(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	u := New(2, 4, 6, 8, 10, 12)
	x := s.Union(u)
	check(x.String(), x.Len(), "{0 1 2 3 4 5 6 7 8 9 10 12}", 12, t)
}

func TestUnite(t *testing.T) {
	s := New(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	s.Unite(New(2, 4, 6, 8, 10, 12))
	check(s.String(), s.Len(), "{0 1 2 3 4 5 6 7 8 9 10 12}", 12, t)
}

func TestIsDisjoint(t *testing.T) {
	s := New(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	u := s.Clone()
	if s.IsDisjoint(u) {
		t.Error("unexpectedly disjoint")
	}
	if u.IsDisjoint(s) {
		t.Error("unexpectedly disjoint")
	}
	w := New(10, 11, 12)
	if !u.IsDisjoint(w) {
		t.Error("unexpectedly not disjoint")
	}
	if !w.IsDisjoint(u) {
		t.Error("unexpectedly not disjoint")
	}
}

func TestIsSubsetOf(t *testing.T) {
	s := New(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	u := s.Clone()
	if !s.IsSubsetOf(u) {
		t.Error("unexpectedly not subset")
	}
	w := New(10, 11, 12)
	if w.IsSubsetOf(s) {
		t.Error("unexpectedly a subset")
	}
	x := New(4, 6, 2)
	if !x.IsSubsetOf(s) {
		t.Error("unexpectedly not subset")
	}
}

func TestIsSupersetOf(t *testing.T) {
	s := New(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	u := s.Clone()
	if !s.IsSupersetOf(u) {
		t.Error("unexpectedly not superset")
	}
	w := New(10, 11, 12)
	if w.IsSupersetOf(s) {
		t.Error("unexpectedly a superset")
	}
	x := New(4, 6, 2)
	if x.IsSupersetOf(s) {
		t.Error("unexpectedly a superset")
	}
	if !s.IsSupersetOf(x) {
		t.Error("unexpectedly not a superset")
	}
}

func TestEqual(t *testing.T) {
	s := New(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	u := s.Clone()
	if !s.Equal(u) {
		t.Errorf("%v != %v", s, u)
	}
	u.Add(-3)
	if s.Equal(u) {
		t.Errorf("%v == %v", s, u)
	}
}

func TestToSlice(t *testing.T) {
	s := New(19, 21, 1, 2, 4, 8)
	u := s.ToSlice()
	check(fmt.Sprintf("%v", u), len(u), "[1 2 4 8 19 21]", s.Len(), t)
}

func TestAll(t *testing.T) {
	s := New(10, 20, 30, 40, 50, 60, 70, 80, 90)
	n := 0
	for v := range s.All() {
		n += v
	}
	if n != 450 {
		t.Errorf("expected 450, got %d", n)
	}
}

func TestAllX(t *testing.T) {
	s := New(10, 20, 30, 40, 50, 60, 70, 80, 90)
	n := 0
	for i, v := range s.AllX() {
		n += v + i
	}
	if n != 486 {
		t.Errorf("expected 486, got %d", n)
	}
	n = 0
	for i, v := range s.AllX(1) {
		n += v + i
	}
	if n != 495 {
		t.Errorf("expected 495, got %d", n)
	}
}

func check(act string, actSize int, exp string, expSize int, t *testing.T) {
	if actSize != expSize {
		t.Errorf("expected %d elements, got %d", expSize, actSize)
	}
	if exp != act {
		t.Errorf("expected %s, got %s", exp, act)
	}
}
