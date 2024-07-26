package wingworld

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/go-puzzles/plog"
	"github.com/superwhys/air-ticket/internal/domains"
	"github.com/superwhys/air-ticket/models"
	"github.com/superwhys/air-ticket/pkg/spiders"
	"github.com/tidwall/gjson"
)

type WingWorldSpider struct {
	spiders.BaseSpider
}

func (s *WingWorldSpider) getHeader() map[string]string {
	return map[string]string{
		"Accept":             "application/json, text/plain, */*",
		"Accept-Language":    "zh-CN,zh;q=0.9",
		"Cache-Control":      "no-cache",
		"Connection":         "keep-alive",
		"Content-Type":       "application/json",
		"Origin":             "https://m.wingworld.cn",
		"Pragma":             "no-cache",
		"Referer":            "https://m.wingworld.cn/",
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "same-site",
		"User-Agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
		"sec-ch-ua":          `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": `"macOS"`,
	}
}

func (s *WingWorldSpider) getData(date, from, to string) map[string]string {
	return map[string]string{
		"departCityCode": from,
		"arriveCityCode": to,
		"departDate":     date,
		"cabin":          "",
	}
}

func (s *WingWorldSpider) ParseResp(resp []byte, rule *domains.CrawlRule) ([]*models.AirTicket, error) {
	// ["data"]["flightlist"]
	flightList := gjson.GetBytes(resp, "data.flightlist").String()
	flights := []*Flight{}

	if err := json.Unmarshal([]byte(flightList), &flights); err != nil {
		return nil, err
	}

	var tickets []*models.AirTicket
	for _, flight := range flights {
		depComp := fmt.Sprintf("%s %s", flight.DepartDate, flight.DepartTime)

		if flight.ArriveDate == nil {
			flight.ArriveDate = flight.DepartDate
		}
		arrComp := fmt.Sprintf("%s %s", flight.ArriveDate, flight.ArriveTime)

		depDatetime, err := time.Parse("2006-01-02 15:04", depComp)
		arrDatetime, err := time.Parse("2006-01-02 15:04", arrComp)
		if err != nil {
			plog.Errorf("time parse error: %v", err)
			continue
		}

		if rule.TickerFilter(depDatetime, arrDatetime) {
			continue
		}

		ticket := &models.AirTicket{
			AirCompany: flight.AirlineName,
			FlightNo:   flight.FlightNo,
			DepTime:    flight.DepartTime,
			ArrTime:    flight.ArriveTime,
			Duration:   flight.Duration,
			Price: func() float64 {
				p := flight.LowCabin.CabinPrice
				if p == nil {
					return 0
				}

				price, ok := p.(float64)
				if !ok {
					return 0
				}
				return price
			}(),
		}

		tickets = append(tickets, ticket)
	}

	return tickets, nil
}

func (s *WingWorldSpider) Crawl(ctx context.Context, rule *domains.CrawlRule) ([]byte, error) {
	headers := s.getHeader()
	data := s.getData(rule.Date, rule.From, rule.To)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.Url.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

type WingWorldSpiderFactory struct{}

func (f *WingWorldSpiderFactory) Source() domains.SpiderSource {
	return domains.WORLDWING
}

func (f *WingWorldSpiderFactory) NewAirTickerSpider() domains.AirTicketSpider {
	base, err := url.Parse("https://gateway.wingworld.cn")
	if err != nil {
		panic(err)
	}
	base.Path = "/flightapi/flight/flightList"

	return &WingWorldSpider{
		BaseSpider: spiders.BaseSpider{
			Url: base,
		},
	}
}
