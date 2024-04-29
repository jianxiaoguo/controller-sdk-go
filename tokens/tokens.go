package tokens

import (
	"encoding/json"
	"fmt"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// List tokens.
func List(c *drycc.Client, results int) ([]api.Token, int, error) {
	body, count, reqErr := c.LimitedRequest("/v2/tokens/", results)

	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.Token{}, -1, reqErr
	}

	var tokens []api.Token
	if err := json.Unmarshal([]byte(body), &tokens); err != nil {
		return []api.Token{}, -1, err
	}

	return tokens, count, reqErr
}

// Delete a token
func Delete(c *drycc.Client, id string) error {
	u := fmt.Sprintf("/v2/tokens/%s/", id)
	res, err := c.Request("DELETE", u, nil)
	if err == nil {
		res.Body.Close()
	}
	return err
}
