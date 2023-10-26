package models

type Draft struct {
	ID        int                    `json:"id"`
	Type      string                 `json:"type"`
	Context   string                 `json:"context"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt string                 `json:"createdAt"`
	Iri       string                 `json:"iri"`
}
