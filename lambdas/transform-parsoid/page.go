package main

type Page struct {
	Header     []*Content   `json:"header,omitempty"`
	Paragraphs []*Paragraph `json:"paragraphs,omitempty"`
	References []*Reference `json:"references,omitempty"`
}

func newPage() *Page {
	new := Page{nil, []*Paragraph{}, []*Reference{}}
	return &new
}

func (page *Page) AddParagraph(paragraph *Paragraph) {
	page.Paragraphs = append(page.Paragraphs, paragraph)
}

func (page *Page) SetHeader(header []*Content) {
	page.Header = header
}
