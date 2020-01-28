package dto

import "github.com/sfqi/library/domain/model"

type LoanResponse struct {
	ID            int    `json:"id"`
	TransactionID string `json:"transaction_id"`
	UserID        int    `json:"user_id"`
	BookID        int    `json:"book_id"`
	Type          string `json:"loan_type"`
}

func CreateLoanResponse(l *model.Loan, loanType string) *LoanResponse {
	return &LoanResponse{
		ID:            l.ID,
		TransactionID: l.TransactionID,
		UserID:        l.UserID,
		BookID:        l.BookID,
		Type:          loanType,
	}
}
