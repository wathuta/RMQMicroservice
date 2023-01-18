package models

type Message struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}
