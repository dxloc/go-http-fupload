package dom

import "strings"

type HtmlAttr struct {
	Key   string
	Value string
}

type HtmlElement struct {
	TagName       string
	Attr          []HtmlAttr
	Children      []HtmlElement
	InnerText     string
	IsEmpty       bool
	IsForceEndTag bool
}

func NewElement(tag string, innerText string, a ...HtmlAttr) HtmlElement {
	return HtmlElement{
		TagName:   tag,
		Attr:      a,
		InnerText: innerText,
	}
}

func (e HtmlElement) String() string {
	var s strings.Builder
	s.WriteString("<" + e.TagName)
	if len(e.Attr) > 0 {
		for _, a := range e.Attr {
			s.WriteString(" " + a.Key + "=\"" + a.Value + "\"")
		}
	}
	if len(e.Children) != 0 || len(e.InnerText) != 0 || e.IsForceEndTag {
		s.WriteString(">")
	} else {
		s.WriteString("/>")
		return s.String()
	}

	if !e.IsEmpty {
		for _, c := range e.Children {
			s.WriteString(c.String())
		}
		s.WriteString(e.InnerText + "</" + e.TagName + ">")
	}

	return s.String()
}

func (e *HtmlElement) AppendAttr(a ...HtmlAttr) {
	e.Attr = append(e.Attr, a...)
}

func (e *HtmlElement) SetAttr(a ...HtmlAttr) {
	e.Attr = a
}

func (e *HtmlElement) AppendChild(c ...HtmlElement) {
	e.Children = append(e.Children, c...)
}

func (e *HtmlElement) SetChild(c ...HtmlElement) {
	e.Children = c
}

func (e *HtmlElement) SetText(s string) {
	e.InnerText = s
}
