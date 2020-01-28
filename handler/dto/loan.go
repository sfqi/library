package dto

type LoanResponse struct {
	ID            int    `json:"id"`
	TransactionID string `json:"transaction_id"`
	UserID        int    `json:"user_id"`
	BookID        int    `json:"book_id"`
	Type          string `json:"loan_type"`
}
