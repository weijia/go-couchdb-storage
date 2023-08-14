package couchdb_storage

import (
	"testing"
)

func TestGetTimestamp(t *testing.T) {
	couchDbConfig := NewCouchDbConfig("test_server")
	couchDbConfig.GetTimestamp()
	couchDbConfig.GetConfig("test", "default_value")
}
