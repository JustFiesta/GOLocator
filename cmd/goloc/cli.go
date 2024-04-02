// CLI tool for comunication with GOLocator microservices

package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	// General flags
	updateFlag := flag.NewFlagSet("update", flag.ExitOnError)
	searchFlag := flag.NewFlagSet("search", flag.ExitOnError)
	travelFlag := flag.NewFlagSet("travel", flag.ExitOnError)

	// Update flags
	updateUsername := updateFlag.String("u", "", "Username for updating location. Only letters and numbers are avalible  in this order: lettersnumbers")
	updateCoordinates := updateFlag.String("c", "", "New coordinates (latitude,longitude)")

	// Search flags
	searchLocation := searchFlag.String("c", "", "Coordinates used to search into")
	searchRadius := searchFlag.String("r", "", "Radius for searching in some area")

	// Travel flags
	travelTimeBound := travelFlag.String("t", "", "Date/time for search in YYYY-MM-DDTHH:MM:SS+UTC format")

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

		isValid := isValidUsername(*updateUsername)
		if !isValid {
			fmt.Println("Username is not valid string! example: user1212 <lettersnumbers>")
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
		if *searchLocation == "" || *searchRadius == "" {
			fmt.Println("Usage: goloc search -c <coordinates> -r <radius>")
			searchFlag.PrintDefaults()
			os.Exit(1)
		}
		latitude, longitude, err := validateCoordinates(*searchLocation)
		if err != nil {
			fmt.Println("Invalid coordinates:", err)
			os.Exit(1)
		}

		radius, err := strconv.ParseFloat(*searchRadius, 64)

        if err != nil {
            fmt.Println("Invalid radius:", err)
            os.Exit(1)
        }

		// user service API comm
		fmt.Printf("Searching for users in %f, %f within %f radius\n", latitude, longitude, radius)

	case "travel":
		travelFlag.Parse(os.Args[2:])
		if *travelTimeBound == "" {
			fmt.Println("Usage: goloc travel -t <YYYY-MM-DDTHH:MM:SS+UTC>")
			travelFlag.PrintDefaults()
			os.Exit(1)
		}
		err := validateDate(*travelTimeBound)
		if !err {
			fmt.Println("Invalid date format")
			os.Exit(1)
		}
		// location service API comm
		fmt.Printf("Checking user traveled distance since: %v)\n", travelTimeBound)

	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}
}

// validateCoordinates - Validate user inputed coordinates AND return both latitude and longitude from one string
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

// validateDate - validates date string, returns true/false
func validateDate(date string) bool {
	// ISO 8601 format
	iso8601Pattern := `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}[+-]\d{2}:\d{2}$`

	// matching format with given string date
	match, err := regexp.MatchString(iso8601Pattern, date)

	if err != nil {
		return false 
	}

	if !match {
		return false 
	}

	// check if date is valid
	_, err = time.Parse(time.RFC3339, date)
	if err != nil {
		return false
	}

	return true
}

// isValidUsername - checks if username is valid according to `^[a-zA-Z0-9]{4,16}$` pattern
func isValidUsername(username string) bool {
	// pattern for username validation
	usernamePattern := `^[a-zA-Z0-9]{4,16}$`

	// matching username with the pattern
	match, err := regexp.MatchString(usernamePattern, username)

	if err != nil {
		return false
	}

	return match
}
