package requests

type SmsRequest struct {
	ReceiverPhone string `json:"receiver_phone" validate:"required,e164"`
	Content       string `json:"content" validate:"required,max=160"`
	IsExpress     *bool  `json:"is_express" validate:"required,boolean"`
}
