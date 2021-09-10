package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	flag "github.com/spf13/pflag"
)

func usage() {
	exec, err := os.Executable()
	if err != nil {
		fmt.Println("Error finding own executable")
		os.Exit(1)
	}

	fmt.Printf("\nUsage: %s [--after DATE] CSVFILE\n", filepath.Base(exec))
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	var (
		after     string
		afterDate time.Time
		err       error
	)

	flag.Usage = usage
	help := flag.BoolP("help", "h", false, "Display help")

	flag.StringVar(&after, "after", "", "Date after which to display new supercharger visits")
	flag.Parse()

	if *help {
		flag.Usage()
	}

	if len(flag.Args()) < 1 {
		fmt.Println("Must provide csv file as positional argument")
		flag.Usage()
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
			fmt.Printf("Couldn't parse %q as date.\n", after)
			fmt.Println("Use one of: '2006-01-02', '01/02/2006', '1/2/2006' formats.")
			flag.Usage()
		}

		afterDate = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 99, t.Nanosecond(), t.Location())
	}

	// Open the file
	csvfile, err := os.Open(flag.Args()[0])
	if err != nil {
		fmt.Println("Couldn't open the csv file", err)
		os.Exit(1)
	}
	defer csvfile.Close()

	r := csv.NewReader(csvfile)

	rows, err := r.ReadAll()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
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

	fmt.Printf("%s;%s;%s\n", "Time", "First Visit", "Supercharger")
	for i, row := range rows {
		if i == 0 {
			if row[0] != "Vehicle Name" && row[1] != "VIN" && row[2] != "Timezone" {
				fmt.Print(
					"Invalid CSV File. Make sure you are using the Charging Data export (from Activity > Export > Charges),",
					"not Charge Summary (Locations)",
				)
				os.Exit(1)
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
			fmt.Printf("Could not parse charge start time: %q\n", record.StartTime)
			os.Exit(1)
		}

		scLocation := record.Supercharger
		if !strings.Contains(record.Supercharger, record.Location[len(record.Location)-4:]) {
			scLocation += record.Location[len(record.Location)-4:]
		}

		if chargeDate.After(afterDate) {
			fmt.Printf("%s;%s;%s\n", chargeDate.Format("03:04PM"), chargeDate.Format("01/02/2006"), scLocation)
		}
	}
}
