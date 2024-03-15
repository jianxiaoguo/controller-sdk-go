package api

// LimitSpec is the definition of GET /v2/limits/specs/
type LimitSpec struct {
	ID       string                 `json:"id"`
	CPU      map[string]interface{} `json:"cpu"`
	Memory   map[string]interface{} `json:"memory"`
	Features map[string]interface{} `json:"features"`
	Keywords []string               `json:"keywords"`
	Disabled bool                   `json:"disabled"`
}

// LimitPlan is the definition of GET /v2/limits/plans/
type LimitPlan struct {
	ID       string                 `json:"id"`
	Spec     LimitSpec              `json:"spec"`
	CPU      int                    `json:"cpu"`
	Memory   int                    `json:"memory"`
	Features map[string]interface{} `json:"features"`
	Disabled bool                   `json:"disabled"`
}
