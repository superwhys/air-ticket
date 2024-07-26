package main

import (
	"github.com/go-puzzles/cores"
	"github.com/go-puzzles/pflags"
	"github.com/go-puzzles/plog"
	"github.com/superwhys/air-ticket/configs"
	"github.com/superwhys/air-ticket/pkg/email"
	"github.com/superwhys/air-ticket/pkg/spiders"
	nanhang "github.com/superwhys/air-ticket/pkg/spiders/nan_hang"
	"github.com/superwhys/air-ticket/pkg/spiders/wingworld"
	"github.com/superwhys/air-ticket/pkg/worker"
)

var (
	crawlCron       = pflags.String("crawl-cron", "0 */1 * * *", "Crawl cron")
	crawlRulesFlags = pflags.Struct("crawl-rules", []*configs.Rule{}, "A list of rules to crawl")
	emailConfFlags  = pflags.Struct("email-conf", &email.EmailConf{}, "Email configuration")
	emailTargets    = pflags.StringSlice("email-targets", []string{}, "Email send targets")
)

func main() {
	pflags.Parse()
	var rules []*configs.Rule
	plog.PanicError(crawlRulesFlags(&rules))

	emailConf := new(email.EmailConf)
	plog.PanicError(emailConfFlags(emailConf))

	plog.Infof("Crawling rules: %v", plog.JsonifyNoIndent(rules))

	gmailSender := email.NewGmailSender(emailConf)
	spiderFactory := spiders.RegisterSpider(
		&nanhang.NanHangSpiderFactory{},
		&wingworld.WingWorldSpiderFactory{},
	)

	spiderOpt := worker.NewOptions(
		spiderFactory,
		rules,
		gmailSender,
		emailTargets(),
	)
	spiderWorker := worker.NewSpiderWorker(spiderOpt)

	c := cores.NewPuzzleCore(
		cores.WithWorker(spiderWorker.Run),
		// cores.WithCronWorker(crawlCron(), spiderWorker.Run),
	)

	plog.PanicError(cores.Run(c))
}
