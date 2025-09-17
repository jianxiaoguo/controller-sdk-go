package drycc

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type versionComparison struct {
	Client string
	Server string
	Error  error
}

func TestCheckAPIVersions(t *testing.T) {
	comparisons := []versionComparison{
		{"1.2", "2.1", ErrAPIMismatch},
		{"2.1", "1.2", ErrAPIMismatch},
		{"2.1", "2.2", ErrAPIMismatch},
		{"2.3", "2.0", nil},
	}

	for _, check := range comparisons {
		err := CheckAPICompatibility(check.Client, check.Server)

		if err != check.Error {
			t.Errorf("%v: Expected %v, Got %v", check, check.Error, err)
		}
	}
}

func TestParseEnv(t *testing.T) {
	expects := map[string]string{
		"F1": "111",
		"F2": "222",
		"F3": "333",
	}

	f, err := os.CreateTemp("", "env")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(f.Name())
	for key, value := range expects {
		fmt.Fprintf(f, "%s=%s\n", key, value)
	}
	f.Seek(0, 0)

	config, err := ParseEnv(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	if len(config) != 3 {
		t.Errorf("Expected %d, Got %d", 3, len(config))
	}
	for key, value := range expects {
		if v, ok := config[key]; !ok || config[key] != value {
			t.Errorf("Expected %s, Got %s", v, value)
		}
	}
}

func TestParseDryccfile(t *testing.T) {
	webPipeline := `
kind: pipeline
ptype: web
build:
  docker: Dockerfile
  arg:
    CODENAME: bookworm
env:
  VERSION: 1.2.1
run:
  command:
  - ./deployment-tasks.sh
  image: worker
  timeout: 100
config:
- jvmconfig
deploy:
  command:
  - bash
  - -ec
  args:
  - bundle exec puma -C config/puma.rb
`
	taskPipeline := `
kind: pipeline
ptype: task
build:
  docker: Dockerfile
  arg:
    CODENAME: bookworm
env:
  VERSION: 1.2.1
run:
  command:
  - ./deployment-tasks.sh
  image: worker
  timeout: 100
config:
- jvmconfig
deploy:
  command:
  - bash
  - -ec
  args:
  - bundle exec puma -C config/puma.rb
`

	tmp, err := os.MkdirTemp("", "drycc")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	configDir := filepath.Join(tmp, "config")
	os.MkdirAll(configDir, 0o777)
	os.WriteFile(filepath.Join(tmp, "web.yml"), []byte(webPipeline), 0o777)
	os.WriteFile(filepath.Join(tmp, "task.yaml"), []byte(taskPipeline), 0o777)
	os.WriteFile(filepath.Join(configDir, "global"), []byte("DEBUG=true\n"), 0o777)
	os.WriteFile(filepath.Join(configDir, "web"), []byte("PORT=8000\nJVM_OPTIONS=-Xms16G\n"), 0o777)

	dryccfile, err := ParseDryccfile(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if config, ok := dryccfile["config"]; ok {
		webConfig := config.(map[string]interface{})["web"].(map[string]interface{})
		globalConfig := config.(map[string]interface{})["global"].(map[string]interface{})

		if globalConfig["DEBUG"] != "true" {
			t.Errorf("Expected %s, Got %s", "true", globalConfig["DEBUG"])
		}
		if webConfig["PORT"].(string) != "8000" {
			t.Errorf("Expected %s, Got %s", "8000", webConfig["PORT"])
		}
		if webConfig["JVM_OPTIONS"].(string) != "-Xms16G" {
			t.Errorf("Expected %s, Got %s", "-Xms16G", webConfig["JVM_OPTIONS"])
		}
	} else {
		t.Error("Config not found")
	}

	if pipeline, ok := dryccfile["pipeline"]; ok {
		if _, ok := pipeline.(map[string]interface{})["web.yml"]; !ok {
			t.Error("web.yml not found")
		}
		if _, ok := pipeline.(map[string]interface{})["task.yaml"]; !ok {
			t.Error("task.yaml not found")
		}
	} else {
		t.Error("Pipeline not found")
	}
}
