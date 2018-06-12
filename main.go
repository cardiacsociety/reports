package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/34South/envr"
	"github.com/pkg/errors"
)

var MySQL = MySQLConnection{}

type SubscriptionReport struct {
	Rows []SubscriptionReportRow
}
type SubscriptionReportRow struct {
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
	fmt.Println("Report: Subscriptions Count by Type")
	r, err := reportSubsCount(QUERY_SUBSCRIPTION_COUNTS)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(r)

	fmt.Println("============================================================================================")
	fmt.Println("Report: Active Subscriptions Count by Type")
	r, err = reportSubsCount(QUERY_ACTIVE_SUBSCRIPTION_COUNTS)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(r)

	// Get all members for next set of reports todo need this for inactive members
	xm, err := allMembers()
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Println(xm)

	fmt.Println("============================================================================================")
	fmt.Println("Report: Membership Title by Year")
	err = reportTitleByYear(xm)
	if err != nil {
		log.Fatalln(err)
	}

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

// reports on the changes of status in order to determine how many members allowed their memberships to lapse,
// and what type of members these were
func reportTitleByYear(members []Member) error {

	// now we have a collection of all member histories, aggregate into some sensible structure
	startYear := oldestTitleYear(members)
	currentYear := time.Now().Year()

	fmt.Println("Titles ranging from", startYear, "-", currentYear)

	type titleCount struct {
		Title string
		Count int
	}

	type reportYear struct {
		Year int
		Data []titleCount

	}


	var report []reportYear

	// Initialise report
	for y := startYear; y <= currentYear; y++ {
		ry := reportYear{Year: y}
		report = append(report, ry)
	}

	for _, m := range members {
		for _, t := range m.TitleHistory {


			report[t.Year].Data = "test"
		}
	}
	//
	//a := map[string]int{"one": 1}
	//b := map[string]int{"two": 2}

	fmt.Println(report)

	return nil
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