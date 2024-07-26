package models

type AirTicket struct {
	// 航空公司
	AirCompany string `json:"airCompany"`
	// 航班号
	FlightNo string `json:"flightNo"`
	// 出发时间
	DepTime string `json:"depTime"`
	// 到达时间
	ArrTime string `json:"arrTime"`
	// 飞行时间总计
	Duration string `json:"duration"`
	// 价格
	Price float64 `json:"price"`
}
