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

type DistanceUnit string

const (
	Mi DistanceUnit = "mi"
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
	Credit                 PricingType = "CREDIT"
	FirstFreeParking       PricingType = "FIRST_FREE_PARKING"
	FreeSite               PricingType = "FREE_SITE"
	NoCharge               PricingType = "NO_CHARGE"
	ParkingLessThanMinTime PricingType = "PARKING_LESS_THAN_MIN_TIME"
	Payment                PricingType = "PAYMENT"
)

type Status string

const (
	Paid Status = "PAID"
)

type UnitOfMeasurement string

const (
	Empty        UnitOfMeasurement = ""
	Free         UnitOfMeasurement = "free"
	FreeCharging UnitOfMeasurement = "free_charging"
	Kwh          UnitOfMeasurement = "kwh"
	Min          UnitOfMeasurement = "min"
)

type VehicleMakeType string

const (
	Tesla VehicleMakeType = "TSLA"
)
