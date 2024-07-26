package domains

import (
	"context"

	"github.com/superwhys/air-ticket/models"
)

type AirCompany string

const (
	NANHANG = "南航"
)

func (c AirCompany) String() string {
	return string(c)
}

var (
	AirCompanies = map[string]AirCompany{
		"南航": NANHANG,
	}
)

type AirTicketSpider interface {
	Crawl(ctx context.Context, date, from, to string) ([]*models.AirTicket, error)
}

type SpiderFactory interface {
	AirCompany() AirCompany
	NewAirTickerSpider() AirTicketSpider
}
