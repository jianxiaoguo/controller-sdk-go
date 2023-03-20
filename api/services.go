package api

// Service is the structure of the service object.
type Service struct {
	ProcfileType string `json:"procfile_type"`
	Port         int    `json:"port"`
	Protocol     string `json:"protocol"`
	TargetPort   int    `json:"target_port"`
}

// Services defines a collection of service objects.
type Services []Service

func (s Services) Len() int           { return len(s) }
func (s Services) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Services) Less(i, j int) bool { return s[i].ProcfileType < s[j].ProcfileType }

// ServiceCreateUpdateRequest is the structure of POST /v2/app/<app id>/services/.
type ServiceCreateUpdateRequest struct {
	ProcfileType string `json:"procfile_type"`
	Port         int    `json:"port"`
	Protocol     string `json:"protocol"`
	TargetPort   int    `json:"target_port"`
}

// ServiceDeleteRequest is the structure of DELETE /v2/app/<app id>/services/.
type ServiceDeleteRequest struct {
	ProcfileType string `json:"procfile_type"`
}
