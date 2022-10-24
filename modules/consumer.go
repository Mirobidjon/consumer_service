package modules

type Consumer struct {
	RecordID string `json:"record_id"`
}

// Response ...
type Response struct {
	Status      string      `json:"status"`
	Description string      `json:"description"`
	Data        interface{} `json:"data"`
}
