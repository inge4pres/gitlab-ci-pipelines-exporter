package cmd

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"net/http/httptest"
	"strings"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestRunWrongLogLevel(t *testing.T) {
	set := flag.NewFlagSet("", 0)
	set.String("log-level", "foo", "")
	set.String("log-format", "json", "")
	fmt.Println(Run(cli.NewContext(nil, set, nil)))
	err := Run(cli.NewContext(nil, set, nil))
	assert.Equal(t, true, strings.HasPrefix(err.Error(), "not a valid logrus Level"))
}

func TestRunWrongLogType(t *testing.T) {
	set := flag.NewFlagSet("", 0)
	set.String("log-level", "debug", "")
	set.String("log-format", "foo", "")
	err := Run(cli.NewContext(nil, set, nil))
	assert.Equal(t, true, strings.HasPrefix(err.Error(), "Invalid log format"))
}

func TestRunInvalidConfigFile(t *testing.T) {
	set := flag.NewFlagSet("", 0)
	set.String("log-level", "debug", "")
	set.String("log-format", "json", "")
	set.String("config", "path_does_not_exist", "")
	err := Run(cli.NewContext(nil, set, nil))
	assert.Equal(t, true, strings.HasPrefix(err.Error(), "couldn't open config file :"))
}

// introduce a test to check the /metrics endpoint body
func TestMetricsRegistryContainsMetricsWhenSet(t *testing.T) {
	// a custom additional metric added to the registry
	aMetric := "test_something"
	aCounter := prometheus.NewCounter(prometheus.CounterOpts{Name: aMetric})
	registry := newMetricsRegistry(aCounter)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	metricsHandlerFor(registry, false).ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), aMetric)
	//
	//// assert also the default metrics are in
	//assert.Contains(t, w.Body.String(), "gitlab_ci_pipeline_last_run_id")
}
