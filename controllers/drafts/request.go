package drafts

type CreateDraftRequest struct {
	Type    string                 `json:"type"`
	Context string                 `json:"context"`
	Data    map[string]interface{} `json:"data"`
}

type UpdateDraftRequest struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}
