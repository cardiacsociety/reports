package main

import (
	"time"
	"github.com/pkg/errors"
)

type Member struct {
	ID            int
	StatusHistory []Status
	TitleHistory  []Title
}

type Status struct {
	Date time.Time
	Year int
	ID   int
	Name string
}

type Title struct {
	Date time.Time
	Year int
	ID   int
	Name string
}

func (m *Member) setStatusHistory() error {

	rows, err := MySQL.Session.Query(QUERY_MEMBER_STATUS_HISTORY, m.ID)
	if err != nil {
		return errors.Wrap(err, "memberStatusHistory")
	}
	defer rows.Close()

	for rows.Next() {

		var s Status
		var date string

		err := rows.Scan(&date, &s.ID, &s.Name)
		if err != nil {
			return errors.Wrap(err, "memberStatusHistory")
		}

		if date != "0000-00-00" {

			s.Date, err = time.Parse("2006-01-02", date)
			if err != nil {
				return errors.Wrap(err, "memberStatusHistory")
			}
			s.Year = s.Date.Year()

			m.StatusHistory = append(m.StatusHistory, s)
		}
	}

	return nil
}

func (m *Member) setTitleHistory() error {

	rows, err := MySQL.Session.Query(QUERY_MEMBER_TITLE_HISTORY, m.ID)
	if err != nil {
		return errors.Wrap(err, "memberTitleHistory")
	}
	defer rows.Close()

	for rows.Next() {

		var t Title
		var date string

		err := rows.Scan(&date, &t.ID, &t.Name)
		if err != nil {
			return errors.Wrap(err, "memberTitleHistory")
		}

		if date != "0000-00-00" {

			t.Date, err = time.Parse("2006-01-02", date)
			if err != nil {
				return errors.Wrap(err, "memberTitleHistory")
			}
			t.Year = t.Date.Year()

			m.TitleHistory = append(m.TitleHistory, t)
		}
	}

	return nil
}

