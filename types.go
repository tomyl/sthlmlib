package main

// LoginRequest is the request body for the login mutation.
type LoginRequest struct {
	Query     string    `json:"query"`
	Variables Variables `json:"variables"`
}

type Variables struct {
	Operation  string     `json:"__operation"`
	LoginInput LoginInput `json:"loginInput"`
}

type LoginInput struct {
	CardNumber string `json:"cardNumber"`
	PinCode    string `json:"pinCode"`
}

// LoginResponse is the response from the login mutation.
type LoginResponse struct {
	Data struct {
		LoginPatron struct {
			Status        string      `json:"status"`
			StatusMessage interface{} `json:"statusMessage"`
		} `json:"loginPatron"`
	} `json:"data"`
}

// ProfileRequest is the request body for fetching the user profile.
type ProfileRequest struct {
	Query     string           `json:"query"`
	Variables ProfileVariables `json:"variables"`
}

type ProfileVariables struct {
	Operation string `json:"__operation"`
}

// ProfileResponse is the response from the getProfile query.
type ProfileResponse struct {
	Data struct {
		Patron Patron `json:"patron"`
	} `json:"data"`
}

// Patron holds all the user's information.
type Patron struct {
	PatronName   string        `json:"patronName"`
	PatronID     string        `json:"patronId"`
	Note         interface{}   `json:"note"`
	CardNumbers  []string      `json:"cardNumbers"`
	Debts        []Debt        `json:"debts"`
	Loans        Loans         `json:"loans"`
	Reservations []Reservation `json:"reservations"`
	Emails       []Email       `json:"emails"`
	PhoneNumbers []PhoneNumber `json:"phoneNumbers"`
}

// Debt represents a fee or charge.
type Debt struct {
	ID               string      `json:"id"`
	Type             string      `json:"type"`
	RefNo            string      `json:"refNo"`
	Titles           []DebtTitle `json:"titles"`
	Amount           float64     `json:"amount"`
	AmountInSubunits int         `json:"amountInSubunits"`
	Date             string      `json:"date"`
	DueDate          string      `json:"dueDate"`
	TypeDisplay      string      `json:"typeDisplay"`
}

type DebtTitle struct {
	Display    string `json:"display"`
	DueDate    string `json:"dueDate"`
	SomeNumber string `json:"someNumber"`
}

// Loans contains different types of loans.
type Loans struct {
	PhysicalLoans []PhysicalLoan `json:"physicalLoans"`
}

// PhysicalLoan represents a loan of a physical item.
type PhysicalLoan struct {
	Loan
	CatalogueRecordID string `json:"catalogueRecordId"`
}

// Loan represents a single loan.
type Loan struct {
	LoanID            string `json:"loanId"`
	LoanDate          string `json:"loanDate"`
	LoanDueDate       string `json:"loanDueDate"`
	LoanDateDiff      string `json:"loanDateDiff"`
	RemainingRenewals int    `json:"remainingRenewals"`
	LoanPerProduct    bool   `json:"loanPerProduct"`
	IsRenewable       bool   `json:"isRenewable"`
	NonRenewableMessage string `json:"nonRenewableMessage"`
	Status            string `json:"status"`
	Branch            struct {
		Name string `json:"name"`
	} `json:"branch"`
	Media struct {
		Key                 string `json:"key"`
		Title               string `json:"title"`
		Author              string `json:"author"`
		Image               string `json:"image"`
		MediaTypeDisplay    string `json:"mediaTypeDisplay"`
		MediaSubTypeDisplay string `json:"mediaSubTypeDisplay"`
	} `json:"media"`
}

// Reservation represents a single reservation.
type Reservation struct {
	CatalogueRecordID string      `json:"catalogueRecordId"`
	Fee               int         `json:"fee"`
	ReservationStatus string      `json:"reservationStatus"`
	Created           string      `json:"created"`
	ValidFrom         string      `json:"validFrom"`
	ValidTo           string      `json:"validTo"`
	QueueNumber       int         `json:"queueNumber"`
	NofCopies         int         `json:"nofCopies"`
	Status            interface{} `json:"status"`
	PickUpNumber      interface{} `json:"pickUpNumber"`
	PickupExpire      interface{} `json:"pickupExpire"`
	Media             struct {
		Title  string `json:"title"`
		Author string `json:"author"`
		Image  string `json:"image"`
	} `json:"media"`
	Branch struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"branch"`
}

// Email represents a user's email address.
type Email struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// PhoneNumber represents a user's phone number.
type PhoneNumber struct {
	ID     string `json:"id"`
	Number string `json:"number"`
	SMS    bool   `json:"sms"`
}
