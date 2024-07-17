package api

type FilerDirEntry struct {
	Name      string `json:"name,omitempty"`
	Path      string `json:"path,omitempty"`
	Size      string `json:"size,omitempty"`
	Type      string `json:"type,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

type FilerDirEntries []FilerDirEntry
