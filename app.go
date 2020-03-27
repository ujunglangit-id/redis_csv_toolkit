package main

import (
	"flag"
	"fmt"
	"github.com/azwarnrst/redis_csv_toolkit/util"
	"github.com/garyburd/redigo/redis"
	"github.com/tokopedia/tdk/go/log"
	"os"
)

var (
	cfg       *util.Config
	redisPool *redis.Pool
	err       error
)

func init() {
	err = log.SetConfig(&log.Config{
		Level:   "info",
		AppName: "Redis Toolkit",
	})
	if err != nil {
		log.Fatal(err)
	}

	// initialize app configuration
	cfg = util.NewConfig()
	if err := cfg.ReadConfig(); err != nil {
		log.Fatalf("[Init][Configuration] %v", err)
	} else {
		log.Infof("[Init][Configuration] file loaded successfully")
	}

	redisHost := flag.String("host", "127.0.0.1", "redis host")
	redisPort := flag.Int("port", 6379, "redis port")
	csvInFile := flag.String("in", "", "input filename to import")
	flag.Parse()

	if *redisHost != "" && *redisPort > 0 {
		cfg.RedisConfig.Host = fmt.Sprintf("%s:%d", *redisHost, *redisPort)
	}

	if *csvInFile != "" {
		cfg.AppConfig.FileName = *csvInFile
	} else {
		log.Errorf("no file input specified ...")
		os.Exit(0)
	}

	redisPool, err = util.InitRedisConnection(cfg)
}

func main() {
	log.Infof("ihik ihik ...")
}
