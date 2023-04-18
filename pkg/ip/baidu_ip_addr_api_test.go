package ip

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaiduIpAddrUnmarshal(t *testing.T) {
	jsonString := `{
        "status": "0",
        "t": "",
        "set_cache_time": "",
        "data": [
            {
                "ExtendedLocation": "",
                "OriginQuery": "223.239.128.138",
                "appinfo": "",
                "disp_type": 0,
                "fetchkey": "223.239.128.138",
                "location": "India",
                "origip": "223.239.128.138",
                "origipquery": "223.239.128.138",
                "resourceid": "6006",
                "role_id": 0,
                "shareImage": 1,
                "showLikeShare": 1,
                "showlamp": "1",
                "titlecont": "IP Address Query",
                "tplt": "ip"
            }
        ]
    }`

	var response BaiduIPAddrResponse
	err := json.Unmarshal([]byte(jsonString), &response)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("%+v\n", response)
	t.Logf("ip address: %s", response.Data[0].OrigIPQuery)
	t.Logf("ip location: %s", response.Data[0].Location)
	assert.Equal(t, "223.239.128.138", response.Data[0].OrigIPQuery)
	assert.Equal(t, "India", response.Data[0].Location)
}

func TestBaiduIpAddrUnmarshalWithEmptyData(t *testing.T) {
	jsonString := `{
        "status": "0",
        "t": "",
        "set_cache_time": "",
        "data": []
    }`

	var response BaiduIPAddrResponse
	err := json.Unmarshal([]byte(jsonString), &response)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("%+v\n", response)
	assert.Equal(t, len(response.Data), 0)
}

func TestBaiduIpAddrUnmarshalWithMissingData(t *testing.T) {
	jsonString := `{
        "status": "0",
        "t": "",
        "set_cache_time": ""
    }`

	var response BaiduIPAddrResponse
	err := json.Unmarshal([]byte(jsonString), &response)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("%+v\n", response)
	assert.Equal(t, len(response.Data), 0)
}
