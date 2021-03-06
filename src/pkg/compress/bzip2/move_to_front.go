// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bzip2

// moveToFrontDecoder implements a move-to-front list. Such a list is an
// efficient way to transform a string with repeating elements into one with
// many small valued numbers, which is suitable for entropy encoding. It works
// by starting with an initial list of symbols and references symbols by their
// index into that list. When a symbol is referenced, it's moved to the front
// of the list. Thus, a repeated symbol ends up being encoded with many zeros,
// as the symbol will be at the front of the list after the first access.
type moveToFrontDecoder struct {
	// Rather than actually keep the list in memory, the symbols are stored
	// as a circular, double linked list with the symbol indexed by head
	// at the front of the list.
	symbols []byte
	next    []uint8
	prev    []uint8
	head    uint8
}

// newMTFDecoder creates a move-to-front decoder with an explicit initial list
// of symbols.
func newMTFDecoder(symbols []byte) *moveToFrontDecoder {
	if len(symbols) > 256 {
		panic("too many symbols")
	}

	m := &moveToFrontDecoder{
		symbols: symbols,
		next:    make([]uint8, len(symbols)),
		prev:    make([]uint8, len(symbols)),
	}

	m.threadLinkedList()
	return m
}

// newMTFDecoderWithRange creates a move-to-front decoder with an initial
// symbol list of 0...n-1.
func newMTFDecoderWithRange(n int) *moveToFrontDecoder {
	if n > 256 {
		panic("newMTFDecoderWithRange: cannot have > 256 symbols")
	}

	m := &moveToFrontDecoder{
		symbols: make([]uint8, n),
		next:    make([]uint8, n),
		prev:    make([]uint8, n),
	}

	for i := 0; i < n; i++ {
		m.symbols[i] = byte(i)
	}

	m.threadLinkedList()
	return m
}

// threadLinkedList creates the initial linked-list pointers.
func (m *moveToFrontDecoder) threadLinkedList() {
	if len(m.symbols) == 0 {
		return
	}

	m.prev[0] = uint8(len(m.symbols) - 1)

	for i := 0; i < len(m.symbols)-1; i++ {
		m.next[i] = uint8(i + 1)
		m.prev[i+1] = uint8(i)
	}

	m.next[len(m.symbols)-1] = 0
}

func (m *moveToFrontDecoder) Decode(n int) (b byte) {
	// Most of the time, n will be zero so it's worth dealing with this
	// simple case.
	if n == 0 {
		return m.symbols[m.head]
	}

	i := m.head
	for j := 0; j < n; j++ {
		i = m.next[i]
	}
	b = m.symbols[i]

	m.next[m.prev[i]] = m.next[i]
	m.prev[m.next[i]] = m.prev[i]
	m.next[i] = m.head
	m.prev[i] = m.prev[m.head]
	m.next[m.prev[m.head]] = i
	m.prev[m.head] = i
	m.head = i

	return
}

// First returns the symbol at the front of the list.
func (m *moveToFrontDecoder) First() byte {
	return m.symbols[m.head]
}
