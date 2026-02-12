package dom

/*
NewScript(src, integrity, crossorigin, referrerpolicy)
*/
func NewScript(a ...string) HtmlElement {
	var attr []HtmlAttr
	if len(a) > 0 {
		attr = append(attr, NewAttr("src", a[0]))
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
		TagName:       "script",
		Attr:          attr,
		IsForceEndTag: true,
	}
}

func NewScriptRaw(script string) HtmlElement {
	if CanMinify() {
		script, _ = Minify(TypeTextJavascript, script)
	}
	return HtmlElement{
		TagName:       "script",
		InnerText:     script,
		IsForceEndTag: true,
	}
}
