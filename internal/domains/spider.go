package domains

import (
	"context"
	"time"

	"github.com/superwhys/air-ticket/models"
)

type SpiderSource string

const (
	NANHANG   SpiderSource = "南航"
	WORLDWING SpiderSource = "WorldWing"
)

func (c SpiderSource) String() string {
	return string(c)
}

type CrawlRule struct {
	Date string
	From string
	To   string
	// StartTime and EndTime use to filter the ticket which not in StartTime and EndTime range
	StartTime time.Time
	EndTime   time.Time
}

func (r *CrawlRule) TickerFilter(dep, arr time.Time) bool {
	if r.StartTime.After(dep) || r.EndTime.Before(arr) {
		return true
	}

	return false
}

type AirTicketSpider interface {
	Crawl(ctx context.Context, rule *CrawlRule) ([]byte, error)
	ParseResp(resp []byte, rule *CrawlRule) ([]*models.AirTicket, error)
}

type SpiderFactory interface {
	Source() SpiderSource
	NewAirTickerSpider() AirTicketSpider
}
