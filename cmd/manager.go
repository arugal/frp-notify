package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github/arugal/frp-manager/logger"
	"github/arugal/frp-manager/models"
	"io/ioutil"
	"net/http"
	"strings"
)

var resp = models.Response{
	Reject:       false,
	RejectReason: "",
	Unchange:     true,
	Content:      nil,
}

var handlerChan = make(chan models.Request, 10)

var config = Config{
	Addr: ":80",
}
var log *logrus.Logger

func init() {
	log = logger.Log
}

func gotify(title string, message string) {
	if config.GotifyAddress == "" || config.GotifyToken == "" {
		return
	}

	client := resty.New()

	format := make(map[string]string)
	format["title"] = title
	format["message"] = message

	_, err := client.R().
		SetFormData(format).
		Post(fmt.Sprintf("http://%s/message?token=%s", config.GotifyAddress, config.GotifyToken))
	if err != nil {
		log.Errorf("gotify error, err:%s address:%s, token:%s", err, config.GotifyAddress, config.GotifyToken)
	}
}

func ipQuery(ip string) string {
	client := resty.New()

	request := client.R()
	request.Header.Add("Content-Type", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.75 Safari/537.36")

	resp, err := request.Get(fmt.Sprintf("http://ip.ws.126.net/ipquery?ip=%s", ip))
	if err != nil {
		log.Errorf("ip query error, err: %v", err)
		return ""
	}
	result := models.ConvertToString(string(resp.Body()), "GBK", "UTF-8")
	return provinceAndCity(result)
}

func provinceAndCity(result string) string {
	i := strings.Index(result, "city:\"") + 6
	result = result[i:]
	i = strings.Index(result, "\"")
	city := result[:i]

	i = strings.Index(result, "province:\"") + 10
	result = result[i:]
	i = strings.Index(result, "\"")
	province := result[:i]

	return province + city
}

func normal(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, resp)
}

func main() {
	config.addFlags()
	flag.Parse()

	r := gin.Default()

	r.POST("/handler", func(ctx *gin.Context) {
		body, err := ioutil.ReadAll(ctx.Request.Body)
		defer normal(ctx)
		if err != nil {
			log.Errorf("read body error, err: %v", err)
			return
		}
		requst := models.Request{}
		err = json.Unmarshal(body, &requst)
		if err != nil {
			log.Errorf("json unmarshal error, err: %v", err)
			return
		}
		handlerChan <- requst
	})

	go func() {
		for requst := range handlerChan {
			if requst.Version != models.APIVersion {
				log.Errorf("Unsupported api version %s", requst.Version)
				return
			}
			log.Infof("%s - %v", requst.Op, requst.Content)
			content, ok := requst.Content.(map[string]interface{})
			if !ok {
				log.Errorf("Unsupported content %v", requst.Content)
				return
			}
			switch requst.Op {
			case models.OpLogin:
				gotify(models.OpLogin, fmt.Sprintf("Version: %v, HostName: %v, Os: %v, Arch: %v",
					content["version"], content["hostname"], content["os"], content["arch"]))
				break
			case models.OpNewProxy:
				gotify(models.OpNewProxy+" - "+fmt.Sprint(content["proxy_name"]),
					fmt.Sprintf("ProxyName: %v, ProxyType: %v, RemotePort: %v",
						content["proxy_name"], content["proxy_type"], content["remote_port"]))
				break
			case models.OpNewAccessIp:
				userRemoteIp := fmt.Sprint(content["user_remote_ip"])
				provinceAndCity := ipQuery(userRemoteIp)
				gotify(models.OpNewAccessIp+" - "+fmt.Sprint(content["proxy_name"]),
					fmt.Sprintf("RemotIp: %s - %s", userRemoteIp, provinceAndCity))
			}
		}
	}()

	err := r.Run(config.Addr)
	if err != nil {
		panic(err)
	}
}
