package models

type UpdatePriority struct {
	RequestID string `json:"request_id"`
	Priority  bool   `json:"priority"`
}
