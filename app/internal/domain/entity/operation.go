package entity

type Operation struct {
	UUID         string  `json:"uuid"`
	CategoryUUID string  `json:"category_uuid"`
	Header       string  `json:"header"`
	Description  string  `json:"description"`
	MoneySum     float64 `json:"money_sum"`
}
