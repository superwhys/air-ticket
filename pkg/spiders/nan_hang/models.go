package nanhang

// Response is the top-level structure
type Response struct {
	Success bool   `json:"success"`
	Data    Data   `json:"data"`
	Version string `json:"version"`
	Server  string `json:"server"`
}

// Data contains the nested data
type Data struct {
	ID         string    `json:"id"`
	CreateTime string    `json:"createTime"`
	Segment    []Segment `json:"segment"`
	Citys      []City    `json:"citys"`
	Airports   []Airport `json:"airports"`
	Planes     []Plane   `json:"planes"`
}

// Segment represents a segment of the flight
type Segment struct {
	DepCity    string     `json:"depCity"`
	ArrCity    string     `json:"arrCity"`
	Date       string     `json:"date"`
	DateFlight DateFlight `json:"dateFlight"`
}

// DateFlight contains the details of the flights on a particular date
type DateFlight struct {
	TransitFlight []TransitFlight `json:"transitFlight"`
}

// TransitFlight contains the transit flights
type TransitFlight struct {
	DepCity             string          `json:"depCity"`
	ArrCity             string          `json:"arrCity"`
	Solutions           []Solution      `json:"solutions"`
	Segments            []FlightSegment `json:"segments"`
	AdultSortPriceTotal float64         `json:"adultSortPriceTotal"`
}

// Solution contains the solutions
type Solution struct {
	AdultFareTotal   float64 `json:"adultFareTotal"`
	ChildFareTotal   float64 `json:"childFareTotal"`
	InfantFareTotal  float64 `json:"infantFareTotal"`
	AdultTaxTotal    float64 `json:"adultTaxTotal"`
	ChildTaxTotal    float64 `json:"childTaxTotal"`
	InfantTaxTotal   float64 `json:"infantTaxTotal"`
	AdultPriceTotal  float64 `json:"adultPriceTotal"`
	ChildPriceTotal  float64 `json:"childPriceTotal"`
	InfantPriceTotal float64 `json:"infantPriceTotal"`
	BrandType        string  `json:"brandType"`
	Fares            []Fare  `json:"fares"`
	MinSolutionFare  bool    `json:"minSolutionFare"`
}

// Fare contains the fare details
type Fare struct {
	Name                       string        `json:"name"`
	AdultPrice                 float64       `json:"adultPrice"`
	ChildPrice                 float64       `json:"childPrice"`
	InfantPrice                float64       `json:"infantPrice"`
	Discount                   string        `json:"discount"`
	AdultFareBasis             string        `json:"adultFareBasis"`
	ChildFareBasis             string        `json:"childFareBasis"`
	InfantFareBasis            string        `json:"infantFareBasis"`
	FareReference              string        `json:"fareReference"`
	ChildFareReference         string        `json:"childFareReference"`
	InfantFareReference        string        `json:"infantFareReference"`
	BrandType                  string        `json:"brandType"`
	Code                       string        `json:"code"`
	AdultFareCode              string        `json:"adultFareCode"`
	AdultFareRule              string        `json:"adultFareRule"`
	ChildFareCode              string        `json:"childFareCode"`
	ChildFareRule              string        `json:"childFareRule"`
	InfantFareCode             string        `json:"infantFareCode"`
	InfantFareRule             string        `json:"infantFareRule"`
	AdultBaggageAllowance      string        `json:"adultbaggageallowance"`
	AdultBaggageAllowanceUnit  string        `json:"adultbaggageallowanceunit"`
	ChildBaggageAllowance      string        `json:"childbaggageallowance"`
	ChildBaggageAllowanceUnit  string        `json:"childbaggageallowanceunit"`
	InfantBaggageAllowance     string        `json:"infantbaggageallowance"`
	InfantBaggageAllowanceUnit string        `json:"infantbaggageallowanceunit"`
	Segments                   []SegmentInfo `json:"segments"`
}

// SegmentInfo contains info about the segment
type SegmentInfo struct {
	Name string `json:"name"`
	Info string `json:"info"`
}

// FlightSegment contains the details of a flight segment
type FlightSegment struct {
	FlightNo           string     `json:"flightNo"`
	AirLine            string     `json:"airLine"`
	CodeShare          string     `json:"codeShare"`
	CodeShareInfo      string     `json:"codeShareInfo"`
	DepPort            string     `json:"depPort"`
	ArrPort            string     `json:"arrPort"`
	DepTime            string     `json:"depTime"`
	ArrTime            string     `json:"arrTime"`
	DepDate            string     `json:"depDate"`
	ArrDate            string     `json:"arrDate"`
	TimeDuringFlight   string     `json:"timeDuringFlight"`
	TimeDuringFlightEn string     `json:"timeDuringFlightEn"`
	Plane              string     `json:"plane"`
	StopNumber         string     `json:"stopNumber"`
	StopNameZh         string     `json:"stopNameZh"`
	StopNameEn         string     `json:"stopNameEn"`
	Meal               string     `json:"meal"`
	Term               string     `json:"term"`
	Rate               string     `json:"rate"`
	DepartureTerminal  string     `json:"departureTerminal"`
	ArrivalTerminal    string     `json:"arrivalTerminal"`
	Taxes              []Tax      `json:"taxes"`
	MealInfos          []MealInfo `json:"mealInfos,omitempty"`
}

// Tax contains the tax details
type Tax struct {
	Code   string  `json:"code"`
	Adult  float64 `json:"adult"`
	Child  float64 `json:"child"`
	Infant float64 `json:"infant"`
	GMP    float64 `json:"gmp"`
	JCP    float64 `json:"jcp"`
}

// MealInfo contains meal information
type MealInfo struct {
	Dep   string `json:"dep"`
	Arr   string `json:"arr"`
	Meals []Meal `json:"meals"`
}

// Meal contains meal details
type Meal struct {
	Cabin      string `json:"cabin"`
	MealNameZh string `json:"mealNameZh"`
	MealNameEn string `json:"mealNameEn"`
}

// City represents a city
type City struct {
	Code   string `json:"code"`
	ZhName string `json:"zhName"`
	EnName string `json:"enName"`
}

// Airport represents an airport
type Airport struct {
	Code   string `json:"code"`
	ZhName string `json:"zhName"`
	EnName string `json:"enName"`
	City   string `json:"city"`
}

// Plane represents a plane
type Plane struct {
	Code       string `json:"code"`
	ZhName     string `json:"zhName"`
	EnName     string `json:"enName"`
	AirportTax string `json:"airportTax"`
}
