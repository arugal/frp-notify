package ip

type BaiduIPAddrResponse struct {
	Status       string            `json:"status"`
	T            string            `json:"t,omitempty"`
	SetCacheTime string            `json:"set_cache_time,omitempty"`
	Data         []BaiduIPAddrData `json:"data"`
}

type BaiduIPAddrData struct {
	ExtendedLocation string `json:"ExtendedLocation,omitempty"`
	OriginQuery      string `json:"OriginQuery,omitempty"`
	AppInfo          string `json:"appinfo,omitempty"`
	DispType         int    `json:"disp_type,omitempty"`
	FetchKey         string `json:"fetchkey,omitempty"`
	Location         string `json:"location"`
	OrigIP           string `json:"origip,omitempty"`
	OrigIPQuery      string `json:"origipquery,omitempty"`
	ResourceID       string `json:"resourceid,omitempty"`
	RoleID           int    `json:"role_id,omitempty"`
	ShareImage       int    `json:"shareImage,omitempty"`
	ShowLikeShare    int    `json:"showLikeShare,omitempty"`
	ShowLamp         string `json:"showlamp,omitempty"`
	TitleCont        string `json:"titlecont,omitempty"`
	Tplt             string `json:"tplt,omitempty"`
}
