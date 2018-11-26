package main

import (
	"fmt"
	"time"
)
import "net"
import "log"
import "runtime"
import "flag"
import "math/rand"
import "net/http"
import "encoding/base64"
import "crypto/md5"
import "strings"
import "encoding/json"
import "github.com/bitly/go-simplejson"
import "io/ioutil"

var first int64
var last int64
var host string
var port int

const APP_ID = 7
const APP_KEY = "sVDIlIiDUm7tWPYWhi6kfNbrqui3ez44"
const APP_SECRET = "0WiCxAU1jh76SbgaaFC7qIaBPm2zkyM1"
const URL = "http://127.0.0.1"

func init() {
	flag.Int64Var(&first, "first", 0, "first uid")
	flag.Int64Var(&last, "last", 0, "last uid")

	flag.StringVar(&host, "host", "127.0.0.1", "host")
	flag.IntVar(&port, "port", 23000, "port")
}

func login(uid int64) string {
	url := URL + "/auth/grant"
	secret := fmt.Sprintf("%x", md5.Sum([]byte(APP_SECRET)))
	s := fmt.Sprintf("%d:%s", APP_ID, secret)
	basic := base64.StdEncoding.EncodeToString([]byte(s))

	v := make(map[string]interface{})
	v["uid"] = uid

	body, _ := json.Marshal(v)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(body)))
	req.Header.Set("Authorization", "Basic " + basic)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	res, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer res.Body.Close()
	
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ""
	}
	obj, err := simplejson.NewJson(b)
	token, _ := obj.Get("data").Get("token").String()
	return token
}


func send(uid int64, count int) {
	ip := net.ParseIP(host)
	addr := net.TCPAddr{ip, port, ""}

	token := login(uid)
	
	conn, err := net.DialTCP("tcp4", nil, &addr)
	defer conn.Close()
	if err != nil {
		log.Println("connect error")
		return
	}
	seq := 1

	auth := &AuthenticationToken{token:token, platform_id:1, device_id:"00000000"}
	SendMessage(conn, &Message{cmd:MSG_AUTH_TOKEN, seq:seq, version:DEFAULT_VERSION, body:auth})	
	ReceiveMessage(conn)

	for i := 0; i < count; i++ {
		r := rand.Int63()
		// 随机一个接收方 id
		receiver := r%(last-first) + first
		if receiver == uid {
			continue
		}
		// log.Println("receiver:", receiver)
		content := fmt.Sprintf("test....%d", i)
		seq++
		msg := &Message{MSG_IM, seq, DEFAULT_VERSION, 0, &IMMessage{uid, receiver, 0, int32(i), content}}
		// 给对方发一个消息
		SendMessage(conn, msg)
		// 等待确认
		for {
			ack := ReceiveMessage(conn)
			if ack.cmd == MSG_ACK {
				break
			}
		}
	}
	log.Printf("%d send complete", uid)
}

// ./benchmark_sender -first=1 -last=30
func main() {
	runtime.GOMAXPROCS(4)
	flag.Parse()
	fmt.Printf("first:%d, last:%d\n", first, last)
	if last <= first {
		return
	}
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	begin := time.Now().UnixNano()
	log.Println("begin test:", begin)

	// logic
	count := 100000
	send(1, count)


	end := time.Now().UnixNano()
	var tps int64 = 0
	if end-begin > 0 {
		tps = int64(1000*1000*1000*count) / (end - begin)
	}
	fmt.Println("tps:", tps)
}
