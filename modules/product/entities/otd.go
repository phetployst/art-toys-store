package entities

type (
	ProductResponse struct {
		ID          uint    `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		ImageURL    string  `json:"image_url"`
	}

	CountProduct struct {
		Count int `json:"count" validate:"required,gte=1"`
	}
)
