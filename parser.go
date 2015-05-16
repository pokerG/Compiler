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

type interCode struct {
	thraddrCodes []string
	current      int
}

var iC *interCode

func (c *interCode) genCode(code string) {
	c.thraddrCodes = append(c.thraddrCodes, code)
	c.current++
}

func (c *interCode) nextQuad() int {
	return c.current + 1
}
func (c *interCode) backPactch(p []int, i int) {
	for _, v := range p {
		sub := strings.Split(c.thraddrCodes[v], "-")
		c.thraddrCodes[v] = strings.Join(sub, strconv.Itoa(i))
	}
}
func (c *interCode) Print() {
	for k, v := range c.thraddrCodes {
		fmt.Println(k, ": ", v)
	}
}

type pair struct {
	state string
	sign  string
}

type attribute struct {
	signal    string
	input     string
	typ       string
	width     int
	addr      string
	truelist  []int
	falselist []int
	quad      int
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
	p.signal_stack.PushBack(&attribute{signal: "$"})
	p.action = make(map[pair]string)
	p.gotoo = make(map[pair]string)
	p.loadTable()
	// fmt.Println(p.action)
}

func (p *Parser) startParsing() {
	iC = &interCode{current: -1}
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
		} else if token_type == DIGITS {
			content = "digits"
		} else if token_type == REAL {
			content = "real"
		} else if token_type == FID {
			content = "fid"
		}
		if content == ";" {
			content = "semic"
		} else if content == "&&" {
			content = "and"
		} else if content == "||" {
			content = "or"
		} else if content == "!" {
			content = "not"
		}

		S := reflect.ValueOf(top.Value).String()
		np := pair{S, content}
		// fmt.Println(e.Value, np, p.action[np])
		if strings.HasPrefix(p.action[np], "shift") { //change to state i
			sub := strings.Split(p.action[np], " ")
			p.state_stack.PushBack(sub[1])
			in := reflect.ValueOf(e.Value).Elem().Field(1).String()
			attr := &attribute{signal: content, input: in}
			p.signal_stack.PushBack(attr)
			e = e.Next()
		} else if strings.HasPrefix(p.action[np], "reduce") { // A->B
			production := strings.TrimLeft(strings.TrimLeft(p.action[np], "reduce"), " ")
			fmt.Println(production)
			sub := strings.Split(p.action[np], "->")
			subb := strings.Split(strings.Trim(sub[1], " "), " ")
			var right []*list.Element

			if subb[0] == "e" {

			} else {
				for i := 0; i < len(subb); i = i + 1 {
					tmp1 := p.state_stack.Back()
					p.state_stack.Remove(tmp1)
					tmp2 := p.signal_stack.Back()
					right = append(right, tmp2)
					// fmt.Println(reflect.ValueOf(tmp2.Value).Elem().Field(0).String(), 1, reflect.ValueOf(tmp2.Value).Elem().Field(1).String(), 2, reflect.ValueOf(tmp2.Value).Elem().Field(2).String())
					p.signal_stack.Remove(tmp2)
				}
			}

			S = reflect.ValueOf(p.state_stack.Back().Value).String()
			subb = strings.Split(strings.Trim(sub[0], " "), " ")
			p.signal_stack.PushBack(&attribute{signal: subb[1]})
			np = pair{S, subb[1]}
			p.state_stack.PushBack(p.gotoo[np])
			p.Translate(production, right)
		} else if p.action[np] == "accept" {
			iC.Print()
			return
		} else {
			fmt.Println("error")
			return
		}
	}
}

func (p *Parser) Translate(production string, right []*list.Element) {
	switch production {
	case "Type -> int":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
			typ:    "int",
			width:  4,
		}
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "Type -> double":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
			typ:    "double",
			width:  8,
		}
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "Type -> string":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
			typ:    "string",
			width:  4,
		}
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "Type -> byte":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
			typ:    "byte",
			width:  1,
		}
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "Array -> [ digits ]":
		num, _ := strconv.Atoi(reflect.ValueOf(right[1].Value).Elem().Field(1).String())
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
			width:  num,
		}
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "Array -> e":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
			width:  1,
		}
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "muilType -> Array Type":
		width1 := reflect.ValueOf(right[0].Value).Elem().Field(3).Int()
		width2 := reflect.ValueOf(right[1].Value).Elem().Field(3).Int()
		ty := reflect.ValueOf(right[0].Value).Elem().Field(2).String()
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
			typ:    ty,
			width:  int(width1 * width2),
		}
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "Varlist -> id":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		name := reflect.ValueOf(right[0].Value).Elem().Field(1).String()
		attr := &attribute{
			signal: sg,
			input:  name,
		}
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "Varlist -> Varlist , id":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		name1 := reflect.ValueOf(right[0].Value).Elem().Field(1).String()
		name2 := reflect.ValueOf(right[2].Value).Elem().Field(1).String()
		attr := &attribute{
			signal: sg,
			input:  name2 + "," + name1,
		}
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "Define -> Varlist muilType":
		name1 := reflect.ValueOf(right[1].Value).Elem().Field(1).String()
		tp := reflect.ValueOf(right[0].Value).Elem().Field(2).String()
		width := reflect.ValueOf(right[0].Value).Elem().Field(3).Int()
		signaltable.Enter(name1, tp, int(width))
	case "declaration -> var Define semic":
	case "constant -> digits":
		in := reflect.ValueOf(right[0].Value).Elem().Field(1).String()
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
			typ:    "int",
			width:  4,
			addr:   in,
		}
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "constant -> real":
		in := reflect.ValueOf(right[0].Value).Elem().Field(1).String()
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
			typ:    "double",
			width:  8,
			addr:   in,
		}
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "F -> constant":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
			typ:    reflect.ValueOf(right[0].Value).Elem().Field(2).String(),
			width:  int(reflect.ValueOf(right[0].Value).Elem().Field(3).Int()),
			addr:   reflect.ValueOf(right[0].Value).Elem().Field(4).String(),
		}
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "F -> ( E )":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
			typ:    reflect.ValueOf(right[1].Value).Elem().Field(2).String(),
			width:  int(reflect.ValueOf(right[1].Value).Elem().Field(3).Int()),
			addr:   reflect.ValueOf(right[1].Value).Elem().Field(4).String(),
		}
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "T -> F":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
			typ:    reflect.ValueOf(right[0].Value).Elem().Field(2).String(),
			width:  int(reflect.ValueOf(right[0].Value).Elem().Field(3).Int()),
			addr:   reflect.ValueOf(right[0].Value).Elem().Field(4).String(),
		}
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "T -> T * F":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		addr2 := reflect.ValueOf(right[0].Value).Elem().Field(4).String()
		addr1 := reflect.ValueOf(right[2].Value).Elem().Field(4).String()
		attr := &attribute{
			signal: sg,
		}
		if reflect.ValueOf(right[0].Value).Elem().Field(3).Int() > reflect.ValueOf(right[2].Value).Elem().Field(3).Int() {
			attr.width = int(reflect.ValueOf(right[0].Value).Elem().Field(3).Int())
			attr.typ = reflect.ValueOf(right[0].Value).Elem().Field(2).String()
		} else {
			attr.width = int(reflect.ValueOf(right[2].Value).Elem().Field(3).Int())
			attr.typ = reflect.ValueOf(right[2].Value).Elem().Field(2).String()
		}
		attr.addr = signaltable.newTemp(attr.typ, attr.width)
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
		iC.genCode(attr.addr + " = " + addr1 + " * " + addr2)
	case "T => T / F":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		addr2 := reflect.ValueOf(right[0].Value).Elem().Field(4).String()
		addr1 := reflect.ValueOf(right[2].Value).Elem().Field(4).String()
		attr := &attribute{
			signal: sg,
		}
		if reflect.ValueOf(right[0].Value).Elem().Field(3).Int() > reflect.ValueOf(right[2].Value).Elem().Field(3).Int() {
			attr.width = int(reflect.ValueOf(right[0].Value).Elem().Field(3).Int())
			attr.typ = reflect.ValueOf(right[0].Value).Elem().Field(2).String()
		} else {
			attr.width = int(reflect.ValueOf(right[2].Value).Elem().Field(3).Int())
			attr.typ = reflect.ValueOf(right[2].Value).Elem().Field(2).String()
		}
		attr.addr = signaltable.newTemp(attr.typ, attr.width)
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
		iC.genCode(attr.addr + " = " + addr1 + " / " + addr2)
	case "E -> T":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
			typ:    reflect.ValueOf(right[0].Value).Elem().Field(2).String(),
			width:  int(reflect.ValueOf(right[0].Value).Elem().Field(3).Int()),
			addr:   reflect.ValueOf(right[0].Value).Elem().Field(4).String(),
		}
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "E -> E + T":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		addr2 := reflect.ValueOf(right[0].Value).Elem().Field(4).String()
		addr1 := reflect.ValueOf(right[2].Value).Elem().Field(4).String()
		attr := &attribute{
			signal: sg,
		}
		if reflect.ValueOf(right[0].Value).Elem().Field(3).Int() > reflect.ValueOf(right[2].Value).Elem().Field(3).Int() {
			attr.width = int(reflect.ValueOf(right[0].Value).Elem().Field(3).Int())
			attr.typ = reflect.ValueOf(right[0].Value).Elem().Field(2).String()
		} else {
			attr.width = int(reflect.ValueOf(right[2].Value).Elem().Field(3).Int())
			attr.typ = reflect.ValueOf(right[2].Value).Elem().Field(2).String()
		}
		attr.addr = signaltable.newTemp(attr.typ, attr.width)
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
		iC.genCode(attr.addr + " = " + addr1 + " + " + addr2)
	case "E -> E - T":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		addr2 := reflect.ValueOf(right[0].Value).Elem().Field(4).String()
		addr1 := reflect.ValueOf(right[2].Value).Elem().Field(4).String()
		attr := &attribute{
			signal: sg,
		}
		if reflect.ValueOf(right[0].Value).Elem().Field(3).Int() > reflect.ValueOf(right[2].Value).Elem().Field(3).Int() {
			attr.width = int(reflect.ValueOf(right[0].Value).Elem().Field(3).Int())
			attr.typ = reflect.ValueOf(right[0].Value).Elem().Field(2).String()
		} else {
			attr.width = int(reflect.ValueOf(right[2].Value).Elem().Field(3).Int())
			attr.typ = reflect.ValueOf(right[2].Value).Elem().Field(2).String()
		}
		attr.addr = signaltable.newTemp(attr.typ, attr.width)
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
		iC.genCode(attr.addr + " = " + addr1 + " - " + addr2)
	case "F -> lvalue":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
			typ:    reflect.ValueOf(right[0].Value).Elem().Field(2).String(),
			width:  int(reflect.ValueOf(right[0].Value).Elem().Field(3).Int()),
			addr:   reflect.ValueOf(right[0].Value).Elem().Field(4).String(),
		}
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "lvalue -> id":
		name := reflect.ValueOf(right[0].Value).Elem().Field(1).String()
		pt := signaltable.loopUp(name)
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		if pt != "" {
			attr := &attribute{
				signal: sg,
				typ:    signaltable.elems[pt].typ,
				width:  signaltable.elems[pt].width,
				addr:   name,
			}
			tmp := p.signal_stack.Back()
			p.signal_stack.Remove(tmp)
			p.signal_stack.PushBack(attr)
		} else {
			panic("No the variable!")
		}
	case "right -> E":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
			typ:    reflect.ValueOf(right[0].Value).Elem().Field(2).String(),
			width:  int(reflect.ValueOf(right[0].Value).Elem().Field(3).Int()),
			addr:   reflect.ValueOf(right[0].Value).Elem().Field(4).String(),
		}
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "right -> right , E":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
			addr:   reflect.ValueOf(right[2].Value).Elem().Field(4).String() + "," + reflect.ValueOf(right[0].Value).Elem().Field(4).String(),
		}
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "Assign -> left = right":
		addrr := reflect.ValueOf(right[0].Value).Elem().Field(4).String()
		addrl := reflect.ValueOf(right[2].Value).Elem().Field(4).String()
		ids := strings.Split(addrl, ",")
		for _, v := range ids {
			if signaltable.loopUp(v) == "" {
				panic("No the variable!")
			}
		}
		iC.genCode(addrl + " = " + addrr)
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
		}
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "left -> lvalue":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
			typ:    reflect.ValueOf(right[0].Value).Elem().Field(2).String(),
			width:  int(reflect.ValueOf(right[0].Value).Elem().Field(3).Int()),
			addr:   reflect.ValueOf(right[0].Value).Elem().Field(4).String(),
		}
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "left -> left , lvalue":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
			addr:   reflect.ValueOf(right[2].Value).Elem().Field(4).String() + "," + reflect.ValueOf(right[0].Value).Elem().Field(4).String(),
		}
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "lvalue -> id [ digits ]":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
			typ:    reflect.ValueOf(right[3].Value).Elem().Field(2).String(),
			width:  int(reflect.ValueOf(right[3].Value).Elem().Field(3).Int()),
			addr:   reflect.ValueOf(right[3].Value).Elem().Field(1).String() + "[" + reflect.ValueOf(right[1].Value).Elem().Field(1).String() + "]",
		}
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "lvalue -> id [ id ]":
		id1 := reflect.ValueOf(right[1].Value).Elem().Field(1).String()
		id2 := reflect.ValueOf(right[3].Value).Elem().Field(1).String()
		if signaltable.loopUp(id1) == "" || signaltable.loopUp(id2) == "" {
			panic("No the variable!")
		}
		if signaltable.elems[id2].typ != "int" {
			panic("subscript should be int")
		}
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
			typ:    reflect.ValueOf(right[3].Value).Elem().Field(2).String(),
			width:  int(reflect.ValueOf(right[3].Value).Elem().Field(3).Int()),
			addr:   reflect.ValueOf(right[3].Value).Elem().Field(1).String() + "[" + reflect.ValueOf(right[1].Value).Elem().Field(1).String() + "]",
		}
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)

	case "B -> ( B or M B )":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
		}
		b1fl := convert(reflect.ValueOf(right[4].Value).Elem().Field(6))
		b1tl := convert(reflect.ValueOf(right[4].Value).Elem().Field(5))
		b2tl := convert(reflect.ValueOf(right[1].Value).Elem().Field(5))
		b2fl := convert(reflect.ValueOf(right[1].Value).Elem().Field(6))
		quad := int(reflect.ValueOf(right[2].Value).Elem().Field(7).Int())
		iC.backPactch(b1fl, quad)
		attr.truelist = merge(b1tl, b2tl)
		attr.falselist = b2fl
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "B -> ( B and M B )":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
		}
		b1fl := convert(reflect.ValueOf(right[4].Value).Elem().Field(6))
		b1tl := convert(reflect.ValueOf(right[4].Value).Elem().Field(5))
		b2tl := convert(reflect.ValueOf(right[1].Value).Elem().Field(5))
		b2fl := convert(reflect.ValueOf(right[1].Value).Elem().Field(6))
		quad := int(reflect.ValueOf(right[2].Value).Elem().Field(7).Int())
		iC.backPactch(b1tl, quad)
		attr.truelist = b2tl
		attr.falselist = merge(b1fl, b2fl)
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "B -> not B":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
		}
		b1fl := convert(reflect.ValueOf(right[0].Value).Elem().Field(6))
		b1tl := convert(reflect.ValueOf(right[0].Value).Elem().Field(5))
		attr.truelist = b1fl
		attr.falselist = b1tl
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "B -> ( B )":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
		}
		b1fl := convert(reflect.ValueOf(right[1].Value).Elem().Field(6))
		b1tl := convert(reflect.ValueOf(right[1].Value).Elem().Field(5))
		attr.truelist = b1tl
		attr.falselist = b1fl
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "B -> E relop E":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
		}
		nextquad := iC.nextQuad()
		attr.truelist = make([]int, 1)
		attr.truelist[0] = nextquad
		attr.falselist = make([]int, 1)
		attr.falselist[0] = nextquad + 1
		addr1 := reflect.ValueOf(right[2].Value).Elem().Field(4).String()
		addr2 := reflect.ValueOf(right[0].Value).Elem().Field(4).String()
		relop := reflect.ValueOf(right[1].Value).Elem().Field(4).String()
		iC.genCode("if " + addr1 + " " + relop + " " + addr2 + " goto -")
		iC.genCode("goto -")
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "relop -> !=":
		fallthrough
	case "relop -> ==":
		fallthrough
	case "relop -> <=":
		fallthrough
	case "relop -> >=":
		fallthrough
	case "relop -> <":
		fallthrough
	case "relop -> >":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
			addr:   reflect.ValueOf(right[0].Value).Elem().Field(1).String(),
		}
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "B -> true":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
		}
		nextquad := iC.nextQuad()
		attr.truelist = make([]int, 1)
		attr.truelist[0] = nextquad
		iC.genCode("goto -")
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "B -> false":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
		}
		nextquad := iC.nextQuad()
		attr.falselist = make([]int, 1)
		attr.falselist[0] = nextquad
		iC.genCode("goto -")
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "M -> e":
		nextquad := iC.nextQuad()
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
			quad:   nextquad,
		}
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "senten -> if B { M sentens }":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
		}
		btl := convert(reflect.ValueOf(right[4].Value).Elem().Field(5))
		bfl := convert(reflect.ValueOf(right[4].Value).Elem().Field(6))
		s1l := convert(reflect.ValueOf(right[1].Value).Elem().Field(5))
		quad := int(reflect.ValueOf(right[2].Value).Elem().Field(7).Int())
		iC.backPactch(btl, quad)
		attr.truelist = merge(bfl, s1l)
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "senten -> if B { M sentens } N else { M sentens }":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
		}
		btl := convert(reflect.ValueOf(right[10].Value).Elem().Field(5))
		quad1 := int(reflect.ValueOf(right[8].Value).Elem().Field(7).Int())
		quad2 := int(reflect.ValueOf(right[2].Value).Elem().Field(7).Int())
		bfl := convert(reflect.ValueOf(right[10].Value).Elem().Field(6))
		s1l := convert(reflect.ValueOf(right[7].Value).Elem().Field(5))
		s2l := convert(reflect.ValueOf(right[1].Value).Elem().Field(5))
		nl := convert(reflect.ValueOf(right[5].Value).Elem().Field(5))
		iC.backPactch(btl, quad1)
		iC.backPactch(bfl, quad2)
		attr.truelist = merge(s1l, merge(nl, s2l))
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "N -> e":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
		}
		nextquad := iC.nextQuad()
		attr.truelist = make([]int, 1)
		attr.truelist[0] = nextquad
		iC.genCode("goto -")
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "senten -> Assign semic":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
		}
		attr.truelist = nil
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "sentens -> sentens M senten":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
		}
		quad := int(reflect.ValueOf(right[1].Value).Elem().Field(7).Int())
		ll := convert(reflect.ValueOf(right[2].Value).Elem().Field(5))
		sl := convert(reflect.ValueOf(right[0].Value).Elem().Field(5))
		iC.backPactch(ll, quad)
		attr.truelist = sl
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "sentens -> e":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
		}
		attr.truelist = nil
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	case "senten -> while M B { M sentens }":
		sg := reflect.ValueOf(p.signal_stack.Back().Value).Elem().Field(0).String()
		attr := &attribute{
			signal: sg,
		}
		s1l := convert(reflect.ValueOf(right[1].Value).Elem().Field(5))
		quad2 := int(reflect.ValueOf(right[2].Value).Elem().Field(7).Int())
		quad1 := int(reflect.ValueOf(right[5].Value).Elem().Field(7).Int())
		btl := convert(reflect.ValueOf(right[4].Value).Elem().Field(5))
		bfl := convert(reflect.ValueOf(right[4].Value).Elem().Field(6))
		fmt.Println(s1l, btl, quad1, quad2)
		iC.backPactch(s1l, quad1)
		iC.backPactch(btl, quad2)
		attr.truelist = bfl
		iC.genCode("goto " + strconv.Itoa(quad1))
		tmp := p.signal_stack.Back()
		p.signal_stack.Remove(tmp)
		p.signal_stack.PushBack(attr)
	}

}

func convert(arr reflect.Value) []int {
	var a []int
	a = make([]int, arr.Len())
	for i := 0; i < arr.Len(); i++ {
		a[i] = int(arr.Index(i).Int())
	}
	return a
}

func merge(a, b []int) []int {
	var c []int
	c = make([]int, len(a)+len(b))
	copy(c, a)
	copy(c[len(a):], b)
	return c
}
