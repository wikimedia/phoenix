package main

type Paragraph struct {
	Name    string     `json:"name"`
	Content []*Content `json:"content"`
}

func newParagraph(name string, content []*Content) *Paragraph {
	return &Paragraph{name, content}
}
