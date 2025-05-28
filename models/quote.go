package models

// Quote represents a quote object
// swagger:model
type Quote struct {
	ID     int    `json:"id"`
	Author string `json:"author"`
	Text   string `json:"quote"`
}
