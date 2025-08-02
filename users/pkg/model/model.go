package model

type (
	RegisterResponse struct {
		Login   string    `json:"login"`
	}

	AuthResponse struct {
		Token   string    `json:"token"`
	}

	User struct {
		Token            string
		Login            string
		HashedPassword   string
	}
)
