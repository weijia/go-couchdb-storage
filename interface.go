package couchdb_storage

type CouchDbInterface interface {
	SimpleFind()
}

type CouchConfigInterface interface {
	GetConfig()
	GetGlobalConfig()
	GetDeviceUuid()
}
