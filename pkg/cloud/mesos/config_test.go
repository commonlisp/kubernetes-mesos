package mesos

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	log "github.com/golang/glog"
)

// test mesos.createDefaultConfig
func Test_createDefaultConfig(t *testing.T) {
	defer log.Flush()

	config := createDefaultConfig()

	assert(t, "Default config has the expected MesosMaster value",
		config.MesosMaster == "localhost:5050")

	assert(t, "Default config has the expected MesosHttpClientTimeout value",
		config.MesosHttpClientTimeout.Duration == time.Duration(10)*time.Second)

	assert(t, "Default config has the expected StateCacheTTL value",
		config.StateCacheTTL.Duration == time.Duration(5)*time.Second)
}

// test mesos.readConfig
func Test_readConfig(t *testing.T) {
	defer log.Flush()

	configString := `
[mesos-cloud]
	mesos-master        = leader.mesos:5050
	http-client-timeout = 500ms
	state-cache-ttl     = 1h`

	reader := bytes.NewBufferString(configString)

	err := readConfig(reader)

	assert(t,
		fmt.Sprintf("Reading configuration does not yield an error: %#v", err),
		err == nil)

	assert(t,
		"Parsed config has the expected MesosMaster value",
		getConfig().MesosMaster == "leader.mesos:5050")

	assert(t,
		"Parsed config has the expected MesosHttpClientTimeout value",
		getConfig().MesosHttpClientTimeout.Duration == time.Duration(500)*time.Millisecond)

	assert(t,
		"Parsed config has the expected StateCacheTTL value",
		getConfig().StateCacheTTL.Duration == time.Duration(1)*time.Hour)
}
