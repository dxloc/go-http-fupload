package dom

func NewUnorderedList(l []string, a ...HtmlAttr) HtmlElement {
	elem := NewElement("ul", "", a...)
	for i := range l {
		elem.AppendChild(NewElement("li", l[i]))
	}
	return elem
}

func NewOrderedList(l []string, a ...HtmlAttr) HtmlElement {
	elem := NewElement("ol", "", a...)
	for i := range l {
		elem.AppendChild(NewElement("li", l[i]))
	}
	return elem
}
