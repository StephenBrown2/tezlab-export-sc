package teslaweb

// https://www.tesla.com/teslaaccount/charging/api/export?startTime=2019-09-01T04:00:00.000Z&endTime=2021-10-01T04:00:00.000Z&vin=

type Export struct {
	Code    int64  `json:"code"`
	Data    Data   `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type Data struct {
	Data string `json:"data"` // base64-encoded csv data
}

type ExportCSV []ExportCSVElement

type ExportCSVElement struct {
	ChargeStartDateTime string      `json:"ChargeStartDateTime"`
	Vin                 string      `json:"Vin"`
	OwnerName           string      `json:"Name"`
	Model               Model       `json:"Model"`
	Country             Country     `json:"Country"`
	SiteLocationName    string      `json:"SiteLocationName"`
	Description         Description `json:"Description"`
	QuantityBase        string      `json:"QuantityBase"`
	QuantityTier1       string      `json:"QuantityTier1"`
	QuantityTier2       string      `json:"QuantityTier2"`
	InvoiceNumber       string      `json:"InvoiceNumber"`
	UnitCostBase        string      `json:"UnitCostBase"`  // i.e. "0.20/kwh", "1.00/min", "N/A"
	UnitCostTier1       string      `json:"UnitCostTier1"` // i.e. "0.11/min", "N/A"
	UnitCostTier2       string      `json:"UnitCostTier2"` // i.e. "0.22/min", "N/A"
	Vat                 float64     `json:"VAT"`
	TotalExcVAT         float64     `json:"Total Exc. VAT"`
	TotalIncVAT         float64     `json:"Total Inc. VAT"`
	Status              Status      `json:"Status"`
	Invoice             string      `json:"Invoice"`
}

type Country string

const (
	Us Country = "US"
)

type Description string

const (
	ChargingFreeSite Description = "CHARGING : FREE_SITE"
	ChargingPayment  Description = "CHARGING : PAYMENT"
	ParkingNoCharge  Description = "PARKING : NO_CHARGE"
)

type Model string

const (
	M3 Model = "m3"
)

type Status string

const (
	Paid Status = "PAID"
)
