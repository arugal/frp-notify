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
	"encoding/json"

	"github/arugal/frp-notify/pkg/logger"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
)

func init() {
	log = logger.Log
}

// AddressQuery 查询 ip 实际地址
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
	req.Header.Add(
		"Content-Type",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
			"AppleWebKit/537.36 (KHTML, like Gecko) "+
			"Chrome/112.0.0.0 Safari/537.36",
	)

	resp, err := req.SetQueryParams(map[string]string{
		"query":       ip,
		"co":          "",
		"resource_id": "6006",
		"oe":          "utf8",
	}).SetHeader("Accept", "application/json").Get("https://opendata.baidu.com/api.php")
	if err != nil {
		log.Errorf("ip query: get error, detail: %v", err)
		return ""
	}

	var response BaiduIPAddrResponse
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		log.Errorf("ip query: decode error, detail: %v", err)
		return ""
	}
	if response.Status != "0" {
		log.Errorf("ip query: status error, detail: %v", response.Status)
		return ""
	}

	if len(response.Data) == 0 {
		log.Errorf("ip query: empty data, detail: %v", response.Data)
		return ""
	}
	return response.Data[0].Location
}
