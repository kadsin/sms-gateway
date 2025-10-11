package requests

type ExmpleRequest struct {
	Ex string `json:"example" validate:"required, min=3, max=50"`
}
