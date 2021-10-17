package main

import "time"

type History struct {
	Code    int64    `json:"code"`
	Charges []Charge `json:"data"`
	Message string   `json:"message"`
	Success bool     `json:"success"`
}

type Charge struct {
	Vin                 string          `json:"vin"`
	ChargeSessionID     string          `json:"chargeSessionId"`
	SiteLocationName    string          `json:"siteLocationName"`
	ChargeStartDateTime time.Time       `json:"chargeStartDateTime"`
	ChargeStopDateTime  time.Time       `json:"chargeStopDateTime"`
	UnlatchDateTime     time.Time       `json:"unlatchDateTime"`
	CountryCode         CountryCode     `json:"countryCode"`
	Credit              interface{}     `json:"credit"`
	DisputeDetails      interface{}     `json:"disputeDetails"`
	Fees                []Fee           `json:"fees"`
	BillingType         Type            `json:"billingType"`
	Invoices            []Invoice       `json:"invoices"`
	FapiaoDetails       interface{}     `json:"fapiaoDetails"`
	ProgramType         string          `json:"programType"`
	VehicleMakeType     VehicleMakeType `json:"vehicleMakeType"`
}

type Fee struct {
	SessionFeeID  int64        `json:"sessionFeeId"`
	FeeType       FeeType      `json:"feeType"`
	CurrencyCode  CurrencyCode `json:"currencyCode"`
	PricingType   PricingType  `json:"pricingType"`
	RateBase      float64      `json:"rateBase"`
	RateTier1     float64      `json:"rateTier1"`
	RateTier2     float64      `json:"rateTier2"`
	RateTier3     interface{}  `json:"rateTier3"`
	RateTier4     interface{}  `json:"rateTier4"`
	UsageBase     int64        `json:"usageBase"`
	UsageTier1    int64        `json:"usageTier1"`
	UsageTier2    int64        `json:"usageTier2"`
	UsageTier3    interface{}  `json:"usageTier3"`
	UsageTier4    interface{}  `json:"usageTier4"`
	TotalBase     float64      `json:"totalBase"`
	TotalTier1    float64      `json:"totalTier1"`
	TotalTier2    float64      `json:"totalTier2"`
	TotalTier3    int64        `json:"totalTier3"`
	TotalTier4    int64        `json:"totalTier4"`
	TotalDue      float64      `json:"totalDue"`
	NetDue        float64      `json:"netDue"`
	Uom           Uom          `json:"uom"`
	IsPaid        bool         `json:"isPaid"`
	Status        Status       `json:"status"`
	ProcessFlagID int64        `json:"processFlagId"`
}

type Invoice struct {
	FileName    string `json:"fileName"`
	ContentID   string `json:"contentId"`
	InvoiceType Type   `json:"invoiceType"`
	BeInvoiceID string `json:"beInvoiceId"`
	ProcessFlag int64  `json:"processFlag"`
}

type Type string

const (
	Immediate Type = "IMMEDIATE"
)

type CountryCode string

const (
	Us CountryCode = "US"
)

type CurrencyCode string

const (
	Usd CurrencyCode = "USD"
)

type FeeType string

const (
	Charging FeeType = "CHARGING"
	Parking  FeeType = "PARKING"
)

type PricingType string

const (
	FreeSite PricingType = "FREE_SITE"
	NoCharge PricingType = "NO_CHARGE"
	Payment  PricingType = "PAYMENT"
)

type Status string

const (
	Paid Status = "PAID"
)

type Uom string

const (
	Free Uom = "free"
	Kwh  Uom = "kwh"
	Min  Uom = "min"
)

type VehicleMakeType string

const (
	Tsla VehicleMakeType = "TSLA"
)
