package request

type AdvertiseReq struct {
	Id     int32  `json:"id"`
	Index  int32  `json:"index"`
	Images string `json:"images"`
	Url    string `json:"url"`
}
