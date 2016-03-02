package conf

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type DatabaseConfig struct {
	Host   string
	Port   int
	User   string
	DBName string
}

type RedisConfig struct {
	Host string
	Port int
}

const (
	DEFAULT_CACHE_DB_NAME = "default"
	DEFAULT_CACHE_EXPIRE  = 7 * 24 * 60 * 60
)

var (
	HttpPort  int
	DebugMode = false
	DBConfig  = make(map[string]*DatabaseConfig)
	RedConfig = make(map[string]*RedisConfig)

	RequestLogDir string
	ErrorLogDir   string

	debugMode         = flag.Bool("debug", false, "debug mode")
	serviceConfigFile = flag.String("service-config", "__unset__", "service config file")
)

func parseJsonInt(config map[string]interface{}, key string, logPrefix string) (res int) {
	if t, ok := config[key].(float64); !ok {
		log.Fatal(logPrefix + key + " not set correctly, should be int value")
	} else {
		res = int(t)
	}

	return
}

func parseJsonString(config map[string]interface{}, key string, logPrefix string) (res string) {
	if t, ok := config[key].(string); !ok {
		log.Fatal(logPrefix + key + " not set correctly, should be string value")
	} else {
		res = t
	}

	return
}

func setupDatabaseConfig(conf interface{}) {
	if conf == nil {
		log.Fatal("No config for database\n")
	}

	databaseConfData, ok := conf.(map[string]interface{})
	if !ok {
		log.Fatal("Database config not set correct\n")
	}

	databaseConfig := make(map[string]*DatabaseConfig)

	for dbName, content := range databaseConfData {
		databaseConfig[dbName] = new(DatabaseConfig)
		if dbConfig, ok := content.(map[string]interface{}); ok {
			databaseConfig[dbName].User = parseJsonString(dbConfig, "user", fmt.Sprintf("[database-conf][%s]", dbName))
			databaseConfig[dbName].Host = parseJsonString(dbConfig, "host", fmt.Sprintf("[database-conf][%s]", dbName))
			databaseConfig[dbName].Port = parseJsonInt(dbConfig, "port", fmt.Sprintf("[database-conf][%s]", dbName))
			databaseConfig[dbName].DBName = parseJsonString(dbConfig, "db", fmt.Sprintf("[database-conf][%s]", dbName))
		} else {
			log.Fatal("database %s config is not correct!!")
		}
	}

	s, _ := json.Marshal(databaseConfig)
	log.Printf("Database config: %v\n", string(s))

	DBConfig = databaseConfig
}

func setupRedisConfig(conf interface{}) {
	if conf == nil {
		log.Fatal("No config for redis\n")
	}

	redisConfData, ok := conf.(map[string]interface{})
	if !ok {
		log.Fatal("redis config not set correct\n")
	}

	redisConfig := make(map[string]*RedisConfig, 1)

	for dbName, content := range redisConfData {
		redisConfig[dbName] = new(RedisConfig)
		if dbConfig, ok := content.(map[string]interface{}); ok {
			redisConfig[dbName].Host = parseJsonString(dbConfig, "host", fmt.Sprintf("[redis-conf][%s]", dbName))
			redisConfig[dbName].Port = parseJsonInt(dbConfig, "port", fmt.Sprintf("[redis-conf][%s]", dbName))
		} else {
			log.Fatal("redis %s config is not correct!!")
		}
	}

	s, _ := json.Marshal(redisConfig)
	log.Printf("redis config: %v\n", string(s))

	RedConfig = redisConfig
}

func mkDir(config map[string]interface{}, key string, logPrefix string) (absPath string) {
	path := parseJsonString(config, key, logPrefix)
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(logPrefix + key + " not set correctly, not valid path")
	}

	err = os.MkdirAll(absPath, 0766)
	if err != nil {
		log.Fatal(logPrefix + key + " not set correctly, not valid path")
	}

	return
}

func init() {
	flag.Parse()

	DebugMode = *debugMode

	if *serviceConfigFile == "__unset__" {
		*serviceConfigFile = "./conf/service.json"
	}

	appConfig := make(map[string]interface{})
	appConfFileName, err := filepath.Abs(*serviceConfigFile)
	if err != nil {
		log.Fatal("Failed to format app_config_file path")
	}

	appConfFileContent, err := ioutil.ReadFile(appConfFileName)
	if err != nil {
		log.Fatal("Failed to read ", appConfFileName, ". error: %s", err)
	}

	err = json.Unmarshal(appConfFileContent, &appConfig)
	if err != nil {
		log.Println(err)
		return
	}

	setupDatabaseConfig(appConfig["database"])
	setupRedisConfig(appConfig["redis"])

	HttpPort = parseJsonInt(appConfig, "http_port", "")

	RequestLogDir = mkDir(appConfig, "request_log_dir", "")

	ErrorLogDir = mkDir(appConfig, "error_log_dir", "")

}
