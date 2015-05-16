package main

import (
	// "fmt"
	"strconv"
	"strings"
)

type item struct {
	name    string
	typ     string
	width   int
	offset  int
	pointer int
}

type SignalTable struct {
	elems   map[string]*item
	offset  int
	numtemp int
}

func NewItem(nm string, tp string, width int, ost int) *item {
	return &item{
		nm,
		tp,
		width,
		ost,
		0,
	}
}

func NewSignalTable() *SignalTable {
	return &SignalTable{
		make(map[string]*item),
		0,
		0,
	}
}

func (s *SignalTable) Enter(name string, tp string, width int) {
	sub := strings.Split(name, ",")
	for _, v := range sub {
		s.elems[v] = NewItem(v, tp, width, s.offset)
		s.offset += width
	}

}

func (s *SignalTable) loopUp(name string) string {
	if _, ok := s.elems[name]; ok {
		return name
	}
	return ""
}

func (s *SignalTable) newTemp(tp string, width int) string {
	name := "temp" + strconv.Itoa(s.numtemp)
	s.elems[name] = NewItem(name, tp, width, s.offset)
	s.offset += width
	return name
}
