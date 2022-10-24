package modules

// ErrorModel ...
type ErrorModel struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
}

// FileUploadedModel ...
type FileUploadedModel struct {
	Filename string `json:"filename"`
	URL      string `json:"url"`
}
