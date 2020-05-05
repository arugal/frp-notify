// Copyright 2020 arugal, zhangwei24@apache.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ip

import (
	"fmt"
	"github/arugal/frp-notify/pkg/logger"
	"strings"

	"github.com/axgle/mahonia"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
)

func init() {
	log = logger.Log
}

// 查询 ip 实际地址
type AddressQuery func(ip string) string

func NewDefaultAddressService() AddressService {
	as := &defaultAddressService{
		client: resty.New(),
	}
	return as
}

type AddressService interface {
	Query(ip string) string
}

// use 126 api
type defaultAddressService struct {
	client *resty.Client
}

func (s *defaultAddressService) Query(ip string) string {
	req := s.client.R()
	req.Header.Add("Content-Type", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.75 Safari/537.36")

	resp, err := req.Get(fmt.Sprintf("http://ip.ws.126.net/ipquery?ip=%s", ip))
	if err != nil {
		log.Errorf("ip query error, err: %v", err)
		return ""
	}
	result := convertToString(string(resp.Body()), "GBK", "UTF-8")
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

func convertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}
