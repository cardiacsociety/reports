package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/34South/envr"
	"github.com/pkg/errors"
	"encoding/json"
)

var MySQL = MySQLConnection{}

type SubscriptionReport struct {
	Rows []SubscriptionReportRow
}
type SubscriptionReportRow struct {
	Subscription string
	Count        int
}

type LapsedReport struct {
	Rows []LapsedReportRow
}
type LapsedReportRow struct {
	LapsedYear int
	Subscription string
	Count        int
}

func init() {
	envr.New("reportEnv", []string{
		"MYSQL_DSN",
	}).Auto()

	MySQL.DSN = os.Getenv("MYSQL_DSN")
}

func main() {

	fmt.Println("CSANZ Report Generator")

	fmt.Print("Connecting to database...")
	err := MySQL.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("ok")

	fmt.Println("============================================================================================")
	fmt.Println("Report: Subscription Count by Type (All)")
	r, err := reportSubsCount(QUERY_SUBSCRIPTION_COUNTS)
	if err != nil {
		log.Fatalln(err)
	}
	printJSON(r)

	fmt.Println("============================================================================================")
	fmt.Println("Report: Subscription Count by Type (Active Members)")
	r, err = reportSubsCount(QUERY_ACTIVE_SUBSCRIPTION_COUNTS)
	if err != nil {
		log.Fatalln(err)
	}
	printJSON(r)

	// Get all members for next set of reports todo need this for inactive members
	xm, err := allMembers()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("============================================================================================")
	fmt.Println("Report: Membership Title by Year Granted")
	r2, err := reportTitleByYear(xm)
	if err != nil {
		log.Fatalln(err)
	}
	printJSON(r2)

	fmt.Println("============================================================================================")
	fmt.Println("Report: Lapsed Memberships By Title and Year Lapsed")
	r3, err := reportLapsedByYear(QUERY_CURRENTLY_LAPSED_MEMBERS_COUNT_TITLE_YEAR)
	if err != nil {
		log.Fatalln(err)
	}
	printJSON(r3)
}

func reportSubsCount(query string) (SubscriptionReport, error) {

	var report SubscriptionReport

	rows, err := MySQL.Session.Query(query)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
		r := SubscriptionReportRow{}
		err := rows.Scan(&r.Subscription, &r.Count)
		if err != nil {
			return report, errors.Wrap(err, "reportSubsCount")
		}
		report.Rows = append(report.Rows, r)
	}

	return report, nil
}

func reportLapsedByYear(query string) (LapsedReport, error) {

	var report LapsedReport

	rows, err := MySQL.Session.Query(query)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
		r := LapsedReportRow{}
		err := rows.Scan(&r.LapsedYear, &r.Subscription, &r.Count)
		if err != nil {
			return report, errors.Wrap(err, "reportLapsedByYear")
		}
		report.Rows = append(report.Rows, r)
	}

	return report, nil
}

// reports on the changes of status in order to determine how many members allowed their memberships to lapse,
// and what type of members these were
func reportTitleByYear(members []Member) (map[int][]map[string]int, error) {

	// now we have a collection of all member histories, aggregate into some sensible structure
	startYear := oldestTitleYear(members)
	currentYear := time.Now().Year()

	fmt.Println("Titles ranging from", startYear, "-", currentYear)

	report := map[int][]map[string]int{}

	// Initialise report
	fmt.Println("Year\tApplicant\tAffiliate\tAssociate\tFellow\tFellow Life\tOrdinary\tLife")

	for y := startYear; y <= currentYear; y++ {

		report[y] = []map[string]int{}

		// todo ... get this from database!!!
		types := []string{"Applicant", "Affiliate", "Associate", "Fellow", "Fellow & Life", "Ordinary", "Life"}

		for _, t := range types {
			c := titleYearCount(members, t, y)
			a := map[string]int{t: c}
			report[y] = append(report[y], a)
		}
	}

	return report, nil
}

// allMembers creates a slice of Member values for all member records
func allMembers() ([]Member, error) {

	var xm []Member

	rows, err := MySQL.Session.Query(QUERY_MEMBER_ID)
	if err != nil {
		return nil, errors.Wrap(err, "allMembers")
	}
	defer rows.Close()

	for rows.Next() {
		m := Member{}
		rows.Scan(&m.ID)

		err := m.setStatusHistory()
		if err != nil {
			return xm, errors.Wrap(err, "allMembers")
		}

		err = m.setTitleHistory()
		if err != nil {
			return xm, errors.Wrap(err, "allMembers")
		}

		xm = append(xm, m)
	}

	return xm, nil
}


// oldestTitleYear scours the member values to locate the oldest title year
func oldestTitleYear(members []Member) int {
	y := time.Now().Year()
	for _, m := range members {
		for _, t := range m.TitleHistory {
			if t.Year < y {
				y = t.Year
			}
		}
	}
	return y
}

// Returns the number of occurrences of a particular title, in the specified year
func titleYearCount(members []Member, title string, year int) int {

	var c int

	for _, m := range members {
		for _, t := range m.TitleHistory {
			if t.Year == year && t.Name == title {
				c++
			}
		}
	}

	return c
}

func printJSON(data interface{}) {

	xb, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(xb))
}