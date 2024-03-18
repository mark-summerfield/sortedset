// Copyright © 2024 Mark Summerfield. All rights reserved.
package rbset

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestAPI(t *testing.T) {
	//                 012345678
	letters := []rune("ZENZEBRAS")
	set := RbSet[rune]{}
	for _, letter := range letters {
		set.Insert(letter)
	}
	var out strings.Builder
	for letter := range set.All() {
		out.WriteString(fmt.Sprintf("%c ", letter))
	}
	text := strings.TrimSpace(out.String())
	expected := "A B E N R S Z"
	if expected != text {
		t.Errorf("expected %q; got %q", expected, text)
	}
	size := set.Len()
	ok := set.Contains('Y')
	if ok {
		t.Error("expected false; got true")
	}
	ok = set.Contains('B')
	if !ok {
		t.Errorf("expected true; got %t", ok)
	}
	size = set.Len()
	deleted := set.Delete('Y')
	if deleted {
		t.Error("expected false; got true")
	}
	deleted = set.Delete('B')
	if !deleted {
		t.Error("expected true; got false")
	}
	size = set.Len()
	if size != 6 {
		t.Errorf("expected 6; got %d", size)
	}
	set.Clear()
	size = set.Len()
	if size != 0 {
		t.Errorf("expected 0; got %d", size)
	}
}

func Test1(t *testing.T) {
	data := []string{"can", "in", "a", "ebony", "go", "be", "dent", "for"}
	expected := []string{"a", "be", "can", "dent", "ebony", "for", "go",
		"in"}
	var set RbSet[string]
	for _, datum := range data {
		set.Insert(datum)
	}
	i := 1
	for word := range set.All() {
		if word != expected[i-1] {
			t.Errorf("expected %q %d; got %q", expected[i-1], i, word)
		}
		i++
	}
}

func Test2(t *testing.T) {
	data := []int{3, 8, 1, 5, 7, 2, 4, 6, 8, 5, 2, 7}
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8}
	set := RbSet[int]{}
	for _, datum := range data {
		set.Insert(datum)
	}
	i := 1
	for element := range set.All() {
		if element != expected[i-1] {
			t.Errorf("expected %d %d; got %d", expected[i-1], i, element)
		}
		i++
	}
}

func TestStringKeyInsertion(t *testing.T) {
	var wordSet RbSet[string]
	for _, word := range []string{"one", "Two", "THREE", "four", "Five"} {
		wordSet.Insert(strings.ToLower(word))
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
	var intSet RbSet[int]
	for _, number := range []int{9, 1, 8, 2, 7, 3, 6, 4, 5, 0} {
		intSet.Insert(number)
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
	var intSet RbSet[int]
	for _, number := range []int{9, 1, 8, 2, 7, 3, 6, 4, 5, 0} {
		intSet.Insert(number)
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
	var intSet RbSet[int]
	intSet.Insert(7)
	passTree(&intSet, t)
}

func passTree(set *RbSet[int], t *testing.T) {
	for _, number := range []int{9, 3, 6, 4, 5, 0} {
		set.Insert(number)
	}
	if set.Len() != 7 {
		t.Errorf("should have %d items", 7)
	}
}

// Thanks to Russ Cox for improving these benchmarks
func BenchmarkContainsSuccess(b *testing.B) {
	b.StopTimer() // Don't time creation and population
	var intSet RbSet[int]
	for i := range 1000000 {
		intSet.Insert(i)
	}
	b.StartTimer() // Time the Contains() method succeeding
	for i := range b.N {
		intSet.Contains(i % 1e6)
	}
}

func BenchmarkContainsFailure(b *testing.B) {
	b.StopTimer() // Don't time creation and population
	intSet := RbSet[int]{}
	for i := range 1000000 {
		intSet.Insert(2 * i)
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

func BenchmarkRbSetInsertion(b *testing.B) {
	var m RbSet[int]
	for i := range 1000000 {
		m.Insert(i)
	}
}

func BenchmarkRbSetIteration(b *testing.B) {
	b.StopTimer() // Don't time creation and population
	var m RbSet[int]
	for i := range 1000000 {
		m.Insert(i)
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
