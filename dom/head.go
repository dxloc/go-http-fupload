package dom

type HtmlHead HtmlElement

func (p *HtmlDocument) Head() *HtmlHead {
	return (*HtmlHead)(&p.Children[0])
}

func (h *HtmlHead) AddElement(e ...HtmlElement) {
	(*HtmlElement)(h).AppendChild(e...)
}

func (h *HtmlHead) SetLanguage(s string) {
	h.Attr[0].Value = s
}

func (h *HtmlHead) SetCharset(s string) {
	h.Children[0].Attr[0].Value = s
}

func (h *HtmlHead) SetTitle(s string) {
	h.Children[2].InnerText = s
}
