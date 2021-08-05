// Package auth handles user management: creation, deletion, and authentication.
package auth

import (
	"encoding/json"
	"fmt"
	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
	"net/http"
)

// Login to the controller and get a oauth url
func Login(c *drycc.Client) (string, error) {
	c.HTTPClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	res, reqErr := c.Request("POST", "/v2/auth/login/", nil)

	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return "", reqErr
	}
	defer res.Body.Close()

	URL := res.Header.Get("Location")
	return URL, reqErr
}

// Token to the controller and get a token
func Token(c *drycc.Client, key string) (api.AuthLoginResponse, error) {
	path := fmt.Sprintf("/v2/auth/token/%s/", key)
	res, reqErr := c.Request("GET", path, nil)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.AuthLoginResponse{}, reqErr
	}
	defer res.Body.Close()
	token := api.AuthLoginResponse{}
	if err := json.NewDecoder(res.Body).Decode(&token); err != nil {
		return api.AuthLoginResponse{}, err
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
