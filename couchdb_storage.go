package couchdb_storage

import (
    "context"
    "fmt"
    "log"

    kivik "github.com/go-kivik/kivik/v4"
    _ "github.com/go-kivik/couchdb/v4" // The CouchDB driver
  "github.com/spf13/viper"
  "github.com/google/uuid"
)

type CouchDbConfig struct {
  KivikClient kivik.Client
  MainDbName string
  MainDb kivik.DB
  DeviceUuid string
}

type ServerConfig struct {
	UUID string `json:"uuid"`
	CouchDbServer string `json:"couch_server"`
	CouchDbUser string `json:"couch_user"`
	MainDb string `json:"main_db"`
	CouchPassword string `json:"couch_password"`
}

func NewCouchDbConfig(CouchDbServer string) *CouchDbConfig {
  const CONFIG_FILE_NAME := "couch_config.json"
  viper.SetConfigName("iot_go.json")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
    uuid := uuid.New()
    key := uuid.String()
		defaultServerConfig := ServerConfig {
      UUID: key,
      CouchDbServer: CouchDbServer,
      CouchDbUser: "test",
      MainDb: "production",
      CouchDbPassword："test"
    }
		defaultLocalConfig, _ := json.Marshal(defaultServerConfig)
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
  server := vip.GetString("couch_server")
  mainDbName := vip.GetString("main_db")
  couchDbUsername := vip.GetString("couch_user")
  couchDbPassword := vip.GetString("couch_password")
  serverUrl := "http://"+couchDbUsername+":"+couchDbPassword+"@"+server+":5984"
  client, err := kivik.New("couch", serverUrl)
  mainDb := client.DB(context.Background(), mainDbName)
  if err != nil {
    log.Fatal(err)
  }

  couchDbConfig := CouchDbConfig{
    KivikClient: client,
    MainDbName: mainDbName,
    MainDb: mainDb
    DeviceUuid: vip.GetString("uuid")
  }
  return &couchDbConfig
}

func (couchDbStorage CouchDbStorage) SimpleFind(query interface{}) []map[string]*json.RawMessage {
  res = couchDbStorage.MainDb.Find(context.TODO(), query)
  for res.Next() {
    log.Printf("Doc ID: %s\n", changes.ID())
    doc := &MyDoc{}
    if err := res.ScanDoc(doc); err != nil {
            log.Fatal(err)
    }
    // 在这里处理文档
    row := db.Get(context.TODO(), changes.ID())
    if err != nil { panic(err) }

    b, err := io.ReadAll(row.Body)
    if err != nil {
            log.Fatal(err)
    }

    fmt.Printf("%s", b)
    fmt.Printf("Got doc: %+v\n", doc)
	  
    var objs []map[string]*json.RawMessage
    if err := json.Unmarshal([]byte(b), &objs); err != nil {
        log.Fatal(err)
    }
    return objs
  }
}

func (couchDbStorage CouchDbStorage) GetConfig(key string) []map[string]*json.RawMessage {
  query:= map["string"] interface {} {
    "selecor": map["string"] interface {
      "type": "config",
      "device_uuid": couchDbStorage.Uuid,
      "name": key,
    },
    "skip": 0,
    "limit": 5,
  }
  couchDbStorage.SimpleFind(query)
}

func (couchDbStorage CouchDbStorage) GetGlobalConfig(key string) []map[string]*json.RawMessage {
  query:= map["string"] interface {} {
    "selecor": map["string"] interface {
      "type": "config",
      "name": key,
    },
    "skip": 0,
    "limit": 5,
  }
  couchDbStorage.SimpleFind(query)
}
