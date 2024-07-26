package worker

import (
	"bytes"
	"context"
	"text/template"

	"github.com/go-puzzles/plog"
	"github.com/pkg/errors"
	"github.com/superwhys/air-ticket/configs"
	"github.com/superwhys/air-ticket/internal/domains"
	"github.com/superwhys/air-ticket/models"
	"github.com/superwhys/air-ticket/pkg/spiders"
	"golang.org/x/sync/errgroup"
)

const (
	sendTmpl = `
Airline: {{.From}} -> {{.To}} 机票检测

机票:
{{range .Tickets}}---------------------------------
航空公司: {{.AirCompany}}
航班号: {{.FlightNo}}
出发时间: {{.DepTime}}
到达时间: {{.ArrTime}}
持续时间: {{.Duration}}
票价: ${{printf "%.2f" .Price}}
{{end}}`
)

type Options struct {
	spiderFactory spiders.SpiderFactory
	crawlRules    []*configs.Rule
	emailSender   domains.EmailSender
	emailTargets  []string
}

func NewOptions(
	sf spiders.SpiderFactory,
	crawlRules []*configs.Rule,
	emailSender domains.EmailSender,
	emailTargets []string,
) *Options {
	return &Options{
		spiderFactory: sf,
		crawlRules:    crawlRules,
		emailSender:   emailSender,
		emailTargets:  emailTargets,
	}
}

type SpiderWorker struct {
	opts *Options
}

func NewSpiderWorker(opts *Options) *SpiderWorker {
	return &SpiderWorker{
		opts: opts,
	}
}

func (s *SpiderWorker) crawlGroup() map[domains.SpiderSource][]*configs.Rule {
	group := make(map[domains.SpiderSource][]*configs.Rule)
	for _, rule := range s.opts.crawlRules {
		group[rule.Source] = append(group[rule.Source], rule)
	}
	return group
}

func (s *SpiderWorker) generateSendMsg(ac domains.SpiderSource, rule *configs.Rule, resp []*models.AirTicket) ([]byte, error) {
	tmpl, err := template.New("test").Parse(sendTmpl)
	if err != nil {
		panic(err)
	}

	data := struct {
		AirCompany string
		From       string
		To         string
		Tickets    []*models.AirTicket
	}{
		AirCompany: ac.String(),
		From:       rule.From,
		To:         rule.To,
		Tickets:    resp,
	}

	var buf bytes.Buffer

	if err = tmpl.Execute(&buf, data); err != nil {
		return nil, errors.Wrap(err, "tmplExecute")
	}

	return buf.Bytes(), nil
}

func (s *SpiderWorker) doCrawl(ctx context.Context, ac domains.SpiderSource, rules []*configs.Rule) func() error {
	return func() error {
		for _, rule := range rules {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			cc := plog.With(ctx, "From", rule.From, "To", rule.To)
			crawlRule, err := rule.CrawlRule()
			if err != nil {
				plog.Errorc(cc, "crawlRule error: %v", err)
				continue
			}

			resp, err := s.opts.spiderFactory.Crawl(cc, ac, crawlRule)
			if err != nil {
				plog.Errorc(cc, "do crawl of %v %v -> %v error: %v", ac, rule.From, rule.To, err)
				continue
			}

			plog.Debugc(cc, "do crawl of %v %v -> %v. resp: %v", ac, rule.From, rule.To, plog.JsonifyNoIndent(resp))

			sendMsg, err := s.generateSendMsg(ac, rule, resp)
			if err != nil {
				plog.Errorc(cc, "generate sendMsg error: %v", err)
				continue
			}

			if err := s.opts.emailSender.SentMsg(cc, s.opts.emailTargets, sendMsg); err != nil {
				plog.Errorc(cc, "send email error: %v", err)
				continue
			}
			plog.Infoc(cc, "do crawl success")
		}

		return nil
	}
}

func (s *SpiderWorker) Run(ctx context.Context) error {
	spiderGroup := s.crawlGroup()

	grp, ctx := errgroup.WithContext(ctx)
	grp.SetLimit(len(spiderGroup))
	for ac, rules := range spiderGroup {
		ac := ac
		rules := rules
		c := plog.With(ctx, "Source", ac.String())
		grp.Go(s.doCrawl(c, ac, rules))
	}

	return grp.Wait()
}
