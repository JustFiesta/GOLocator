package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	// General flags
	updateFlag := flag.NewFlagSet("update", flag.ExitOnError)
	searchFlag := flag.NewFlagSet("search", flag.ExitOnError)
	travelFlag := flag.NewFlagSet("travel", flag.ExitOnError)

	// Update flags
	updateUsername := updateFlag.String("u", "", "Username for updating location")
	updateCoordinates := updateFlag.String("c", "", "New coordinates (latitude,longitude)")

	// Search flags
	searchTime := searchFlag.String("t", "", "Date/time for search in YYYY-MM-DDTHH:MM:SS+UTC format (optional)")

	// Travel flags
	travelCoordinates := travelFlag.String("c", "", "Coordinates for travel (latitude,longitude)")

	// Check if 2 arguments were passed
	if len(os.Args) < 2 {
		fmt.Println("Usage: goloc [update|search|travel] <options>")
		os.Exit(1)
	}

	// Set flag according to first argument
	switch os.Args[1] {
	case "update":
		updateFlag.Parse(os.Args[2:])
		if *updateUsername == "" || *updateCoordinates == "" {
			fmt.Println("Usage: goloc update -u <username> -c <coordinates>")
			updateFlag.PrintDefaults()
			os.Exit(1)
		}
		latitude, longitude, err := validateCoordinates(*updateCoordinates)
		if err != nil {
			fmt.Println("Invalid coordinates:", err)
			os.Exit(1)
		}
		// user service API comm
		fmt.Printf("Updating location for user %s to coordinates (%f,%f)\n", *updateUsername, latitude, longitude)
	case "search":
		searchFlag.Parse(os.Args[2:])
		if *searchTime == "" {
			fmt.Println("Usage: goloc search -t <time>")
			searchFlag.PrintDefaults()
			os.Exit(1)
		}
		// user service API comm
		fmt.Printf("Searching for user travel distance till %s\n", *searchTime)
	case "travel":
		travelFlag.Parse(os.Args[2:])
		if *travelCoordinates == "" {
			fmt.Println("Usage: goloc travel -c <coordinates>")
			travelFlag.PrintDefaults()
			os.Exit(1)
		}
		latitude, longitude, err := validateCoordinates(*travelCoordinates)
		if err != nil {
			fmt.Println("Invalid coordinates:", err)
			os.Exit(1)
		}
		// location service API comm
		fmt.Printf("Tracking user travel to coordinates (%f,%f)\n", latitude, longitude)
	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}
}

func validateCoordinates(coordStr string) (latitude, longitude float64, err error) {
	// Regexp pattern to validate latitude and longitude
	coordPattern := `^[-+]?([1-8]?\d(\.\d{1,8})?|90(\.0{1,8})?),\s*[-+]?(180(\.0{1,8})?|((1[0-7]\d)|([1-9]?\d))(\.\d{1,8})?)$`
	coordRegex := regexp.MustCompile(coordPattern)

	if !coordRegex.MatchString(coordStr) {
		return 0, 0, fmt.Errorf("Invalid coordinates format. Use decimal format (latitude,longitude)")
	}

	// Split coordinates string by comma
	parts := strings.Split(coordStr, ",")

	// Convert latitude and longitude to float64
	latitude, err = strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("Invalid latitude")
	}

	longitude, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("Invalid longitude")
	}

	// Validate latitude range (-90 to 90) and longitude range (-180 to 180)
	if latitude < -90 || latitude > 90 {
		return 0, 0, fmt.Errorf("Latitude out of range (-90 to 90)")
	}

	if longitude < -180 || longitude > 180 {
		return 0, 0, fmt.Errorf("Longitude out of range (-180 to 180)")
	}

	return latitude, longitude, nil
}
