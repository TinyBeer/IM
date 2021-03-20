package view

import (
	"fmt"
	"log"
)

// menu page
type Page struct {
	head        string
	description string // description about the page
	parent      *Page
	count       int
	options     map[int]*Option
}

func NewPage(head string, description string, parent *Page) *Page {
	return &Page{
		head:        head,
		description: description,
		parent:      parent,
		count:       0,
		options:     make(map[int]*Option, DEFAULT_OPTION_MAP_CAP),
	}
}

func (p *Page) SelectOption(index int) {
	if index < 0 || index > p.count {
		log.Printf("SelectOption 参数错误")
		return
	} else if p.options[index].method == nil {
		return
	} else {
		p.options[index].method()
	}

}

func (p *Page) AddOption(content string, method func()) {
	p.count++
	p.options[p.count] = NewOption(content, method)
}

func (p *Page) GetParent() *Page {
	if p.parent == nil {
		log.Println("没有上级页面了")
		return p
	}
	return p.parent
}

func (p *Page) SetDescription(description string) {
	p.description = description
}

func (p *Page) Show() {

	fmt.Println(p.head)

	if p.description != "" {
		fmt.Println(p.description)
	}

	for i := 1; i <= p.count; i++ {
		fmt.Printf("%d.%s\n", i, p.options[i].content)
	}
	fmt.Println("请选择:")
}
