package api

import (
	"bytes"
	"fmt"
	"text/template"
)

// ConfigTags is the key, value for tag
type ConfigTags map[string]interface{}

// ConfigValues value for env
type ConfigVar struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type ConfigValue struct {
	Ptype string `json:"ptype,omitempty"`
	Group string `json:"group,omitempty"`
	ConfigVar
}

type PtypeValue struct {
	Env []ConfigVar `json:"env,omitempty"`
	Ref []string    `json:"ref,omitempty"`
}

type ConfigInfo struct {
	Ptype map[string]PtypeValue  `json:"ptype,omitempty"`
	Group map[string][]ConfigVar `json:"group,omitempty"`
}

// ValuesRefs is the key, value for refs
type ValuesRefs map[string][]string

// ConfigSet is the definition of POST /v2/apps/<app id>/config/.
type ConfigSet struct {
	Values []ConfigValue `json:"values"`
}

// ConfigUnset is the definition of POST /v2/apps/<app id>/config/.
type ConfigUnset struct {
	Values []ConfigValue `json:"values"`
}

// Config is the structure of an app's config.
type Config struct {
	// Owner is the app owner. It cannot be updated with config.Set(). See app.Transfer().
	Owner string `json:"owner,omitempty"`
	// App is the app name. It cannot be updated at all right now.
	App string `json:"app,omitempty"`
	// Values are exposed as environment variables to the app.
	Values []ConfigValue `json:"values,omitempty"`
	// Typed values are exposed as environment variables to the app.
	ValuesRefs ValuesRefs `json:"values_refs,omitempty"`
	// Limits is used to set process resources limits. The key is the process name
	// and the value is a limit plan. Ex: std1.xlarge.c1m1
	Limits map[string]interface{} `json:"limits,omitempty"`
	// Timeout is used to set termination grace period. The key is the process name
	// and the value is a number in seconds, e.g. 30
	Timeout map[string]interface{} `json:"termination_grace_period,omitempty"`
	// Healthcheck is map of healthchecks for each process that the application uses.
	Healthcheck map[string]*Healthchecks `json:"healthcheck,omitempty"`
	// Tags restrict applications to run on k8s nodes with that label.
	Tags map[string]ConfigTags `json:"tags,omitempty"`
	// Registry is a key-value pair to provide authentication for container registries.
	// The key is the username and the value is the password.
	Registry map[string]map[string]interface{} `json:"registry,omitempty"`
	// Created is the time that the application was created and cannot be updated.
	Created string `json:"created,omitempty"`
	// Updated is the last time the configuration was changed and cannot be updated.
	Updated string `json:"updated,omitempty"`
	// UUID is a unique string reflecting the configuration in its current state.
	// It changes every time the configuration is changed and cannot be updated.
	UUID string `json:"uuid,omitempty"`
}

// ConfigHookRequest defines the request for configuration from the config hook.
type ConfigHookRequest struct {
	User string `json:"receive_user"`
	App  string `json:"receive_repo"`
}

// Healthchecks is a map of healthcheck probes.
// The valid values are "startupProbe" "livenessProbe" and "readinessProbe".
type Healthchecks map[string]*Healthcheck

// Healthcheck is the structure for an application healthcheck.
// Healthchecks only need to provide information about themselves.
// All the information is pushed to the server and handled by kubernetes.
type Healthcheck struct {
	InitialDelaySeconds int             `json:"initialDelaySeconds"`
	TimeoutSeconds      int             `json:"timeoutSeconds"`
	PeriodSeconds       int             `json:"periodSeconds"`
	SuccessThreshold    int             `json:"successThreshold"`
	FailureThreshold    int             `json:"failureThreshold"`
	Exec                *ExecProbe      `json:"exec,omitempty"`
	HTTPGet             *HTTPGetProbe   `json:"httpGet,omitempty"`
	TCPSocket           *TCPSocketProbe `json:"tcpSocket,omitempty"`
}

// String displays the HealthcheckHTTPGetProbe in a readable format.
func (h Healthcheck) String() string {
	var doc bytes.Buffer
	tmpl, err := template.New("healthcheck").Parse(`Initial Delay (seconds): {{.InitialDelaySeconds}}
Timeout (seconds): {{.TimeoutSeconds}}
Period (seconds): {{.PeriodSeconds}}
Success Threshold: {{.SuccessThreshold}}
Failure Threshold: {{.FailureThreshold}}
Exec Probe: {{or .Exec "N/A"}}
HTTP GET Probe: {{or .HTTPGet "N/A"}}
TCP Socket Probe: {{or .TCPSocket "N/A"}}`)
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(&doc, h); err != nil {
		panic(err)
	}
	return doc.String()
}

// KVPair is a key/value pair used to parse values from
// strings into a formal structure.
type KVPair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (k KVPair) String() string {
	return k.Name + "=" + k.Value
}

// ExecProbe executes a command within a Pod.
type ExecProbe struct {
	Command []string `json:"command"`
}

// String displays the ExecProbe in a readable format.
func (e ExecProbe) String() string {
	return fmt.Sprintf(`Command=%s`, e.Command)
}

// HTTPGetProbe performs an HTTP GET request to the Pod
// with the given path, port and headers.
type HTTPGetProbe struct {
	Path        string    `json:"path,omitempty"`
	Port        int       `json:"port"`
	HTTPHeaders []*KVPair `json:"httpHeaders,omitempty"`
}

// String displays the HTTPGetProbe in a readable format.
func (h HTTPGetProbe) String() string {
	return fmt.Sprintf(`Path="%s" Port=%d HTTPHeaders=%s`,
		h.Path,
		h.Port,
		h.HTTPHeaders)
}

// TCPSocketProbe attempts to open a socket connection to the
// Pod on the given port.
type TCPSocketProbe struct {
	Port int `json:"port"`
}

// String displays the TCPSocketProbe in a readable format.
func (t TCPSocketProbe) String() string {
	return fmt.Sprintf("Port=%d", t.Port)
}
