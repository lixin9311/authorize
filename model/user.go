package model

// User is
type User struct {
	ID       uint64 `json:"id,omitempty"`
	Username string `json:"username" form:"username" query:"username"`
	Email    string `json:"email" form:"email" query:"email"`
	Password string `json:"password,omitempty" form:"password" query:"password"`
	Token    string `json:"token,omitempty"`
}
