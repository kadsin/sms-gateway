package requests

type SmsRequest struct {
	ClientId      string `json:"client_id" validate:"required,uuid"`
	ReceiverPhone string `json:"receiver_phone" validate:"required,e164"`
	Content       string `json:"content" validate:"required,max=160"`
	IsExpress     *bool  `json:"is_express" validate:"required,boolean"`
}
