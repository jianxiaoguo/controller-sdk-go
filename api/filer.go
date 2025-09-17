package api

// FilerDirEntry represents a directory entry in the file system.
type FilerDirEntry struct {
	Name      string `json:"name,omitempty"`
	Path      string `json:"path,omitempty"`
	Size      string `json:"size,omitempty"`
	Type      string `json:"type,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

// FilerDirEntries is a collection of FilerDirEntry.
type FilerDirEntries []FilerDirEntry
