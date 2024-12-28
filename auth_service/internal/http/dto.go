package http

type LoginUserDto struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	AppID    int    `json:"app_id,omitempty"`
}

type RegisterDto struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}
