package tezlab

import (
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/goinggo/timezone"
	log "github.com/sirupsen/logrus"
)

type ChargeRecord struct {
	VehicleName          string `csv:"Vehicle Name"`
	Vin                  string `csv:"VIN"`
	Timezone             string `csv:"Timezone"`
	StartTime            string `csv:"Start Time"`
	EndTime              string `csv:"End Time"`
	OdometerMi           string `csv:"Odometer (mi)"`
	PowerDrawnKWh        string `csv:"Power Drawn (kWh)"`
	ChargeEnergyAddedKWh string `csv:"Charge Energy Added (kWh)"`
	StartRangeMi         string `csv:"Start Range (mi)"`
	EndRangeMi           string `csv:"End Range (mi)"`
	RangeAddedMi         string `csv:"Range Added (mi)"`
	DurationS            string `csv:"Duration (s)"`
	Supercharging        string `csv:"Supercharging"`
	Supercharger         string `csv:"Supercharger"`
	MaxChargerPowerKW    string `csv:"Max Charger Power (kW)"`
	FastChargerPresent   string `csv:"Fast Charger Present"`
	ConnectorType        string `csv:"Connector Type"`
	Location             string `csv:"Location"`
	Coordinate           string `csv:"Coordinate"`
	UserChargeLocation   string `csv:"User Charge Location"`
	CostUSD              string `csv:"Cost (USD)"`
	TemperatureF         string `csv:"Temperature (F)"`
	CO2G                 string `csv:"CO2 (g)"`
	GCO2KWh              string `csv:"gCO2/kWh"`
	EqFuelBurnGal        string `csv:"Eq. Fuel Burn (gal)"`
}

func NewRecordFromCSV(row []string) ChargeRecord {
	return ChargeRecord{
		VehicleName:          row[0],
		Vin:                  row[1],
		Timezone:             row[2],
		StartTime:            row[3],
		EndTime:              row[4],
		OdometerMi:           row[5],
		PowerDrawnKWh:        row[6],
		ChargeEnergyAddedKWh: row[7],
		StartRangeMi:         row[8],
		EndRangeMi:           row[9],
		RangeAddedMi:         row[10],
		DurationS:            row[11],
		Supercharging:        row[12],
		Supercharger:         row[13],
		MaxChargerPowerKW:    row[14],
		FastChargerPresent:   row[15],
		ConnectorType:        row[16],
		Location:             row[17],
		Coordinate:           row[18],
		UserChargeLocation:   row[19],
		CostUSD:              row[20],
		TemperatureF:         row[21],
		CO2G:                 row[22],
		GCO2KWh:              row[23],
		EqFuelBurnGal:        row[24],
	}
}

func (r ChargeRecord) IsSupercharger() bool {
	log.Debugf("record.Supercharging == %q", r.Supercharging)
	log.Debugf("record.UserChargeLocation == %q (%t)", r.UserChargeLocation, strings.Contains(r.UserChargeLocation, "Super"))

	return r.Supercharging == "true" || strings.Contains(r.UserChargeLocation, "Super")
}

func (r ChargeRecord) SuperchargerLocation() string {
	scLocation := r.Supercharger

	if scLocation == "" { // User-entered Supercharger
		scLocation = r.Location
	}

	log.Debugf("scLocation: %q", scLocation)

	if !strings.Contains(scLocation, r.Location[len(r.Location)-4:]) && scLocation[len(scLocation)-4] != ',' {
		scLocation += r.Location[len(r.Location)-4:]
		log.Debugf("scLocation new: %q", scLocation)
	}

	if strings.HasSuffix(scLocation, "COs") {
		fmt.Println(r)
		scLocation = strings.TrimRight(scLocation, "s")
		log.Debugf("scLocation trimmed: %q", scLocation)
	}

	return scLocation
}

func (r ChargeRecord) StartTimeInRecordTimezone() time.Time {
	log.Debugf("Loading record Timezone: %q", r.Timezone)

	recordLoc, err := time.LoadLocation(r.Timezone)
	if err != nil {
		log.Errorf("Could not parse charge location timezone: %q - %s\n", r.Timezone, err)
		os.Exit(1)
	}

	log.Debugf("Parsing time %q in Location %q", r.StartTime, recordLoc)

	chargeDate, err := time.ParseInLocation("2006-01-02 03:04PM", r.StartTime, recordLoc)
	if err != nil {
		log.Errorf("Could not parse charge start time: %q - %s\n", r.StartTime, err)
		os.Exit(1)
	}

	return chargeDate
}

func (r ChargeRecord) LocationTimezone() *time.Location {
	coord := strings.Split(r.Coordinate, ",")
	log.Debugf("Coordinates: %v", coord)

	lat, err := strconv.ParseFloat(coord[0], 64)
	if err != nil {
		log.Errorf("Could not parse charge location latitude as float: %q - %s\n", coord[0], err)
		os.Exit(1)
	}

	lon, err := strconv.ParseFloat(coord[1], 64)
	if err != nil {
		log.Errorf("Could not parse charge location longitude as float: %q - %s\n", coord[1], err)
		os.Exit(1)
	}

	log.Debug("Retrieving Timezone from GeoNames")

	gtz, err := timezone.RetrieveGeoNamesTimezone(lat, lon, username())
	if err != nil {
		log.Errorf("Could not retrieve Timezone from GeoNames: %s\n", err)
		os.Exit(1)
	}

	log.Debugf("Loading location from %q", gtz.TimezoneID)

	location, err := time.LoadLocation(gtz.TimezoneID)
	if err != nil {
		log.Errorf("Could not parse TimezoneID: %q - %s\n", gtz.TimezoneID, err)
		os.Exit(1)
	}

	return location
}

func username() string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		panic("Can't read BuildInfo")
	}

	s := strings.Split(buildInfo.Path, "/")

	return s[1]
}
