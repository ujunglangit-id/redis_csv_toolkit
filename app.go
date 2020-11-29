package main

import (
	"flag"
	"fmt"
	"github.com/azwarnrst/redis_csv_toolkit/util"
	"github.com/garyburd/redigo/redis"
	"log"
	"os"
	"time"
)

var (
	cfg              *util.Config
	redisPool        *redis.Pool
	err              error
	redisImportLimit = 1000
	isTemporaryKey   bool
)

func main() {
	// initialize app configuration
	cfg = util.NewConfig()
	if err := cfg.ReadConfig(); err != nil {
		log.Fatalf("[Init][Configuration] %v", err)
	} else {
		log.Printf("[Init][Configuration] file loaded successfully")
	}

	redisHost := flag.String("host", "", "redis host")
	importLimit := flag.Int("limit", 30000, "Max concurrent import to redis")
	redisPort := flag.Int("port", 0, "redis port")
	csvInFile := flag.String("in", "", "input filename to import")
	ttl := flag.Int("ttl", 0, "redis key expiration time (seconds)")
	flag.Parse()

	if *redisHost != "" && *redisPort > 0 {
		cfg.RedisConfig.Host = fmt.Sprintf("%s:%d", *redisHost, *redisPort)
	}

	if *csvInFile != "" {
		cfg.AppConfig.FileName = *csvInFile
	} else {
		log.Printf("no file input specified ...")
		os.Exit(0)
	}

	if *importLimit > 0 {
		redisImportLimit = *importLimit
	}

	if *ttl > 0 {
		isTemporaryKey = true
	}

	start := time.Now()
	log.Printf("Target redis host : %s", cfg.RedisConfig.Host)
	redisConn, err := redis.Dial("tcp", cfg.RedisConfig.Host)
	if err != nil {
		log.Fatal(err)
	}
	csvUtil := util.NewCsv(cfg, redisConn, redisImportLimit)
	err = csvUtil.ParseCsv(isTemporaryKey, *ttl)
	if err != nil {
		log.Printf("parse error : %v\n", err)
	}
	t := time.Since(start)
	log.Printf("import completed in %f seconds", float64(t)/float64(time.Second))
	fmt.Println("")
}
