package api

// Ptype defines the structure of ptype deployment.
type Ptype struct {
	Name              string `json:"name"`
	Release           string `json:"release"`
	Ready             string `json:"ready"`
	UpToDate          int    `json:"up_to_date"`
	AvailableReplicas int    `json:"available_replicas"`
	Started           string `json:"started"`
	Garbage           bool   `json:"garbage"`
}

// Ptypes defines a collection of app Ptypes.
type Ptypes []Ptype

func (d Ptypes) Len() int           { return len(d) }
func (d Ptypes) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d Ptypes) Less(i, j int) bool { return d[i].Name < d[j].Name }

// PtypeState defines a ptype deployment state.
type PtypeState struct {
	Container      string            `json:"container"`
	Image          string            `json:"image"`
	Command        []string          `json:"command,omitempty"`
	Args           []string          `json:"args,omitempty"`
	StartupProbe   Healthcheck       `json:"startup_probe,omitempty"`
	LivenessProbe  Healthcheck       `json:"liveness_probe,omitempty"`
	ReadinessProbe Healthcheck       `json:"readiness_probe,omitempty"`
	Limits         map[string]string `json:"limits,omitempty"`
	VolumeMounts   []VolumeMount     `json:"volume_mounts,omitempty"`
	NodeSelector   map[string]string `json:"node_selector,omitempty"`
}

type VolumeMount struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
}

// PtypesState defines a collection of container state.
type PtypeStates []PtypeState
