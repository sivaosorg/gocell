package models

// Define table by Gin tag
type Common struct {
	Id     int  `json:"_id" db:"id"`
	Submit bool `json:"submit" db:"submit"`
}

type CommonRequest struct {
}

type CommonResponse struct {
}
