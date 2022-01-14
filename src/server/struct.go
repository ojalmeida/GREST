package server

type Response struct {
	Status   int         `json:"status"`
	Response interface{} `json:"response"`
	Errors   []string    `json:"errors"`
}

type GetPayload struct {
	Must map[string]string `json:"must"`
}

type PostPayload struct {
	Set map[string]string `json:"set"`
}

type PutPayload struct {
	Must map[string]string `json:"must"`
	Set  map[string]string `json:"set"`
}

type DeletePayload struct {
	Must map[string]string `json:"must"`
}
