package request

type TableRequest struct {
	TableId int `json:"tableId" binding:"required"`
}
