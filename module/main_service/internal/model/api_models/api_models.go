package api_models

type APIResponse struct {
	ResponseInfo `json:"Response"`
}

type ResponseInfo struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
	Data    string `json:"Data"`
}
