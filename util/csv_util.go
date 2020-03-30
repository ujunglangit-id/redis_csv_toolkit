package util

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/tokopedia/tdk/go/log"
	"io"
	"os"
	"strconv"
)

//Csv is struct
type Csv struct {
	cfg            *Config
	redisPool      *redis.Pool
	concurentLimit int
}

//NewCsv is method
func NewCsv(config *Config, pool *redis.Pool, limit int) *Csv {
	return &Csv{
		cfg:            config,
		redisPool:      pool,
		concurentLimit: limit,
	}
}

//OpenFile is method
func (c *Csv) OpenFile(file, fileType string) (*os.File, error) {
	if fileType == "input" {
		return os.Open(c.cfg.AppConfig.FileLocation + file)
	}
	if fileType == "output" {
		return os.OpenFile(c.cfg.AppConfig.FileLocation+file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	}
	return nil, errors.New("[Migration][OpenFile] wrong fileType")
}

func (c *Csv) ParseCsv() (err error) {
	var fileName = c.cfg.AppConfig.FileName
	shopList := []int{}
	isEOF := false
	csvFileInput, err := c.OpenFile(fileName, "input")
	if err != nil {
		return
	}

	reader := csv.NewReader(bufio.NewReader(csvFileInput))
	_, err = reader.Read()
	if err != nil {
		return
	}
	for !isEOF {
		if len(shopList) < c.concurentLimit {
			record, err := reader.Read()
			if err == io.EOF {
				isEOF = true
				break
			} else if err != nil {
				return err
			}

			shopID, err := strconv.Atoi(record[0])
			shopList = append(shopList, shopID)
		}

		if len(shopList) >= c.concurentLimit || isEOF {
			err = c.importRedis(shopList)
			if err != nil {
				log.Errorf("error import redis ")
			}
			shopList = []int{}
		}
	}

	return
}

func (c *Csv) importRedis(shopList []int) (err error) {
	for _, v := range shopList {
		if err := c.redisPool.Get().Send("SET", fmt.Sprintf(c.cfg.AppConfig.KeyFormat, v), 1); err != nil {
			log.Fatal(err)
		}
	}
	if err := c.redisPool.Get().Flush(); err != nil {
		log.Fatal(err)
	}

	go func() {
		for i := 0; i < len(shopList); i++ {
			_, err := c.redisPool.Get().Receive()
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
	return
}
