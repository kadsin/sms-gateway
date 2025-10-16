package requests

type ReportsRequest struct {
	ClientId string `json:"client_id" validate:"required,uuid"`
}
