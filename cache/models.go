package cache

type Session struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	IsBlocked bool   `json:"is_blocked"`
}
