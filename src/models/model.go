package models

type Resp struct {
	Status  int64       `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
type Req struct {
	PageIndex int `json:"page_index"`
	PageSize  int `json:"page_size"`
}
