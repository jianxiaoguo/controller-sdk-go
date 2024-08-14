package api

import (
	"bytes"
	"text/template"
)

// AppSettings is the structure of an app's settings.
type AppSettings struct {
	// Owner is the app owner. It cannot be updated with AppSettings.Set(). See app.Transfer().
	Owner string `json:"owner,omitempty"`
	// App is the app name. It cannot be updated at all right now.
	App string `json:"app,omitempty"`
	// Created is the time that the application settings was created and cannot be updated.
	Created string `json:"created,omitempty"`
	// Updated is the last time the application settings was changed and cannot be updated.
	Updated string `json:"updated,omitempty"`
	// UUID is a unique string reflecting the application settings in its current state.
	// It changes every time the application settings is changed and cannot be updated.
	UUID string `json:"uuid,omitempty"`
	// Routable determines if the application should be exposed by the router.
	Routable     *bool                 `json:"routable,omitempty"`
	Allowlist    []string              `json:"allowlist,omitempty"`
	Autorollback *bool                 `json:"autorollback,omitempty"`
	Autoscale    map[string]*Autoscale `json:"autoscale,omitempty"`
	Label        Labels                `json:"label,omitempty"`
}

// NewRoutable returns a default value for the AppSettings.Routable field.
func NewRoutable() *bool {
	b := true
	return &b
}

// NewAutorollback returns a default value for the AppSettings.Autorollback field.
func NewAutorollback() *bool {
	b := true
	return &b
}

// Allowlist is the structure of POST /v2/app/<app id>/allowlist/.
type Allowlist struct {
	Addresses []string `json:"addresses"`
}

// String displays the Autoscale rule in a readable format.
func (a Autoscale) String() string {
	var doc bytes.Buffer
	tmpl, err := template.New("autoscale").Parse(`Min Replicas: {{.Min}}
Max Replicas: {{.Max}}
CPU: {{.CPUPercent}}%`)
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(&doc, a); err != nil {
		panic(err)
	}
	return doc.String()
}

// Autoscales contains a hash of process types and the autoscale rules
type Autoscales map[string]*Autoscale

// Autoscale is a per proc type scaling information
type Autoscale struct {
	Min        int `json:"min"`
	Max        int `json:"max"`
	CPUPercent int `json:"cpu_percent"`
}

// Labels can contain any user-defined key value
type Labels map[string]interface{}
