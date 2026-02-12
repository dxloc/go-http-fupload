package dom

/*
NewStyleSheet(href, integrity, crossorigin, referrerpolicy)
*/
func NewStyleSheet(a ...string) HtmlElement {
	attr := []HtmlAttr{
		{Key: "rel", Value: "stylesheet"},
	}
	if len(a) > 0 {
		attr = append(attr, NewAttr("href", a[0]))
	}
	if len(a) > 1 {
		attr = append(attr, NewAttr("integrity", a[1]))
	}
	if len(a) > 2 {
		attr = append(attr, NewAttr("crossorigin", a[2]))
	}
	if len(a) > 3 {
		attr = append(attr, NewAttr("referrerpolicy", a[3]))
	}
	return HtmlElement{
		TagName: "link",
		Attr:    attr,
	}
}

func NewStyleSheetRaw(css string) HtmlElement {
	if CanMinify() {
		css, _ = Minify(TypeTextCss, css)
	}
	return NewElement("style", css)
}
