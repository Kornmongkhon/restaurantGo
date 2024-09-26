package request

type OrderRequest struct {
	TableId      int     `json:"table_id"`
	Total_Amount float64 `json:"total_amount"`
	Create_At    string  `json:"create_at"`
}
