// Package ps provides methods for managing app processes.
package ps

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
	"golang.org/x/net/websocket"
)

// List lists an app's processes.
func List(c *drycc.Client, appID string, results int) (api.PodsList, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/pods/", appID)
	body, count, reqErr := c.LimitedRequest(u, results)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.Pods{}, -1, reqErr
	}

	var procs []api.Pods
	if err := json.Unmarshal([]byte(body), &procs); err != nil {
		return []api.Pods{}, -1, err
	}

	return procs, count, reqErr
}

// Exec a command in a container.
func Exec(c *drycc.Client, appID, podID string, command api.Command) (*websocket.Conn, error) {
	scheme := "ws"
	if c.ControllerURL.Scheme == "https" {
		scheme = "wss"
	}
	path := fmt.Sprintf("v2/apps/%s/pods/%s/exec/", appID, podID)
	u := url.URL{Scheme: scheme, Host: c.ControllerURL.Host, Path: path}
	config, err := websocket.NewConfig(u.String(), c.ControllerURL.String())
	if err != nil {
		return nil, err
	}
	config.Header = http.Header{
		"User-Agent":           {c.UserAgent},
		"Authorization":        {"token " + c.Token},
		"X-Drycc-Builder-Auth": {c.HooksToken},
	}
	conn, err := websocket.DialConfig(config)
	if err != nil {
		return nil, err
	}
	websocket.JSON.Send(conn, command)
	return conn, nil
}

// Logs retrieves logs from an pod. The number of log lines fetched can be set by the lines
func Logs(c *drycc.Client, appID, podID string, request api.PodLogsRequest) (*websocket.Conn, error) {
	scheme := "ws"
	if c.ControllerURL.Scheme == "https" {
		scheme = "wss"
	}
	path := fmt.Sprintf("v2/apps/%s/pods/%s/logs/", appID, podID)
	endpoint := url.URL{Scheme: scheme, Host: c.ControllerURL.Host, Path: path}

	config, err := websocket.NewConfig(endpoint.String(), c.ControllerURL.String())
	if err != nil {
		return nil, err
	}
	config.Header = http.Header{
		"User-Agent":           {c.UserAgent},
		"Authorization":        {"token " + c.Token},
		"X-Drycc-Builder-Auth": {c.HooksToken},
	}
	conn, err := websocket.DialConfig(config)
	if err != nil {
		return nil, err
	}
	err = websocket.JSON.Send(conn, request)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Scale increases or decreases an app's processes. The processes are specified in the target argument,
// a key-value map, where the key is the process name and the value is the number of replicas
func Scale(c *drycc.Client, appID string, targets map[string]int) error {
	u := fmt.Sprintf("/v2/apps/%s/scale/", appID)

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
	u := fmt.Sprintf("/v2/apps/%s/pods/restart/", appID)
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

// Describe pod state
func Describe(c *drycc.Client, appID string, podID string, results int) (api.PodState, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/pods/%s/describe/", appID, podID)

	body, count, reqErr := c.LimitedRequest(u, results)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.PodState{}, -1, reqErr
	}

	var podState api.PodState
	if err := json.Unmarshal([]byte(body), &podState); err != nil {
		return api.PodState{}, -1, err
	}
	return podState, count, reqErr
}

// ByType organizes processes of an app by process type.
func ByType(processes api.PodsList) api.PodTypes {
	var pts api.PodTypes

	for _, process := range processes {
		exists := false
		// Is processtype for process already exists, append to it.
		for i, pt := range pts {
			if pt.Type == process.Type {
				exists = true
				pts[i].PodsList = append(pts[i].PodsList, process)
				break
			}
		}

		// Is processtype for process doesn't exist, create a new one
		if !exists {
			pts = append(pts, api.PodType{
				Type:     process.Type,
				PodsList: api.PodsList{process},
			})
		}
	}

	// Sort the pods alphabetically by name.
	for _, pt := range pts {
		sort.Sort(pt.PodsList)
	}

	// Sort ProcessTypes alphabetically by process name
	sort.Sort(pts)

	return pts
}
