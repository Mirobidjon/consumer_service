package restservice

type Phone struct {
	ID     string `json:"id"`
	Number string `json:"number"`
}

type PhoneNumber struct {
	Number string `json:"number"`
}

type PhoneList struct {
	Phones []Phone `json:"phones"`
	Count  int32   `json:"count"`
}
