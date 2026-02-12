package dom

import "strconv"

func Bold(s string) string {
	return "<b>" + s + "</b>"
}

func Italic(s string) string {
	return "<i>" + s + "</i>"
}

func Br() string {
	return "<br/>"
}

func Hr() string {
	return "<hr/>"
}

func NewHeading(level int, text string, attr ...HtmlAttr) HtmlElement {
	return NewElement("h"+strconv.FormatInt(int64(level), 10), text, attr...)
}

func NewParagraph(text string, attr ...HtmlAttr) HtmlElement {
	return HtmlElement{
		TagName:       "p",
		Attr:          attr,
		InnerText:     text,
		IsForceEndTag: true,
	}
}

func NewDiv(text string, attr ...HtmlAttr) HtmlElement {
	return HtmlElement{
		TagName:       "div",
		Attr:          attr,
		InnerText:     text,
		IsForceEndTag: true,
	}
}

func NewImg(src, alt string, attr ...HtmlAttr) HtmlElement {
	elem := NewElement("img", "", HtmlAttr{Key: "src", Value: src})
	elem.IsEmpty = true
	if alt != "" {
		elem.AppendAttr(HtmlAttr{Key: "alt", Value: alt})
	}
	if len(attr) > 0 {
		elem.AppendAttr(attr...)
	}

	return elem
}

func NewSpan(text string, attr ...HtmlAttr) HtmlElement {
	return HtmlElement{
		TagName:       "span",
		Attr:          attr,
		InnerText:     text,
		IsForceEndTag: true,
	}
}
