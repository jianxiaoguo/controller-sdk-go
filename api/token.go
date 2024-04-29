package api

// Token is the structure of the token object.
type Token struct {
	UUID    string `json:"uuid"`
	Owner   string `json:"owner"`
	Alias   string `json:"alias"`
	Key     string `json:"fuzzy_key"`
	Created string `json:"created"`
	Updated string `json:"updated"`
}
