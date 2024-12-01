package structs

// All response sturcts
// Response struct have Res prefix

type ResPlaceOrder struct {
	Id        string  `json:"id" validate:"required"`
	Symbol    string  `json:"symbol" validate:"required"`
	Volume    float64 `json:"volume" validate:"required"`
	OrderType string  `json:"order_type" validate:"required"`
	Price     string  `json:"price" validate:"required"`
	UserId    string  `json:"user_id" validate:"required"`
	CreatedAt string  `json:"created_at" validate:"required"`
	UpdatedAt string  `json:"updated_at" validate:"required"`
}
