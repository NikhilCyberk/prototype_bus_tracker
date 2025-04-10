// models/message.go
package models

type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}
