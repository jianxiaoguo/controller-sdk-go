package drycc

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

// ParseEnv parses environment variables from a file.
func ParseEnv(fileame string) (map[string]interface{}, error) {
	contents, err := os.ReadFile(fileame)
	if err != nil {
		return nil, err
	}
	configMap := make(map[string]interface{})
	regex := regexp.MustCompile(`^([A-z0-9_\-\.]+)=([\s\S]*)$`)
	for _, config := range strings.Split(string(contents), "\n") {
		// Skip config that starts with an comment
		if len(config) == 0 || config[0] == '#' {
			continue
		}

		if regex.MatchString(config) {
			captures := regex.FindStringSubmatch(config)
			configMap[captures[1]] = captures[2]
		} else {
			return nil, fmt.Errorf("'%s' does not match the pattern 'key=var', ex: MODE=test", config)
		}
	}

	return configMap, nil
}

// ParseDryccfile parses a Drycc configuration file.
func ParseDryccfile(dryccpath string) (map[string]interface{}, error) {
	config := make(map[string]interface{})
	if entries, err := os.ReadDir(path.Join(dryccpath, "config")); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				if env, err := ParseEnv(path.Join(dryccpath, "config", entry.Name())); err == nil {
					config[entry.Name()] = env
				} else {
					return nil, err
				}
			}
		}
	}
	pipeline := make(map[string]interface{})
	if entries, err := os.ReadDir(dryccpath); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && (strings.HasSuffix(entry.Name(), ".yaml") || strings.HasSuffix(entry.Name(), ".yml")) {
				data := make(map[string]interface{})
				if bytes, err := os.ReadFile(path.Join(dryccpath, entry.Name())); err == nil {
					if err = yaml.Unmarshal([]byte(bytes), data); err == nil {
						pipeline[entry.Name()] = data
					} else {
						return nil, err
					}
				} else {
					return nil, err
				}
			}
		}
	}

	dryccfile := make(map[string]interface{})
	if len(config) > 0 {
		dryccfile["config"] = config
	}
	if len(pipeline) > 0 {
		dryccfile["pipeline"] = pipeline
	}
	return dryccfile, nil
}

// CheckAPICompatibility checks if the server and client API versions are compatible.
func CheckAPICompatibility(serverAPIVersion, clientAPIVersion string) error {
	sVersion := strings.Split(serverAPIVersion, ".")
	aVersion := strings.Split(clientAPIVersion, ".")

	// If API Versions are invalid, return a mismatch error.
	if len(sVersion) < 2 || len(aVersion) < 2 {
		return ErrAPIMismatch
	}

	// If major versions are different, return a mismatch error.
	if sVersion[0] != aVersion[0] {
		return ErrAPIMismatch
	}

	// If server is older than client, return mismatch error.
	if sVersion[1] < aVersion[1] {
		return ErrAPIMismatch
	}

	return nil
}
