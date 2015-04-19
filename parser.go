package main

import (
	// "bytes"
	"container/list"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	. "github.com/pokerG/Compiler/common"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type pair struct {
	state string
	sign  string
}
type Parser struct {
	token_stream list.List
	parse_tree   list.List
	state_stack  list.List
	signal_stack list.List
	action       map[pair]string
	gotoo        map[pair]string
	dict         []string
	token_index  int
}

func (p *Parser) createParser(token_stream list.List) {
	p.token_stream = token_stream
	token := &Token{}
	token.token_type = ENDSIGNAL
	token.content = "#"
	p.token_stream.PushBack(token)
	p.token_index = 0
	p.state_stack.PushBack("0")
	p.signal_stack.PushBack("$")
	p.action = make(map[pair]string)
	p.gotoo = make(map[pair]string)
	p.loadTable()
	// fmt.Println(p.action)
}

func (p *Parser) startParsing() {
	e := p.token_stream.Front()
	// for e := p.token_stream.Front(); e != nil; e = e.Next() {
	for {
		top := p.state_stack.Back()
		token_type := reflect.ValueOf(e.Value).Elem().Field(0).Int()
		content := reflect.ValueOf(e.Value).Elem().Field(1).String()
		if token_type == IDENTIFIER {
			content = "id"
		} else if token_type == END_OF_FILE {
			content = "$"
		} else if token_type == NUMBER || token_type == STRING || token_type == CHARACTER {
			content = "constant"
		}
		if content == ";" {
			content = "semic"
		}
		S := reflect.ValueOf(top.Value).String()
		np := pair{S, content}
		fmt.Println(e.Value, np, p.action[np])
		if strings.HasPrefix(p.action[np], "shift") { //change to state i
			sub := strings.Split(p.action[np], " ")
			p.state_stack.PushBack(sub[1])
			p.signal_stack.PushBack(content)
			e = e.Next()
		} else if strings.HasPrefix(p.action[np], "reduce") { // A->B
			fmt.Println(strings.TrimLeft(p.action[np], "reduce "))

			sub := strings.Split(p.action[np], "->")
			subb := strings.Split(strings.Trim(sub[1], " "), " ")
			for i := 0; i < len(subb); i = i + 1 {
				tmp := p.state_stack.Back()
				p.state_stack.Remove(tmp)
				tmp = p.signal_stack.Back()
				p.signal_stack.Remove(tmp)
			}
			S = reflect.ValueOf(p.state_stack.Back().Value).String()
			subb = strings.Split(strings.Trim(sub[0], " "), " ")
			p.signal_stack.PushBack(subb[1])
			np = pair{S, subb[1]}
			p.state_stack.PushBack(p.gotoo[np])

		} else if p.action[np] == "accept" {
			return
		} else {
			fmt.Println("error")
			return
		}
	}
}

func (p *Parser) loadTable() {
	fn, _ := os.Open("LALR分析表.htm")
	doc, _ := goquery.NewDocumentFromReader(fn)

	// fmt.Println(doc, err)
	var term, noterm int
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			sc := s.Find("td")
			sterm, _ := sc.Next().Attr("colspan")
			snoterm, _ := sc.Last().Attr("colspan")
			term, _ = strconv.Atoi(sterm)
			noterm, _ = strconv.Atoi(snoterm)
			// fmt.Println(term, noterm)
		} else if i == 1 {
			sc := s.Find("td")
			scc := sc.First()
			// fmt.Println(scc.Text())
			for j := 0; j < term+noterm; j++ {
				// fmt.Println(scc.Text())
				p.dict = append(p.dict, scc.Text())
				sc = sc.Next()
				scc = sc.First()
			}
			// fmt.Println(p.dict)
		} else {
			sc := s.Find("td")
			scc := sc.First()
			state := scc.Text()
			for j := 1; j <= term; j++ {
				sc = sc.Next()
				scc = sc.First()
				np := pair{state, p.dict[j-1]}
				p.action[np] = scc.Text()
			}
			for j := 0; j < noterm; j++ {
				sc = sc.Next()
				scc = sc.First()
				np := pair{state, p.dict[term+j]}
				p.gotoo[np] = scc.Text()
			}
		}

	})
}
