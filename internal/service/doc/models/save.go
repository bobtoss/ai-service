package models

type SaveDoc struct {
	ID        string `json:"id"`
	RequestID string `json:"request_id"`
	Document  string `json:"document"`
	Name      string `json:"name"`
	CompanyId string `json:"company_id"`
	Priority  bool   `json:"priority"`
}

type SaveDocResponse struct {
	Status bool `json:"status"`
}
