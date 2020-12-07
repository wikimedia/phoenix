package common

// ChangeEvent is a JSON object representing a change to a document
type ChangeEvent struct {
	Title      string `json:"title"`
	ServerName string `json:"server_name"`
	Revision   int    `json:"revision"`
}

// NodeStoredEvent is a JSON object that corresponds to a Node being added to the Content Store.
type NodeStoredEvent struct {
	ID string `json:"id"`
}
