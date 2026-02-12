package dom

func NewAttr(key, value string) HtmlAttr {
	return HtmlAttr{Key: key, Value: value}
}

func NewClass(c string) HtmlAttr {
	return NewAttr("class", c)
}

func NewId(id string) HtmlAttr {
	return NewAttr("id", id)
}

func NewHref(href string) HtmlAttr {
	return NewAttr("href", href)
}
