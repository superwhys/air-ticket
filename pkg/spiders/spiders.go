package spiders

import (
	"context"
	"net/url"

	"github.com/superwhys/air-ticket/internal/domains"
	"github.com/superwhys/air-ticket/models"
	"github.com/superwhys/air-ticket/pkg/errors"
)

type SpiderFactory map[domains.SpiderSource]domains.SpiderFactory

type BaseSpider struct {
	Url *url.URL
}

var (
	spiders = make(SpiderFactory)
)

func RegisterSpider(sf ...domains.SpiderFactory) SpiderFactory {
	for _, s := range sf {
		spiders[s.Source()] = s
	}

	return spiders
}

func (s SpiderFactory) RegisterSpider(sf domains.SpiderFactory) SpiderFactory {
	s[sf.Source()] = sf
	return s
}

func (s SpiderFactory) GetSpiderFactory(airCompany domains.SpiderSource) (domains.SpiderFactory, error) {
	SpiderFactory, ok := s[airCompany]
	if !ok {
		return nil, errors.ErrUnknownAirCompany
	}
	return SpiderFactory, nil
}

func (s SpiderFactory) Crawl(ctx context.Context, ac domains.SpiderSource, rule *domains.CrawlRule) ([]*models.AirTicket, error) {
	sf, err := s.GetSpiderFactory(ac)
	if err != nil {
		return nil, err
	}

	spider := sf.NewAirTickerSpider()
	resp, err := spider.Crawl(ctx, rule)
	if err != nil {
		return nil, err
	}

	return spider.ParseResp(resp, rule)
}
