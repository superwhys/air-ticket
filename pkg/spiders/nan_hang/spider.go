package nanhang

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/superwhys/air-ticket/internal/domains"
	"github.com/superwhys/air-ticket/models"
	"github.com/superwhys/air-ticket/pkg/spiders"
	"github.com/tidwall/gjson"
)

type NanHangSpider struct {
	spiders.BaseSpider
}

func (s *NanHangSpider) getHeader(date, from, to string) map[string]string {
	return map[string]string{
		"Accept":             "application/json, text/javascript, */*; q=0.01",
		"Accept-Language":    "zh-CN,zh;q=0.9",
		"Cache-Control":      "no-cache",
		"Content-Type":       "application/json",
		"Origin":             fmt.Sprintf("%v://%v", s.Url.Scheme, s.Url.Host),
		"Pragma":             "no-cache",
		"Priority":           "u=1, i",
		"Referer":            fmt.Sprintf("https://b2c.csair.com/B2C40/newTrips/static/main/page/booking/index.html?t=S&c1=%s&c2=%s&d1=%s&at=2&ct=0&it=0", from, to, date),
		"Sec-Ch-Ua":          "\"Not/A)Brand\";v=\"8\", \"Chromium\";v=\"126\", \"Google Chrome\";v=\"126\"",
		"Sec-Ch-Ua-Mobile":   "?0",
		"Sec-Ch-Ua-Platform": "\"macOS\"",
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "same-origin",
		"User-Agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
		"X-Requested-With":   "XMLHttpRequest",
	}
}

func (s *NanHangSpider) getParams() url.Values {
	params := url.Values{}
	params.Add("type__1188", "CqAxRD0QDteWuD0534+O7fGCD97bkteF4D")
	return params
}

func (s *NanHangSpider) getJsonData(date, from, to string) map[string]any {
	return map[string]interface{}{
		"depCity":       from,
		"arrCity":       to,
		"flightDate":    strings.ReplaceAll(date, "-", ""),
		"adultNum":      "2",
		"childNum":      "0",
		"infantNum":     "0",
		"cabinOrder":    "0",
		"airLine":       1,
		"flyType":       2,
		"international": "0",
		"action":        "0",
		"segType":       "1",
		"cache":         0,
		"preUrl":        "",
		"tariffRules": []map[string]string{
			{
				"bigCustomer":         "",
				"domesticFlightGroup": "ZSKQ37",
				"isLowPriceRation":    "0",
				"flightLine":          fmt.Sprintf("%v-%v", from, to),
			},
		},
		"isMember": "",
	}
}

func (s *NanHangSpider) ParseResp(resp []byte, rule *domains.CrawlRule) ([]*models.AirTicket, error) {
	// ["data"]["segment"][0]["dateFlight"]["transitFlight"]
	transitFlight := gjson.GetBytes(resp, "data.segment.0.dateFlight.transitFlight").String()
	flights := []*TransitFlight{}
	if err := json.Unmarshal([]byte(transitFlight), &flights); err != nil {
		return nil, err
	}

	var tickets []*models.AirTicket
	for _, flight := range flights {
		segments := flight.Segments

		lastIdx := len(segments) - 1
		depTime := segments[0].DepTime
		arrTime := segments[lastIdx].ArrTime
		depComp := fmt.Sprintf("%s%s", segments[0].DepDate, segments[0].DepTime)
		arrComp := fmt.Sprintf("%s%s", segments[lastIdx].DepDate, segments[lastIdx].DepTime)
		depDatetime, err := time.Parse("200601021504", depComp)
		arrDatetime, err := time.Parse("200601021504", arrComp)
		if err != nil {
			continue
		}

		if rule.TickerFilter(depDatetime, arrDatetime) {
			continue
		}

		ticket := &models.AirTicket{
			AirCompany: domains.NANHANG.String(),
			FlightNo:   segments[0].FlightNo,
			DepTime:    depTime,
			ArrTime:    arrTime,
		}

		ticket.Duration = fmt.Sprintf("%.2fh", arrDatetime.Sub(depDatetime).Hours())
		ticket.Price = flight.AdultSortPriceTotal
		tickets = append(tickets, ticket)
	}

	return tickets, nil
}

func (s *NanHangSpider) Crawl(ctx context.Context, rule *domains.CrawlRule) ([]byte, error) {
	headers := s.getHeader(rule.Date, rule.From, rule.To)
	data := s.getJsonData(rule.Date, rule.From, rule.To)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "POST", s.Url.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

type NanHangSpiderFactory struct{}

func (f *NanHangSpiderFactory) Source() domains.SpiderSource {
	return domains.NANHANG
}

func (f *NanHangSpiderFactory) NewAirTickerSpider() domains.AirTicketSpider {
	base, err := url.Parse("https://b2c.csair.com")
	if err != nil {
		panic(err)
	}
	base.Path = "/portal/main/flight/czaddon/query"

	return &NanHangSpider{
		BaseSpider: spiders.BaseSpider{
			Url: base,
		},
	}
}
