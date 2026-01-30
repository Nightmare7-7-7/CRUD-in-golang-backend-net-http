package models

type User struct {
	Id        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	OldEmail  string `json:"old_email,omitempty"`
	NewEmail  string `json:"new_email,omitempty"`
}

type CreateUser struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type UpdateEmail struct {
	OldEmail string `json:"old_email,omitempty"`
	NewEmail string `json:"new_email,omitempty"`
}

type ReturnApiUser struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    []User `json:"data,omitempty"`
}
