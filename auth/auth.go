// Package auth handles user management: creation, deletion, and authentication.
package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// Login to the controller and get a oauth url
func Login(c *drycc.Client, username, password string) (string, error) {
	c.HTTPClient.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}
	var err error
	var body []byte
	if username != "" && password != "" {
		body, err = json.Marshal(api.AuthLoginRequest{Username: username, Password: password})
		if err != nil {
			return "", err
		}
	} else {
		body = nil
	}

	res, err := c.Request("POST", "/v2/auth/login/", body)
	if err != nil && !drycc.IsErrAPIMismatch(err) {
		return "", err
	}
	defer res.Body.Close()
	if username != "" && password != "" {
		login := api.AuthLoginResponse{}
		if err := json.NewDecoder(res.Body).Decode(&login); err != nil {
			return "", err
		}
		return login.Key, nil
	}

	url := res.Header.Get("Location")
	return url, err
}

// Token to the controller and get a token
func Token(c *drycc.Client, key, alias string) (api.AuthTokenResponse, error) {
	path := fmt.Sprintf("/v2/auth/token/%s/?alias=%s", key, alias)
	res, reqErr := c.Request("GET", path, nil)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.AuthTokenResponse{}, reqErr
	}
	defer res.Body.Close()
	token := api.AuthTokenResponse{}
	if err := json.NewDecoder(res.Body).Decode(&token); err != nil {
		return api.AuthTokenResponse{}, err
	}
	return token, reqErr
}

// Whoami retrives the user object for the authenticated user.
func Whoami(c *drycc.Client) (api.User, error) {
	res, err := c.Request("GET", "/v2/auth/whoami/", nil)
	if err != nil {
		return api.User{}, err
	}
	defer res.Body.Close()

	resUser := api.User{}
	if err = json.NewDecoder(res.Body).Decode(&resUser); err != nil {
		return api.User{}, err
	}

	return resUser, nil
}
