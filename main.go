package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/StephenBrown2/tezlab-export-sc/tezlab"

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

func parseDate(s string) (t time.Time) {
	var err error

	validFormats := []string{"2006-01-02", "01/02/2006", "1/2/2006"}
	for _, format := range validFormats {
		t, err = time.Parse(format, s)
		if err == nil {
			break
		}
	}

	if err != nil {
		fmt.Printf("Couldn't parse %q as date.\n", s)
		fmt.Printf("Use one of: '%s' formats.\n", strings.Join(validFormats, "', '"))
		flag.Usage()
	}

	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 99, t.Nanosecond(), t.Location())
}

func main() {
	var (
		after, tzAfter         string
		afterDate, tzAfterDate time.Time
		debug, trace           bool
		err                    error
	)

	tzAfterDate = time.Now().AddDate(0, 0, -7)

	outputHeader := []string{"Index", "Time", "Date", "Supercharger"}

	flag.Usage = usage
	help := flag.BoolP("help", "h", false, "Display help")

	flag.StringVar(&after, "after", "", "Date after which to display new supercharger visits")
	flag.StringVar(&tzAfter, "tz-after", "", fmt.Sprintf("Date after which to fetch timezone for supercharger visits (Default: one week ago - %s)", tzAfterDate.Format("1/2/2006")))
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
		afterDate = parseDate(after)
		tzAfterDate = afterDate
	}

	if tzAfter != "" {
		tzAfterDate = parseDate(tzAfter)
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

	var record tezlab.ChargeRecord
	w := csv.NewWriter(os.Stdout)

	for i, row := range rows {
		if i == 0 {
			if len(row) != 26 || row[13] != "Supercharger" {
				log.Error(
					"Invalid CSV File! Make sure you are using the Charging Data export",
					"(from Activity > Export > Charges), not Charge Summary (Locations)",
				)
				log.Errorf("Got CSV header: '%s'\n", strings.Join(row, "', '"))
				os.Exit(1)
			}

			w.Write(outputHeader)

			continue
		}

		log.Tracef("Loading record from row: '%s'", strings.Join(row, "', '"))
		record = tezlab.NewRecordFromCSV(row)
		log.Tracef("Loaded record: %+v", record)

		if !record.IsSupercharger() {
			log.Debug("Skipping because not a supercharger")
			continue
		}

		scLocation := record.SuperchargerLocation()

		if _, found := seen[scLocation]; found {
			log.Debugf("Skipping because %s already seen", scLocation)
			continue
		}

		seen[scLocation] = struct{}{}

		chargeDate := record.StartTimeInRecordTimezone()

		if chargeDate.After(tzAfterDate) {
			log.Debugf("Charge date before: %s", chargeDate)
			chargeDate = chargeDate.In(record.LocationTimezone())
			log.Debugf("Charge date after: %s", chargeDate)
		}

		if chargeDate.After(afterDate) {
			log.Trace(outputHeader)
			w.Write([]string{fmt.Sprint(len(seen)), chargeDate.Format("03:04PM MST"), chargeDate.Format("01/02/2006"), scLocation})
			w.Flush()
		}
	}
}
