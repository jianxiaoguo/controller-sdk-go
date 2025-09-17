// Package limits provides methods for managing resource limits of apps.
package limits

import (
	"encoding/json"
	"fmt"
	"strings"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// Specs is list all available limit specs
func Specs(c *drycc.Client, keywords string, results int) ([]api.LimitSpec, int, error) {
	u := "/v2/limits/specs/"
	if keywords != "" {
		u += fmt.Sprintf("?keywords=%s", keywords)
	}

	body, count, reqErr := c.LimitedRequest(u, results)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.LimitSpec{}, -1, reqErr
	}
	var limitSpecs []api.LimitSpec
	if err := json.Unmarshal([]byte(body), &limitSpecs); err != nil {
		return []api.LimitSpec{}, -1, err
	}
	return limitSpecs, count, reqErr
}

// Plans is list all available limit plans
func Plans(c *drycc.Client, specID string, cpu, memory, results int) ([]api.LimitPlan, int, error) {
	var queryArray []string
	if cpu > 0 {
		queryArray = append(queryArray, fmt.Sprintf("cpu=%d", cpu))
	}
	if memory > 0 {
		queryArray = append(queryArray, fmt.Sprintf("memory=%d", memory))
	}
	if specID != "" {
		queryArray = append(queryArray, fmt.Sprintf("spec-id=%s", specID))
	}
	u := "/v2/limits/plans/"
	if len(queryArray) > 0 {
		u = fmt.Sprintf("%s?%s", u, strings.Join(queryArray, "&"))
	}

	body, count, reqErr := c.LimitedRequest(u, results)

	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.LimitPlan{}, -1, reqErr
	}
	var limitPlans []api.LimitPlan
	if err := json.Unmarshal([]byte(body), &limitPlans); err != nil {
		return []api.LimitPlan{}, -1, err
	}
	return limitPlans, count, reqErr
}

// GetPlan is get a available Plan
func GetPlan(c *drycc.Client, planID string) (api.LimitPlan, error) {
	u := fmt.Sprintf("/v2/limits/plans/%s/", planID)
	res, reqErr := c.Request("GET", u, nil)

	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.LimitPlan{}, reqErr
	}
	defer res.Body.Close()
	limitPlan := api.LimitPlan{}
	if err := json.NewDecoder(res.Body).Decode(&limitPlan); err != nil {
		return api.LimitPlan{}, err
	}
	return limitPlan, reqErr
}
