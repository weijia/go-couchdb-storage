package couchdb_storage

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	_ "github.com/go-kivik/couchdb/v4" // The CouchDB driver
	kivik "github.com/go-kivik/kivik/v4"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type CouchDbStorage struct {
	KivikClient *kivik.Client
	MainDbName  string
	MainDb      *kivik.DB
	DeviceUuid  string
}

type ServerConfig struct {
	UUID            string `json:"uuid"`
	CouchDbServer   string `json:"couch_server"`
	CouchDbUser     string `json:"couch_user"`
	MainDb          string `json:"main_db"`
	CouchDbPassword string `json:"couch_password"`
}

func NewCouchDbConfig(CouchDbServer string) *CouchDbStorage {
	const CONFIG_FILE_NAME string = "couch_config.json"
	viper.SetConfigName(CONFIG_FILE_NAME)
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		uuid := uuid.New()
		key := uuid.String()
		defaultServerConfig := ServerConfig{
			UUID:            key,
			CouchDbServer:   CouchDbServer,
			CouchDbUser:     "test",
			MainDb:          "production",
			CouchDbPassword: "test",
		}
		defaultLocalConfig, _ := json.MarshalIndent(defaultServerConfig, "", "  ")
		/*******************  使用 ioutil.WriteFile 写入文件 *****************/
		err2 := os.WriteFile("./"+CONFIG_FILE_NAME, defaultLocalConfig, 0666) //写入文件(字节数组)
		if err2 != nil {
			log.Fatal(err2)
		}

		secondErr := viper.ReadInConfig()
		if secondErr != nil {
			log.Fatal(secondErr)
		}
		// viper.WriteConfig()
	}
	server := viper.GetString("couch_server")
	mainDbName := viper.GetString("main_db")
	couchDbUsername := viper.GetString("couch_user")
	couchDbPassword := viper.GetString("couch_password")
	serverUrl := "http://" + couchDbUsername + ":" + couchDbPassword + "@" + server + ":5984"
	client, err := kivik.New("couch", serverUrl)
	mainDb := client.DB(mainDbName)
	if err != nil {
		log.Fatal(err)
	}

	couchDbConfig := CouchDbStorage{
		KivikClient: client,
		MainDbName:  mainDbName,
		MainDb:      mainDb,
		DeviceUuid:  viper.GetString("uuid"),
	}
	return &couchDbConfig
}

func (couchDbStorage CouchDbStorage) SimpleFind(query interface{}) map[string]interface{} {
	res := couchDbStorage.MainDb.Find(context.TODO(), query)
	for res.Next() {
		// log.Printf("Doc ID: %s\n", res.ID())
		var doc interface{}
		if err := res.ScanDoc(&doc); err != nil {
			log.Fatal(err)
		}
		return doc.(map[string]interface{})
	}
	return nil
}

func (couchDbStorage CouchDbStorage) GetConfig(key string, defaultValue string) string {
	base := map[string]interface{}{
		"type":        "config",
		"device_uuid": couchDbStorage.DeviceUuid,
		"name":        key,
	}
	return couchDbStorage.GetOrCreateConfig(base, defaultValue)
}

func (couchDbStorage CouchDbStorage) GetGlobalConfig(key string, defaultValue string) string {
	base := map[string]interface{}{
		"type": "config",
		"name": key,
	}
	return couchDbStorage.GetOrCreateConfig(base, defaultValue)
}

func (couchDbStorage CouchDbStorage) GetOrCreateConfig(base map[string]interface{}, defaultValue string) string {
	query := map[string]interface{}{
		"selector": base,
		"skip":     0,
		"limit":    5,
	}
	res := couchDbStorage.SimpleFind(query)
	if res == nil {
		base["value"] = defaultValue
		base["created"] = couchDbStorage.GetTimestamp()
		couchDbStorage.MainDb.CreateDoc(context.TODO(), base)
		return defaultValue
	}
	return res["value"].(string)
}

func (couchDbStorage CouchDbStorage) GetTimestamp() int64 {
	return time.Now().Unix()
}

func (couchDbStorage CouchDbStorage) GetDeviceUuid() string {
	return couchDbStorage.DeviceUuid
}
