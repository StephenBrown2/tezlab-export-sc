package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	var (
		after     string
		afterDate time.Time
		err       error
	)

	flag.StringVar(&after, "after", "", "Date after which to display new supercharger visits")
	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatalln("Must provide csv file as positional argument")
	}

	if after != "" {
		var t time.Time
		for _, format := range []string{"2006-01-02", "01/02/2006", "1/2/2006"} {
			t, err = time.Parse(format, after)
			if err == nil {
				break
			}
		}

		if err != nil {
			log.Fatalf("Couldn't parse %q as date %s\n", after, err)
		}

		afterDate = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 99, t.Nanosecond(), t.Location())
	}

	// Open the file
	csvfile, err := os.Open(flag.Args()[0])
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	defer csvfile.Close()

	r := csv.NewReader(csvfile)

	rows, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	seen := make(map[string]struct{})

	type ChargeRow struct {
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

	var record ChargeRow

	for i, row := range rows {
		if i == 0 {
			if row[0] != "Vehicle Name" && row[1] != "VIN" && row[2] != "Timezone" {
				log.Fatal("Invalid CSV File. Make sure you are using the Charging Data export, not Charge Summary")
			}

			continue
		}

		record = ChargeRow{
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

		if record.Supercharging != "true" {
			continue
		}

		if _, found := seen[record.Supercharger]; found {
			continue
		}

		seen[record.Supercharger] = struct{}{}

		chargeDate, err := time.Parse("2006-01-02 03:04PM", record.StartTime)
		if err != nil {
			log.Fatalf("Could not parse charge start time: %q\n", record.StartTime)
		}

		if chargeDate.After(afterDate) {
			fmt.Printf("%s;%s;%s\n", chargeDate.Format("01/02/2006"), chargeDate.Format("03:04PM"), record.Supercharger)
		}
	}
}
