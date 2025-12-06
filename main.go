package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
)

func main() {
	cardNumber := flag.String("card-number", "", "Library card number")
	pinCode := flag.String("pin", "", "PIN code")
	dump := flag.Bool("dump", false, "Dump the raw JSON response")
	ical := flag.Bool("ical", false, "Output loans as an iCal file to stdout")
	group := flag.Bool("group", false, "Group ical entries for the same day")
	flag.Parse()

	if *cardNumber == "" || *pinCode == "" {
		log.Fatal("card-number and pin are required")
	}

	client, err := NewClient(*cardNumber, *pinCode)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	if *dump {
		dumpProfile(client)
		return
	}

	profile, err := client.GetProfile()
	if err != nil {
		log.Fatalf("Failed to get profile: %v", err)
	}

	if *ical {
		icalData, err := GenerateICal(profile, *group)
		if err != nil {
			log.Fatalf("Failed to generate iCal data: %v", err)
		}
		fmt.Print(icalData)
		return
	}

	printProfile(profile)
}

func dumpProfile(client *Client) {
	requestBody := ProfileRequest{
		Query: getProfileQuery,
		Variables: ProfileVariables{
			Operation: "getProfile",
		},
	}
	var responseBody json.RawMessage
	if err := client.query(&requestBody, &responseBody); err != nil {
		log.Fatalf("Failed to get raw profile: %v", err)
	}
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(responseBody), "", "  "); err != nil {
		log.Fatalf("Failed to format JSON: %v", err)
	}
	fmt.Println(prettyJSON.String())
}

func printProfile(profile *Patron) {
	fmt.Printf("Welcome, %s!\n", profile.PatronName)

	fmt.Println("\nLoans:")
	if len(profile.Loans.PhysicalLoans) == 0 {
		fmt.Println("  No loans.")
	} else {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Title\tAuthor\tDue Date\tLibrary")
		fmt.Fprintln(w, "-----\t------\t--------\t-------")
		for _, loan := range profile.Loans.PhysicalLoans {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", loan.Media.Title, loan.Media.Author, loan.LoanDueDate, loan.Branch.Name)
		}
		w.Flush()
	}

	fmt.Println("\nReservations:")
	if len(profile.Reservations) == 0 {
		fmt.Println("  No reservations.")
	} else {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Title\tAuthor\tStatus\tExpires\tPickup Number\tLibrary")
		fmt.Fprintln(w, "-----\t------\t------\t-------\t-------------\t-------")
		for _, res := range profile.Reservations {
			pickupExpire := ""
			if res.PickupExpire != nil {
				pickupExpire = *res.PickupExpire
			}
			pickupNumber := ""
			if res.PickUpNumber != nil {
				pickupNumber = *res.PickUpNumber
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", res.Media.Title, res.Media.Author, res.ReservationStatus, pickupExpire, pickupNumber, res.Branch.Name)
		}
		w.Flush()
	}
}

func printReservations(profile *Patron) {
	fmt.Println("\nReservations:")
	if len(profile.Reservations) == 0 {
		fmt.Println("  No reservations.")
	} else {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Title\tAuthor\tStatus\tExpires\tPickup Number")
		fmt.Fprintln(w, "-----\t------\t------\t-------\t-------------")
		for _, res := range profile.Reservations {
			pickupExpire := ""
			if res.PickupExpire != nil {
				pickupExpire = *res.PickupExpire
			}
			pickupNumber := ""
			if res.PickUpNumber != nil {
				pickupNumber = *res.PickUpNumber
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", res.Media.Title, res.Media.Author, res.ReservationStatus, pickupExpire, pickupNumber)
		}
		w.Flush()
	}
}
