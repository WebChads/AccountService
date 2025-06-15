package dtos

// Response represents common API response
// swagger:model Response
type Response struct {
	Status  int `json:"status" example:"200"`
	Message any `json:"message"`
}
