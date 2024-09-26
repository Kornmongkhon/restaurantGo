package enums

type StatusCode struct {
	Code    string
	Message string
}

var (
	Success = StatusCode{"S0000", "Success"}
	Invalid = StatusCode{"I0001", "Invalid request"}
	Error   = StatusCode{"E9999", "The system has a problem. Please contact the system administrator."}
)

func (s StatusCode) String() string {
	return s.Message
}

func (s StatusCode) GetCode() string {
	return s.Code
}
