// Package ps provides methods for managing app processes.
package pts

import (
	"encoding/json"
	"fmt"
	"sort"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// List lists an app's processes.
func List(c *drycc.Client, appID string, results int) (api.Ptypes, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/ptypes/", appID)
	body, count, reqErr := c.LimitedRequest(u, results)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.Ptype{}, -1, reqErr
	}

	var ptypes []api.Ptype
	if err := json.Unmarshal([]byte(body), &ptypes); err != nil {
		return []api.Ptype{}, -1, err
	}

	return ptypes, count, reqErr
}

// Describe Ptype state
func Describe(c *drycc.Client, appID string, ptype string, results int) (api.PtypeStates, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/ptypes/%s-%s/describe/", appID, appID, ptype)

	body, count, reqErr := c.LimitedRequest(u, results)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.PtypeStates{}, -1, reqErr
	}

	var PtypeStates api.PtypeStates
	if err := json.Unmarshal([]byte(body), &PtypeStates); err != nil {
		return api.PtypeStates{}, -1, err
	}
	return PtypeStates, count, reqErr
}

// Scale increases or decreases an app's processes. The processes are specified in the target argument,
// a key-value map, where the key is the process name and the value is the number of replicas
func Scale(c *drycc.Client, appID string, targets map[string]int) error {
	u := fmt.Sprintf("/v2/apps/%s/ptypes/scale/", appID)

	body, err := json.Marshal(targets)

	if err != nil {
		return err
	}

	res, err := c.Request("POST", u, body)
	if err != nil && !drycc.IsErrAPIMismatch(err) {
		return err
	}
	defer res.Body.Close()
	return err
}

// Restart restarts an app's processes. To restart all app processes, pass empty strings for
// procType and name. To restart an specific process, pass an procType by leave name empty.
// To restart a specific instance, pass a procType and a name.
func Restart(c *drycc.Client, appID string, targets map[string]string) error {
	u := fmt.Sprintf("/v2/apps/%s/ptypes/restart/", appID)
	body, err := json.Marshal(targets)
	if err != nil {
		return err
	}
	res, err := c.Request("POST", u, body)
	if err != nil && !drycc.IsErrAPIMismatch(err) {
		return err
	}
	defer res.Body.Close()
	return err
}

// ByType organizes process types of an app by process type.
func ByType(ptypes api.Ptypes) api.Ptypes {
	// Sort ProcessTypes alphabetically by process types name
	sort.Sort(ptypes)

	return ptypes
}
