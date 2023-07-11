// Package apps provides methods for managing drycc apps.
package apps

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
	"golang.org/x/net/websocket"
)

// ErrNoLogs is returned when logs are missing from an app.
var ErrNoLogs = errors.New(
	`There are currently no log messages. Please check the following things:
1) Logger and fluentd pods are running: kubectl --namespace=drycc get pods.
2) The application is writing logs to the logger component by checking that an entry in the ring buffer was created: kubectl --namespace=drycc logs <logger pod>
3) Making sure that the container logs were mounted properly into the fluentd pod: kubectl --namespace=drycc exec <fluentd pod> ls /var/log/containers
3a) If the above command returns saying /var/log/containers cannot be found then please see the following github issue for a workaround: https://github.com/drycc/logger/issues/50`)

// List lists apps on a Drycc controller.
func List(c *drycc.Client, results int) (api.Apps, int, error) {
	body, count, reqErr := c.LimitedRequest("/v2/apps/", results)

	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.App{}, -1, reqErr
	}

	var apps []api.App
	if err := json.Unmarshal([]byte(body), &apps); err != nil {
		return []api.App{}, -1, err
	}

	return apps, count, reqErr
}

// New creates a new app with the given appID. Passing an empty string will result in
// a randomized app name.
//
// If the app name already exists, the error drycc.ErrDuplicateApp will be returned.
func New(c *drycc.Client, appID string) (api.App, error) {
	body := []byte{}

	if appID != "" {
		req := api.AppCreateRequest{ID: appID}
		b, err := json.Marshal(req)

		if err != nil {
			return api.App{}, err
		}
		body = b
	}

	res, reqErr := c.Request("POST", "/v2/apps/", body)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.App{}, reqErr
	}
	defer res.Body.Close()

	app := api.App{}
	if err := json.NewDecoder(res.Body).Decode(&app); err != nil {
		return api.App{}, err
	}

	return app, reqErr
}

// Get app details from a controller.
func Get(c *drycc.Client, appID string) (api.App, error) {
	u := fmt.Sprintf("/v2/apps/%s/", appID)

	res, reqErr := c.Request("GET", u, nil)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.App{}, reqErr
	}
	defer res.Body.Close()

	app := api.App{}

	if err := json.NewDecoder(res.Body).Decode(&app); err != nil {
		return api.App{}, err
	}

	return app, reqErr
}

// Logs retrieves logs from an app. The number of log lines fetched can be set by the lines
// argument. Setting lines = -1 will retrieve all app logs.
func Logs(c *drycc.Client, appID string, request api.AppLogsRequest) (*websocket.Conn, error) {
	scheme := "ws"
	if c.ControllerURL.Scheme == "https" {
		scheme = "wss"
	}
	path := fmt.Sprintf("v2/apps/%s/logs", appID)
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

// Run a one-time command in your app. This will start a kubernetes job with the
// same container image and environment as the rest of the app.
func Run(c *drycc.Client, appID string, command string, volumes map[string]interface{}) (api.AppRunResponse, error) {
	req := api.AppRunRequest{
		Command: command,
		Volumes: volumes,
	}
	body, err := json.Marshal(req)

	if err != nil {
		return api.AppRunResponse{}, err
	}

	u := fmt.Sprintf("/v2/apps/%s/run", appID)

	res, reqErr := c.Request("POST", u, body)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.AppRunResponse{}, reqErr
	}

	arr := api.AppRunResponse{}

	if err = json.NewDecoder(res.Body).Decode(&arr); err != nil {
		return api.AppRunResponse{}, err
	}

	return arr, reqErr
}

// Delete an app.
func Delete(c *drycc.Client, appID string) error {
	u := fmt.Sprintf("/v2/apps/%s/", appID)

	res, err := c.Request("DELETE", u, nil)
	if err == nil {
		res.Body.Close()
	}
	return err
}

// Transfer an app to another user.
func Transfer(c *drycc.Client, appID string, username string) error {
	u := fmt.Sprintf("/v2/apps/%s/", appID)

	req := api.AppUpdateRequest{Owner: username}
	body, err := json.Marshal(req)

	if err != nil {
		return err
	}

	res, err := c.Request("POST", u, body)
	if err == nil {
		res.Body.Close()
	}
	return err
}
