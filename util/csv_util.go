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
	redisConn      redis.Conn
	concurentLimit int
}

//NewCsv is method
func NewCsv(config *Config, conn redis.Conn, limit int) *Csv {
	return &Csv{
		cfg:            config,
		redisConn:      conn,
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

func (c *Csv) ParseCsv(isTemporary bool, keyTtl int) (err error) {
	var lineCount int
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
				//break
			} else if err != nil {
				return err
			}
			if !isEOF {
				shopID, err := strconv.Atoi(record[0])
				if err != nil {
					log.Errorf("invalid shop id : %v\n", err)
				} else {
					lineCount++
					shopList = append(shopList, shopID)
				}
			}
		}

		if len(shopList) >= c.concurentLimit || isEOF {
			log.Infof("execute pipeline import")
			err = c.importRedis(shopList, isTemporary, keyTtl)
			if err != nil {
				log.Errorf("error import redis ")
			}
			shopList = []int{}
		}
	}

	log.Infof("total shop ID %d", lineCount)

	return
}

func (c *Csv) importRedis(shopList []int, isTemporary bool, keyTtl int) (err error) {
	log.Printf("length : %d", len(shopList))
	for _, v := range shopList {
		if isTemporary {
			if err := c.redisConn.Send("SETEX", fmt.Sprintf(c.cfg.AppConfig.KeyFormat, v), keyTtl, 1); err != nil {
				log.Errorf("pipeline error : %v\n", err)
			}
		} else {
			if err := c.redisConn.Send("SET", fmt.Sprintf(c.cfg.AppConfig.KeyFormat, v), 1); err != nil {
				log.Errorf("pipeline error : %v\n", err)
			}
		}
	}
	if err := c.redisConn.Flush(); err != nil {
		log.Errorf("flush error : %v\n", err)
	}

	go c.receive(c.redisConn, len(shopList))
	return
}

func (c *Csv) receive(conn redis.Conn, length int) {
	for i := 0; i < length; i++ {
		_, err := conn.Receive()
		if err != nil {
			log.Errorf("pipeline receive error : %v\n", err)
			break
		}
	}
}
