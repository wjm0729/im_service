/**
 * Copyright (c) 2014-2015, GoBelieve
 * All rights reserved.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 59 Temple Place, Suite 330, Boston, MA  02111-1307  USA
 */
package main

import (
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/golang/glog"
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"time"
)

var redis_pool *redis.Pool
var config *Config

var (
	VERSION       string
	BUILD_TIME    string
	GO_VERSION    string
	GIT_COMMIT_ID string
	GIT_BRANCH    string
)

func NewRedisPool(server, password string, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     100,
		MaxActive:   500,
		IdleTimeout: 480 * time.Second,
		Dial: func() (redis.Conn, error) {
			timeout := time.Duration(2) * time.Second
			c, err := redis.DialTimeout("tcp", server, timeout, 0, 0)
			if err != nil {
				return nil, err
			}
			if len(password) > 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if db > 0 && db < 16 {
				if _, err := c.Do("SELECT", db); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
	}
}

type AuthReq struct {
	Uid         int64  `json:"uid"`
	Device_id   string `json:"device_id"`
	User_name   string `json:"user_name"`
	Platform_id int    `json:"platform_id"`
}

type GrantResp struct {
	Data AuthResp `json:"data"`
}

type AuthResp struct {
	Token string `json:"token"`
}

func token(w http.ResponseWriter, req *http.Request) {
	conn := redis_pool.Get()
	defer conn.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal("request error")
		return
	}

	var authReq AuthReq
	json.Unmarshal(body, &authReq)

	log.Info("auth Req %s", authReq)

	// token
	has := md5.Sum([]byte(authReq.Device_id + authReq.User_name))
	token := fmt.Sprintf("%x", has)

	conn.Do("HMSET", fmt.Sprintf("access_token_%s", token),
		"app_id", authReq.Platform_id,
		"user_id", authReq.Uid,
		"notification_on", "1",
		"forbidden", "0")

	w.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(AuthResp{Token: token})

	w.Write(bytes)
}



func grant(w http.ResponseWriter, req *http.Request) {
	conn := redis_pool.Get()
	defer conn.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal("request error")
		return
	}

	var authReq AuthReq
	json.Unmarshal(body, &authReq)

	log.Info("auth Req %s", authReq)

	// token
	has := md5.Sum([]byte(string(int64(time.Now().Nanosecond()) + authReq.Uid)))
	token := fmt.Sprintf("%x", has)

	conn.Do("HMSET", fmt.Sprintf("access_token_%s", token),
		"app_id", authReq.Platform_id,
		"user_id", authReq.Uid,
		"notification_on", "1",
		"forbidden", "0")

	w.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(GrantResp{Data: AuthResp{Token: token}})

	w.Write(bytes)
}

func main() {
	fmt.Printf("Version:     %s\nBuilt:       %s\nGo version:  %s\nGit branch:  %s\nGit commit:  %s\n", VERSION, BUILD_TIME, GO_VERSION, GIT_BRANCH, GIT_COMMIT_ID)
	rand.Seed(time.Now().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())

	var configFile = "./auth.properties"
	flag.Parse()
	if len(flag.Args()) == 0 {
		_, err := os.Stat(configFile)
		if err != nil {
			fmt.Println("usage: auth auth.properties")
			return
		}
	} else {
		configFile = flag.Args()[0]
	}

	config = read_cfg(configFile)

	redis_pool = NewRedisPool(config.redis_address, config.redis_password, config.redis_db)

	log.Info("auth service start")

	http.HandleFunc("/auth/token", token)
	http.HandleFunc("/auth/grant", grant)
	http.ListenAndServe(fmt.Sprintf(":%d", config.port), nil)
}
