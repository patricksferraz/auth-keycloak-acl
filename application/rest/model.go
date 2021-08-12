package rest

import "time"

type Auth struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AccessToken struct {
	AccessToken string `json:"access_token" binding:"required"`
}

type JWT struct {
	AccessToken      string `json:"access_token,omitempty"`
	IDToken          string `json:"id_token,omitempty"`
	ExpiresIn        int    `json:"expires_in,omitempty"`
	RefreshExpiresIn int    `json:"refresh_expires_in,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
	NotBeforePolicy  int    `json:"not_before_policy,omitempty"`
	SessionState     string `json:"session_state,omitempty"`
	Scope            string `json:"scope,omitempty"`
}

type Claims struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles,omitempty"`
}

type HTTPError struct {
	Code  int    `json:"code,omitempty" example:"400"`
	Error string `json:"error,omitempty" example:"status bad request"`
}

type CreateUserRequest struct {
	Username   string `json:"username" binding:"required"`
	EmployeeID string `json:"employee_id" binding:"required"`
}

type CreateUserResponse struct {
	ID string `json:"id"`
}

type IDRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type SetPasswordRequest struct {
	Password  string `json:"password" binding:"required"`
	Temporary bool   `json:"temporary"`
}

type HTTPResponse struct {
	Code    int    `json:"code,omitempty" example:"200"`
	Message string `json:"message,omitempty" example:"a message"`
}

type Filter struct {
	Username *string `json:"username" form:"username"`
	Enabled  *bool   `json:"enabled" form:"enabled"`
	PageSize *int    `json:"page_size" form:"page_size" default:"10"`
	Page     *int    `json:"page" form:"page" default:"0"`
}

type SearchUsersRequest struct {
	Filter `json:",inline"`
}

type User struct {
	ID         string    `json:"id"`
	Username   string    `json:"username,omitempty"`
	Enabled    bool      `json:"enabled,omitempty"`
	EmployeeID string    `json:"employee_id,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}
