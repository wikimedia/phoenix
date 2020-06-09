package common

// ChangeEvent is a JSON object representing a change to a document
type ChangeEvent struct {
	Title      string `json:"title"`
	ServerName string `json:"server_name"`
	Revision   int    `json:"revision"`
}
