package server

type Response struct {
	Status   int         `json:"status"`
	Response interface{} `json:"response"`
	Errors   []string    `json:"errors"`
}

type GetPayload struct {
	Must map[string]string `json:"match"`
}

type PostPayload struct {
	Set map[string]string
}

type PutPayload struct {
	Must map[string]string
	Set  map[string]string
}

type DeletePayload struct {
	Must map[string]string
}
