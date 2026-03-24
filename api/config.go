package api

import (
	"bytes"
	"fmt"
	"text/template"
)

// ConfigTags is the key, value for tag
type ConfigTags map[string]any

// ConfigVar represents a configuration variable for an app.
type ConfigVar struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

// ConfigValue represents a configuration value with its type and group.
type ConfigValue struct {
	Ptype string `json:"ptype,omitempty"`
	Group string `json:"group,omitempty"`
	ConfigVar
}

// PtypeValue represents values for a specific process type.
type PtypeValue struct {
	Env []ConfigVar `json:"env,omitempty"`
	Ref []string    `json:"ref,omitempty"`
}

// ConfigInfo represents the complete configuration information for an app.
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
	// App is the app name. It cannot be updated at all right now.
	App string `json:"app,omitempty"`
	// Values are exposed as environment variables to the app.
	Values []ConfigValue `json:"values,omitempty"`
	// Typed values are exposed as environment variables to the app.
	ValuesRefs ValuesRefs `json:"values_refs,omitempty"`
	// Limits is used to set process resources limits. The key is the process name
	// and the value is a limit plan. Ex: std1.xlarge.c1m1
	Limits map[string]any `json:"limits,omitempty"`
	// Timeout is used to set termination grace period. The key is the process name
	// and the value is a number in seconds, e.g. 30
	Timeout map[string]any `json:"termination_grace_period,omitempty"`
	// Lifecycle is a map of lifecycles for each process type.
	Lifecycle map[string]*Lifecycle `json:"lifecycle,omitempty"`
	// Healthcheck is map of healthchecks for each process that the application uses.
	Healthcheck map[string]*Healthcheck `json:"healthcheck,omitempty"`
	// Tags restrict applications to run on k8s nodes with that label.
	Tags map[string]ConfigTags `json:"tags,omitempty"`
	// Registry is a key-value pair to provide authentication for container registries.
	// The key is the username and the value is the password.
	Registry map[string]map[string]any `json:"registry,omitempty"`
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

// Lifecycle defines actions to take in the container lifecycle.
type Lifecycle struct {
	PostStart  **LifecycleHandler `json:"postStart,omitempty"`
	PreStop    **LifecycleHandler `json:"preStop,omitempty"`
	StopSignal string             `json:"stopSignal,omitempty"`
}

// LifecycleHandler defines actions to take in the container lifecycle.
type LifecycleHandler struct {
	Exec      *ExecAction      `json:"exec,omitempty"`
	HTTPGet   *HTTPGetAction   `json:"httpGet,omitempty"`
	Sleep     *SleepAction     `json:"sleep,omitempty"`
	TCPSocket *TCPSocketAction `json:"tcpSocket,omitempty"`
}

// String displays the LifecycleHandler in a readable format.
func (l LifecycleHandler) String() string {
	var doc bytes.Buffer
	tmpl, err := template.New("healthcheck").Parse(`Exec Probe: {{or .Exec "N/A"}}
GRPC Action: {{or .GRPC "N/A"}}
HTTP GET Action: {{or .HTTPGet "N/A"}}
Sleep Action: {{or .Sleep "N/A"}}
TCP Socket Action: {{or .TCPSocket "N/A"}}`)
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(&doc, l); err != nil {
		panic(err)
	}
	return doc.String()
}

// Healthcheck defines a container healthcheck.
type Healthcheck struct {
	StartupProbe   **ContainerProbe `json:"startupProbe,omitempty"`
	LivenessProbe  **ContainerProbe `json:"livenessProbe,omitempty"`
	ReadinessProbe **ContainerProbe `json:"readinessProbe,omitempty"`
}

// ContainerProbe defines a container healthcheck probe.
type ContainerProbe struct {
	InitialDelaySeconds int              `json:"initialDelaySeconds"`
	TimeoutSeconds      int              `json:"timeoutSeconds"`
	PeriodSeconds       int              `json:"periodSeconds"`
	SuccessThreshold    int              `json:"successThreshold"`
	FailureThreshold    int              `json:"failureThreshold"`
	Exec                *ExecAction      `json:"exec,omitempty"`
	GRPC                *GRPCAction      `json:"grpc,omitempty"`
	HTTPGet             *HTTPGetAction   `json:"httpGet,omitempty"`
	TCPSocket           *TCPSocketAction `json:"tcpSocket,omitempty"`
}

// String displays the ContainerProbe in a readable format.
func (c ContainerProbe) String() string {
	var doc bytes.Buffer
	tmpl, err := template.New("healthcheck").Parse(`Initial Delay (seconds): {{.InitialDelaySeconds}}
Timeout (seconds): {{.TimeoutSeconds}}
Period (seconds): {{.PeriodSeconds}}
Success Threshold: {{.SuccessThreshold}}
Failure Threshold: {{.FailureThreshold}}
Exec Probe: {{or .Exec "N/A"}}
GRPC Probe: {{or .GRPC "N/A"}}
HTTP GET Probe: {{or .HTTPGet "N/A"}}
TCP Socket Probe: {{or .TCPSocket "N/A"}}`)
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(&doc, c); err != nil {
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
type ExecAction struct {
	Command []string `json:"command"`
}

// String displays the ExecAction in a readable format.
func (e ExecAction) String() string {
	return fmt.Sprintf(`Command=%s`, e.Command)
}

// HTTPGetAction performs an HTTP GET request to the Pod
// with the given path, port and headers.
type HTTPGetAction struct {
	Path        string    `json:"path,omitempty"`
	Port        int       `json:"port"`
	HTTPHeaders []*KVPair `json:"httpHeaders,omitempty"`
}

// String displays the HTTPGetAction in a readable format.
func (h HTTPGetAction) String() string {
	return fmt.Sprintf(`Path="%s" Port=%d HTTPHeaders=%s`,
		h.Path,
		h.Port,
		h.HTTPHeaders)
}

// TCPSocketAction attempts to open a socket connection to the
// Pod on the given port.
type TCPSocketAction struct {
	Port int `json:"port"`
}

// String displays the TCPSocketAction in a readable format.
func (t TCPSocketAction) String() string {
	return fmt.Sprintf("Port=%d", t.Port)
}

// GRPCAction performs an GRPC request to the Pod
// with the given path, port and headers.
type GRPCAction struct {
	Port    int    `json:"port"`
	Service string `json:"service,omitempty"`
}

// String displays the GRPCAction in a readable format.
func (g GRPCAction) String() string {
	return fmt.Sprintf(`Port=%d Service="%s"`, g.Port, g.Service)
}

// SleepAction pauses for a specified number of seconds.
type SleepAction struct {
	Seconds int `json:"seconds"`
}

// String displays the SleepAction in a readable format.
func (s SleepAction) String() string {
	return fmt.Sprintf(`Seconds=%d`, s.Seconds)
}
