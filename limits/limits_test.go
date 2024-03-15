package limits

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
	"github.com/stretchr/testify/require"
)

const specsFixture string = `
{
	"results": [
	  {
		"id": "std1",
		"cpu": {
		  "name": "Unknown CPU",
		  "cores": 32,
		  "clock": "3100MHZ",
		  "boost": "3700MHZ",
		  "threads": 64
		},
		"memory": {
		  "size": "64GB",
		  "type": "DDR4-ECC"
		},
		"features": {
		  "gpu": {
			"name": "Unknown Integrated GPU",
			"tmus": 1,
			"rops": 1,
			"cores": 128,
			"memory": {
			  "size": "shared",
			  "type": "shared"
			}
		  },
		  "network": "10G"
		},
		"keywords": [
		  "amd",
		  "intel",
		  "unknown"
		],
		"disabled": false
	  }
	],
	"count": 1
  }
`
const plansFixture string = `
{
	"results": [{
			"id": "std1.large.c1m1",
			"spec": {
				"id": "std1",
				"cpu": {
					"name": "Unknown CPU",
					"cores": 32,
					"clock": "3100MHZ",
					"boost": "3700MHZ",
					"threads": 64
				},
				"memory": {
					"size": "64GB",
					"type": "DDR4-ECC"
				},
				"features": {
					"gpu": {
						"name": "Unknown Integrated GPU",
						"tmus": 1,
						"rops": 1,
						"cores": 128,
						"memory": {
							"size": "shared",
							"type": "shared"
						}
					},
					"network": "10G"
				},
				"keywords": [
					"amd",
					"intel",
					"unknown"
				],
				"disabled": false
			},
			"cpu": 1,
			"memory": 1,
			"features": {
				"gpu": 1,
				"network": 1
			},
			"disabled": false
		},
		{
			"id": "std1.large.c1m2",
			"spec": {
				"id": "std1",
				"cpu": {
					"name": "Unknown CPU",
					"cores": 32,
					"clock": "3100MHZ",
					"boost": "3700MHZ",
					"threads": 64
				},
				"memory": {
					"size": "64GB",
					"type": "DDR4-ECC"
				},
				"features": {
					"gpu": {
						"name": "Unknown Integrated GPU",
						"tmus": 1,
						"rops": 1,
						"cores": 128,
						"memory": {
							"size": "shared",
							"type": "shared"
						}
					},
					"network": "10G"
				},
				"keywords": [
					"amd",
					"intel",
					"unknown"
				],
				"disabled": false
			},
			"cpu": 1,
			"memory": 2,
			"features": {
				"gpu": 1,
				"network": 1
			},
			"disabled": false
		}
	],
	"count": 2
}
`

const getPlanFixture string = `
{
	"id": "std1.large.c1m1",
	"spec": {
	  "id": "std1",
	  "cpu": {
		"name": "Unknown CPU",
		"cores": 32,
		"clock": "3100MHZ",
		"boost": "3700MHZ",
		"threads": 64
	  },
	  "memory": {
		"size": "64GB",
		"type": "DDR4-ECC"
	  },
	  "features": {
		"gpu": {
		  "name": "Unknown Integrated GPU",
		  "tmus": 1,
		  "rops": 1,
		  "cores": 128,
		  "memory": {
			"size": "shared",
			"type": "shared"
		  }
		},
		"network": "10G"
	  },
	  "keywords": [
		"amd",
		"intel",
		"unknown"
	  ],
	  "disabled": false
	},
	"cpu": 1,
	"memory": 1,
	"features": {
		"gpu": 1,
		"network": 1
	},
	"disabled": false
}
`

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)
	if req.URL.Path == "/v2/limits/specs/" && req.Method == "GET" {
		res.Write([]byte(specsFixture))
		return
	}

	if req.URL.Path == "/v2/limits/plans/" && req.Method == "GET" {
		res.Write([]byte(plansFixture))
		return
	}
	if req.URL.Path == "/v2/limits/plans/std1.large.c1m1/" && req.Method == "GET" {
		res.Write([]byte(getPlanFixture))
		return
	}

	fmt.Printf("Unrecongized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestSpecs(t *testing.T) {
	t.Parallel()
	expected := []api.LimitSpec{
		{
			ID: "std1",
			CPU: map[string]interface{}{
				"name":    "Unknown CPU",
				"cores":   32,
				"clock":   "3100MHZ",
				"boost":   "3700MHZ",
				"threads": 64,
			},
			Memory: map[string]interface{}{
				"size": "64GB",
				"type": "DDR4-ECC",
			},
			Features: map[string]interface{}{
				"gpu": map[string]interface{}{
					"name":  "Unknown Integrated GPU",
					"tmus":  1,
					"rops":  1,
					"cores": 128,
					"memory": map[string]interface{}{
						"size": "shared",
						"type": "shared",
					},
				},
				"network": "10G",
			},
			Keywords: []string{
				"amd",
				"intel",
				"unknown",
			},
			Disabled: false,
		},
	}
	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, _, err := Specs(drycc, "", 100)

	if err != nil {
		t.Fatal(err)
	}
	a, _ := json.Marshal(expected)
	b, _ := json.Marshal(actual)
	require.JSONEq(t, string(a), string(b), "Expected %v, Got %v", expected, actual)
}

func TestPlans(t *testing.T) {
	t.Parallel()
	spec := api.LimitSpec{
		ID: "std1",
		CPU: map[string]interface{}{
			"name":    "Unknown CPU",
			"cores":   32,
			"clock":   "3100MHZ",
			"boost":   "3700MHZ",
			"threads": 64,
		},
		Memory: map[string]interface{}{
			"size": "64GB",
			"type": "DDR4-ECC",
		},
		Features: map[string]interface{}{
			"gpu": map[string]interface{}{
				"name":  "Unknown Integrated GPU",
				"tmus":  1,
				"rops":  1,
				"cores": 128,
				"memory": map[string]interface{}{
					"size": "shared",
					"type": "shared",
				},
			},
			"network": "10G",
		},
		Keywords: []string{
			"amd",
			"intel",
			"unknown",
		},
		Disabled: false,
	}

	expected := []api.LimitPlan{
		{
			ID:     "std1.large.c1m1",
			Spec:   spec,
			CPU:    1,
			Memory: 1,
			Features: map[string]interface{}{
				"gpu":     1,
				"network": 1,
			},
			Disabled: false,
		},
		{
			ID:     "std1.large.c1m2",
			Spec:   spec,
			CPU:    1,
			Memory: 2,
			Features: map[string]interface{}{
				"gpu":     1,
				"network": 1,
			},
			Disabled: false,
		},
	}
	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, _, err := Plans(drycc, "", 0, 0, 100)

	if err != nil {
		t.Fatal(err)
	}

	a, _ := json.Marshal(expected)
	b, _ := json.Marshal(actual)
	require.JSONEq(t, string(a), string(b), "Expected %v, Got %v", expected, actual)
}

func TestGetPlan(t *testing.T) {
	spec := api.LimitSpec{
		ID: "std1",
		CPU: map[string]interface{}{
			"name":    "Unknown CPU",
			"cores":   32,
			"clock":   "3100MHZ",
			"boost":   "3700MHZ",
			"threads": 64,
		},
		Memory: map[string]interface{}{
			"size": "64GB",
			"type": "DDR4-ECC",
		},
		Features: map[string]interface{}{
			"gpu": map[string]interface{}{
				"name":  "Unknown Integrated GPU",
				"tmus":  1,
				"rops":  1,
				"cores": 128,
				"memory": map[string]interface{}{
					"size": "shared",
					"type": "shared",
				},
			},
			"network": "10G",
		},
		Keywords: []string{
			"amd",
			"intel",
			"unknown",
		},
		Disabled: false,
	}

	expected := api.LimitPlan{
		ID:     "std1.large.c1m1",
		Spec:   spec,
		CPU:    1,
		Memory: 1,
		Features: map[string]interface{}{
			"gpu":     1,
			"network": 1,
		},
		Disabled: false,
	}
	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()
	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := GetPlan(drycc, "std1.large.c1m1")

	if err != nil {
		t.Fatal(err)
	}

	a, _ := json.Marshal(expected)
	b, _ := json.Marshal(actual)
	require.JSONEq(t, string(a), string(b), "Expected %v, Got %v", expected, actual)
}
