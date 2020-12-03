package main

type LoginRequiredParams struct {
	UserID string `json:"u"`
	Jwt    string `json:"j"`
	Data   interface{}
}

type errResult struct {
	Msg string
}
