package mesos

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	log "github.com/golang/glog"
)

func assert(t *testing.T, fact string, condition bool) {
	if !condition {
		t.Fatalf("assertion failure: %#v", fact)
	}
}

// test mesos.newMesosCloud with no config
func Test_newMesosCloud_NoConfig(t *testing.T) {
	defer log.Flush()

	resetConfigToDefault()

	mesosCloud, err := newMesosCloud(nil)

	assert(t,
		fmt.Sprintf("Creating a new Mesos cloud provider without config does not yield an error: %#v", err),
		err == nil)

	assert(t, "Mesos client with default config has the expected HTTP client timeout value",
		mesosCloud.client.httpClient.Timeout == time.Duration(10)*time.Second)

	assert(t, "Mesos client with default config has the expected state cache TTL value",
		mesosCloud.client.state.ttl == time.Duration(5)*time.Second)
}

// test mesos.newMesosCloud with custom config
func Test_newMesosCloud_WithConfig(t *testing.T) {
	defer log.Flush()

	resetConfigToDefault()

	configString := `
[mesos-cloud]
	http-client-timeout = 500ms
	state-cache-ttl = 1h`

	reader := bytes.NewBufferString(configString)

	mesosCloud, err := newMesosCloud(reader)

	assert(t,
		fmt.Sprintf("Creating a new Mesos cloud provider with a custom config does not yield an error: %#v", err),
		err == nil)

	assert(t, "Mesos client with a custom config has the expected HTTP client timeout value",
		mesosCloud.client.httpClient.Timeout == time.Duration(500)*time.Millisecond)

	assert(t, "Mesos client with a custom config has the expected state cache TTL value",
		mesosCloud.client.state.ttl == time.Duration(1)*time.Hour)
}
