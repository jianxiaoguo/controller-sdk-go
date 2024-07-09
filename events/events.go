// Package ps provides methods for managing app processes.
package events

import (
	"encoding/json"
	"fmt"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// List events of an app process.
func ListPodEvents(c *drycc.Client, appID string, podName string, results int) (api.AppEvents, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/events/?pod_name=%s", appID, podName)
	body, count, reqErr := c.LimitedRequest(u, results)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.AppEvent{}, -1, reqErr
	}

	var events []api.AppEvent
	if err := json.Unmarshal([]byte(body), &events); err != nil {
		return []api.AppEvent{}, -1, err
	}

	return events, count, reqErr
}

// List events of an app ptype.
func ListPtypeEvents(c *drycc.Client, appID string, ptype string, results int) (api.AppEvents, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/events/?ptype_name=%s-%s", appID, appID, ptype)
	body, count, reqErr := c.LimitedRequest(u, results)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.AppEvent{}, -1, reqErr
	}

	var events []api.AppEvent
	if err := json.Unmarshal([]byte(body), &events); err != nil {
		return []api.AppEvent{}, -1, err
	}

	return events, count, reqErr
}
