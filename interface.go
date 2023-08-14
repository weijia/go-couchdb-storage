type CouchDb interface{
    SimpleFind()
}

type CouchConfig interface{
    GetConfig()
    GetGloabalConfig()
    GetDeviceUuid()
}

