package response

type CustomResponse struct {
	Code    string `json:"code"`           // Status code, e.g., "S0000", "E9999"
	Message string `json:"message"`        // Message describing the status
	Data    any    `json:"data,omitempty"` // Data field, will be omitted if empty
}
