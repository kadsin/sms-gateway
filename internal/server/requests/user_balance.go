package requests

type UserBalanceRequest struct {
	Email   string  `json:"email" validate:"required,email"`
	Balance float32 `json:"balance" validate:"required,numeric"`
}
