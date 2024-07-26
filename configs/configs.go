package configs

import (
	"time"

	"github.com/superwhys/air-ticket/internal/domains"
)

type Rule struct {
	Source domains.SpiderSource
	Date   string
	From   string
	To     string
	// StartTime and EndTime use to filter the ticket which not in StartTime and EndTime range
	StartTime string
	EndTime   string
}

func (r *Rule) CrawlRule() (*domains.CrawlRule, error) {
	start, err := time.Parse("2006-01-02 15:04", r.StartTime)
	end, err := time.Parse("2006-01-02 15:04", r.EndTime)
	if err != nil {
		return nil, err
	}

	return &domains.CrawlRule{
		Date:      r.Date,
		From:      r.From,
		To:        r.To,
		StartTime: start,
		EndTime:   end,
	}, nil
}
