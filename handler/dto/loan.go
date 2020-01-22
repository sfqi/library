package dto

type LoanResponse struct {
	ID            int    `json:"id"`
	TransactionID string `json:"transaction_id"`
	UserID        int    `json:"user_id"`
	BookID        int    `json:"book_id"`
	Type          string `json:"loan_type"`
}

func CreateLoanResponse(id int, transactionID string, userID int, bookID int, loanType string) (*LoanResponse, error) {
	return &LoanResponse{
		ID:            id,
		TransactionID: transactionID,
		UserID:        userID,
		BookID:        bookID,
		Type:          loanType,
	}, nil
}
