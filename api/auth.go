package api

// AuthLoginRequest represents the request structure for authentication login.
type AuthLoginRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// AuthLoginResponse represents the response structure for authentication login.
type AuthLoginResponse struct {
	Key string `json:"key,omitempty"`
}

// AuthTokenResponse is the definition of /v2/auth/login/.
type AuthTokenResponse struct {
	Token    string `json:"token"`
	Username string `json:"username,omitempty"`
}
