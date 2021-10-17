package teslaweb

type CountryCode string

const (
	Us CountryCode = "US"
)

type CurrencyCode string

const (
	Usd CurrencyCode = "USD"
)

type Description string

const (
	ChargingFreeSite Description = "CHARGING : FREE_SITE"
	ChargingPayment  Description = "CHARGING : PAYMENT"
	ParkingNoCharge  Description = "PARKING : NO_CHARGE"
)

type FeeType string

const (
	Charging FeeType = "CHARGING"
	Parking  FeeType = "PARKING"
)

type Model string

const (
	Model3 Model = "m3"
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

type Type string

const (
	Immediate Type = "IMMEDIATE"
)

type Uom string

const (
	Free Uom = "free"
	Kwh  Uom = "kwh"
	Min  Uom = "min"
)

type VehicleMakeType string

const (
	Tesla VehicleMakeType = "TSLA"
)
