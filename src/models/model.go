package models

type Output struct {
	Status  int64       `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
type Query struct {
	PageIndex int `form:"page_index"`
	PageSize  int `form:"page_size"`
}
