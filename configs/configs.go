package configs

import "github.com/superwhys/air-ticket/internal/domains"

type Rule struct {
	AirCompany domains.AirCompany
	Date       string
	From       string
	To         string
}
