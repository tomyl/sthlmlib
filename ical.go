package main

import (
	"fmt"
	"time"

	"github.com/arran4/golang-ical"
)

// GenerateICal creates an iCalendar file from the user's loans.
func GenerateICal(patron *Patron, group bool) (string, error) {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodPublish)
	cal.SetCalscale("GREGORIAN")
	cal.SetName("Library Loans")
	cal.SetDescription("Due dates for library loans and reservations")

	if group {
		loansByDate := make(map[string][]PhysicalLoan)
		for _, loan := range patron.Loans.PhysicalLoans {
			loansByDate[loan.LoanDueDate] = append(loansByDate[loan.LoanDueDate], loan)
		}

		for dueDateStr, loans := range loansByDate {
			dueDate, err := time.Parse("2006-01-02", dueDateStr)
			if err != nil {
				fmt.Printf("Skipping loans due to invalid date format: %v\n", err)
				continue
			}

			summary := fmt.Sprintf("Return %d items", len(loans))
			var description string
			for _, loan := range loans {
				description += fmt.Sprintf("- %s by %s\n", loan.Media.Title, loan.Media.Author)
			}

			event := cal.AddEvent(dueDate.Format("2006-01-02"))
			event.SetCreatedTime(time.Now())
			event.SetDtStampTime(time.Now())
			event.SetModifiedAt(time.Now())
			event.SetSummary(summary)
			event.SetDescription(description)
			event.SetProperty(ics.ComponentPropertyDtStart, dueDate.Format("20060102"), ics.WithValue("DATE"))
		}
	} else {
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
	}

	for _, res := range patron.Reservations {
		if res.PickupExpire != nil {
			event := cal.AddEvent(res.ID)
			event.SetCreatedTime(time.Now())
			event.SetDtStampTime(time.Now())
			event.SetModifiedAt(time.Now())
			summary := "Pickup: " + res.Media.Title
			if res.PickUpNumber != nil {
				summary += fmt.Sprintf(" (#%s)", *res.PickUpNumber)
			}
			event.SetSummary(summary)
			event.SetDescription(fmt.Sprintf("Author: %s\nLibrary: %s", res.Media.Author, res.Branch.Name))

			expireDate, err := time.Parse("2006-01-02", *res.PickupExpire)
			if err != nil {
				// Log the error and skip this event
				fmt.Printf("Skipping reservation due to invalid date format: %v\n", err)
				continue
			}
			event.SetProperty(ics.ComponentPropertyDtStart, expireDate.Format("20060102"), ics.WithValue("DATE"))
		}
	}

	return cal.Serialize(), nil
}