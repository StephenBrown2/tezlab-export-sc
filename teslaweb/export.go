package teslaweb

import "time"

// https://www.tesla.com/teslaaccount/charging/api/export?startTime=2019-09-01T04:00:00.000Z&endTime=2021-10-01T04:00:00.000Z&vin=

type Export struct {
	Code    int64  `json:"code"`
	Data    Verify `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type Data struct {
	Data string `json:"data"` // base64-encoded csv data
}

type ExportCSV []ExportCSVRow

type ExportCSVRow struct {
	ChargeStartDateTime time.Time   `csv:"ChargeStartDateTime"`
	Vin                 string      `csv:"Vin"`
	OwnerName           string      `csv:"Name"`
	Model               Model       `csv:"Model"`
	Country             CountryCode `csv:"Country"`
	SiteLocationName    string      `csv:"SiteLocationName"`
	Description         Description `csv:"Description"`
	QuantityBase        string      `csv:"QuantityBase"`
	QuantityTier1       string      `csv:"QuantityTier1"`
	QuantityTier2       string      `csv:"QuantityTier2"`
	InvoiceNumber       string      `csv:"InvoiceNumber"`
	UnitCostBase        string      `csv:"UnitCostBase"`  // i.e. "0.20/kwh", "1.00/min", "N/A"
	UnitCostTier1       string      `csv:"UnitCostTier1"` // i.e. "0.11/min", "N/A"
	UnitCostTier2       string      `csv:"UnitCostTier2"` // i.e. "0.22/min", "N/A"
	Vat                 float64     `csv:"VAT"`
	TotalExcVAT         float64     `csv:"Total Exc. VAT"`
	TotalIncVAT         float64     `csv:"Total Inc. VAT"`
	Status              Status      `csv:"Status"`
	Invoice             string      `csv:"Invoice"`
}
