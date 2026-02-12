package dom

type HtmlDocument HtmlElement

const (
	DefaultTitle   string = "Go HTTP file server"
	DefaultCharset string = "utf-8"
	DefaultLang    string = "en"
)

func NewDocument(title, lang, charset string) *HtmlDocument {
	if title == "" {
		title = DefaultTitle
	}
	if charset == "" {
		charset = DefaultCharset
	}
	if lang == "" {
		lang = DefaultLang
	}
	return &HtmlDocument{
		TagName: "html",
		Attr:    []HtmlAttr{{Key: "lang", Value: lang}},
		Children: []HtmlElement{
			{
				TagName: "head",
				Children: []HtmlElement{
					{TagName: "meta", Attr: []HtmlAttr{{Key: "charset", Value: charset}}},
					{TagName: "meta", Attr: []HtmlAttr{
						{Key: "name", Value: "viewport"},
						{Key: "content", Value: "width=device-width, initial-scale=1.0"},
					}},
					{TagName: "title", InnerText: title},
				},
			},
			{TagName: "body"},
		},
	}
}

func (p HtmlDocument) Serialize() string {
	elem := HtmlElement(p)
	return "<!DOCTYPE html>" + elem.String()
}
