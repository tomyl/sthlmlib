package main

import (
	"fmt"
	"time"

	"github.com/arran4/golang-ical"
)

// GenerateICal creates an iCalendar file from the user's loans.
func GenerateICal(patron *Patron) (string, error) {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodPublish)
	cal.SetCalscale("GREGORIAN")
	cal.SetName("Library Loans")
	cal.SetDescription("Due dates for library loans")

	for _, loan := range patron.Loans.PhysicalLoans {
		event := cal.AddEvent(loan.LoanID)
		event.SetCreatedTime(time.Now())
		event.SetDtStampTime(time.Now())
		event.SetModifiedAt(time.Now())
		event.SetSummary("Return: " + loan.Media.Title)
		event.SetDescription(fmt.Sprintf("Author: %s\nLibrary: %s", loan.Media.Author, loan.Branch.Name))

		dueDate, err := time.Parse("2006-01-02", loan.LoanDueDate)
		if err != nil {
			// Log the error and skip this event
			fmt.Printf("Skipping loan due to invalid date format: %v\n", err)
			continue
		}
		event.SetProperty(ics.ComponentPropertyDtStart, dueDate.Format("20060102"), ics.WithValue("DATE"))
	}

	return cal.Serialize(), nil
}

