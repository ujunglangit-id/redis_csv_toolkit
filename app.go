package main

import (
	"flag"
	"fmt"
	"github.com/azwarnrst/redis_csv_toolkit/util"
	"github.com/garyburd/redigo/redis"
	"github.com/tokopedia/tdk/go/log"
	"os"
	"time"
)

var (
	cfg              *util.Config
	redisPool        *redis.Pool
	err              error
	redisImportLimit = 1000
)

func main() {
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
	importLimit := flag.Int("limit", 10000, "Max concurrent import to redis")
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

	if *importLimit > 0 {
		redisImportLimit = *importLimit
	}

	start := time.Now()
	redisPool, err = util.InitRedisConnection(cfg)
	csvUtil := util.NewCsv(cfg, redisPool, redisImportLimit)
	err = csvUtil.ParseCsv()

	t := time.Since(start)
	fmt.Printf("import completed in %f seconds\n", float64(t)/float64(time.Second))
}
