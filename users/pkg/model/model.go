package model

type (
	LoginResponse struct {
		Login   string    `json:"login"`
	}

	AuthResponse struct {
		Auth    string    `json:"token"`
	}

	User struct {
		Token      string    `json:"token"`
		Login      string    `json:"login"`
		Password   string    `json:"pswd"`
	}
)