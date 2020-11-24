package model

type Setting struct {
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
	IsDeleted bool        `json:"is_deleted"`
}
