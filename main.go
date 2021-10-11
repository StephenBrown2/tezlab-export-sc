package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/goinggo/timezone"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

func usage() {
	exec, err := os.Executable()
	if err != nil {
		log.Error("Error finding own executable")
		os.Exit(1)
	}

	fmt.Printf("\nUsage: %s [options] CSVFILE\n", filepath.Base(exec))
	flag.PrintDefaults()
	os.Exit(2)
}

func username() string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		panic("Can't read BuildInfo")
	}
	s := strings.Split(buildInfo.Path, "/")
	return s[len(s)-2]
}

func main() {
	var (
		after, tzAfter         string
		afterDate, tzAfterDate time.Time
		debug, trace           bool
		err                    error
	)

	tzAfterDate = time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())

	outputHeader := strings.Join([]string{"Time", "First Visit", "Supercharger"}, ";")

	flag.Usage = usage
	help := flag.BoolP("help", "h", false, "Display help")

	flag.StringVar(&after, "after", "", "Date after which to display new supercharger visits")
	flag.StringVar(&tzAfter, "tz-after", "", fmt.Sprintf("Date after which to fetch timezone for supercharger visits (Default: %s)", tzAfterDate.Format("1/2/2006")))
	flag.BoolVar(&debug, "debug", false, "Display debugging information for troubleshooting")
	flag.BoolVar(&trace, "trace", false, "Display more information for troubleshooting (overrides --debug)")

	flag.Parse()

	if *help {
		flag.Usage()
	}

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	if trace {
		log.SetLevel(log.TraceLevel)
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
		tzAfterDate = afterDate
	}

	if tzAfter != "" {
		var t time.Time
		for _, format := range []string{"2006-01-02", "01/02/2006", "1/2/2006"} {
			t, err = time.Parse(format, tzAfter)
			if err == nil {
				break
			}
		}

		if err != nil {
			fmt.Printf("Couldn't parse %q as date.\n", tzAfter)
			fmt.Println("Use one of: '2006-01-02', '01/02/2006', '1/2/2006' formats.")
			flag.Usage()
		}

		tzAfterDate = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 99, t.Nanosecond(), t.Location())
	}

	// Open the file
	csvfile, err := os.Open(flag.Args()[0])
	if err != nil {
		log.Error("Couldn't open the csv file", err)
		os.Exit(1)
	}
	defer csvfile.Close()

	r := csv.NewReader(csvfile)

	rows, err := r.ReadAll()
	if err != nil {
		log.Error(err)
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

	for i, row := range rows {
		if i == 0 {
			if row[13] != "Supercharger" {
				log.Error(
					"Invalid CSV File! Make sure you are using the Charging Data export",
					"(from Activity > Export > Charges), not Charge Summary (Locations)",
				)
				log.Errorf("Got CSV header: '%s'\n", strings.Join(row, "', '"))
				os.Exit(1)
			}

			fmt.Println(outputHeader)

			continue
		}

		log.Tracef("Loading record from row: '%s'", strings.Join(row, "', '"))
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
		log.Tracef("Loaded record: %+v", record)

		if record.Supercharging != "true" && !strings.Contains(record.UserChargeLocation, "Super") {
			log.Debug("Skipping because not a supercharger")
			continue
		}

		scLocation := record.Supercharger

		if scLocation == "" { // User-entered Supercharger
			scLocation = record.Location
		}

		log.Debugf("scLocation: %q", scLocation)
		if !strings.Contains(scLocation, record.Location[len(record.Location)-4:]) && scLocation[len(scLocation)-4] != ',' {
			scLocation += record.Location[len(record.Location)-4:]
			log.Debugf("scLocation new: %q", scLocation)
		}

		if strings.HasSuffix(scLocation, "COs") {
			scLocation = strings.TrimRight(scLocation, "s")
			log.Debugf("scLocation trimmed: %q", scLocation)
		}

		if _, found := seen[scLocation]; found {
			log.Debugf("Skipping because %s already seen", scLocation)
			continue
		}

		seen[scLocation] = struct{}{}

		log.Debugf("Loading record Timezone: %q", record.Timezone)
		recordLoc, err := time.LoadLocation(record.Timezone)
		if err != nil {
			log.Errorf("Could not parse charge location timezone: %q - %s\n", record.Timezone, err)
			os.Exit(1)
		}

		log.Debugf("Parsing time %q in Location %q", record.StartTime, recordLoc)
		chargeDate, err := time.ParseInLocation("2006-01-02 03:04PM", record.StartTime, recordLoc)
		if err != nil {
			log.Errorf("Could not parse charge start time: %q - %s\n", record.StartTime, err)
			os.Exit(1)
		}

		if chargeDate.After(tzAfterDate) {
			coord := strings.Split(record.Coordinate, ",")
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

			log.Debugf("Charge date before: %s", chargeDate)
			chargeDate = chargeDate.In(location)
			log.Debugf("Charge date after: %s", chargeDate)
		}

		if chargeDate.After(afterDate) {
			log.Trace(outputHeader)
			fmt.Printf("%s;%s;%s\n", chargeDate.Format("03:04PM MST"), chargeDate.Format("01/02/2006"), scLocation)
		}
	}
}
