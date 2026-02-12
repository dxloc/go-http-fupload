package dom

type HtmlBody HtmlElement

func (p *HtmlDocument) Body() *HtmlBody {
	return (*HtmlBody)(&p.Children[1])
}

func (b *HtmlBody) AddElement(elem ...HtmlElement) {
	(*HtmlElement)(b).AppendChild(elem...)
}
