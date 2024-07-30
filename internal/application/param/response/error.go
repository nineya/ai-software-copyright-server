package response

type ErrorResponse struct {
	Status  int    `json:"status"`
	Title   string `json:"title"`
	Message string `json:"message"`
}
