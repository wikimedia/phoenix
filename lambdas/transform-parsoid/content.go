package main

// import "golang.org/x/net/html"

// // Content represents a piece of content on a page
// type Content struct {
// 	Schema   string     `json:"schema,omitempty"`
// 	Rel      string     `json:"rel,omitempty"`
// 	TypeOf   string     `json:"typeOf,omitempty"`
// 	Text     string     `json:"text,omitempty"`
// 	DataMW   string     `json:"attributes,omitempty"`
// 	Parent   *Content   `json:"-"`
// 	Children []*Content `json:"content,omitempty"`
// }

// // AddChild to content
// func (content *Content) AddChild(child *Content) {
// 	content.Children = append(content.Children, child)
// }

// //AddAttribtues to content
// func (content *Content) AddAttribtues(token html.Token) {
// 	attributes := getAttribtues(token)
// 	content.Rel = attributes["rel"]
// 	content.TypeOf = attributes["typeof"]
// 	content.DataMW = attributes["data-mw"]
// }

// func newContent(tagName string, parent *Content) *Content {
// 	new := Content{tagName, "", "", "", "", parent, []*Content{}}
// 	if parent != nil {
// 		parent.AddChild(&new)
// 	}
// 	return &new
// }
