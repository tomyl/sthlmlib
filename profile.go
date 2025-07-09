package main

import "fmt"

const getProfileQuery = `
  fragment debt on Debt {
    id
    type
    refNo
    titles {
      display
      dueDate
      someNumber
    }
    amount
    amountInSubunits
    date
    dueDate
    typeDisplay
  }

  fragment loan on Loan {
    loanId
    loanDate
    loanDueDate
    loanDateDiff
    remainingRenewals
    loanPerProduct
    isRenewable
    status
    branch {
      name
    }
    media {
      key
      title
      author
      image
      mediaTypeDisplay
      mediaSubTypeDisplay
    }
  }

  fragment reservation on Reservation {
    catalogueRecordId
    fee
    reservationStatus
    created
    validFrom
    validTo
    queueNumber
    nofCopies
    status
    pickUpNumber
    pickupExpire
    media {
      title
      author
      image
    }
    branch {
      name
      slug
    }
  }

  query getProfile {
    patron {
      ... on Patron {
        patronName
        patronId
        note
        cardNumbers
        debts {
          ...debt
        }
        loans {
          physicalLoans {
            catalogueRecordId
            ...loan
          }
        }
        reservations {
          ...reservation
        }
        emails {
          id
          email
        }
        phoneNumbers {
          id
          number
          sms
        }
      }
      ... on Response {
        status
        statusMessage
      }
    }
  }
`

// GetProfile fetches the user's complete profile.
func (c *Client) GetProfile() (*Patron, error) {
	requestBody := ProfileRequest{
		Query: getProfileQuery,
		Variables: ProfileVariables{
			Operation: "getProfile",
		},
	}

	var responseBody ProfileResponse
	if err := c.query(&requestBody, &responseBody); err != nil {
		return nil, fmt.Errorf("could not get profile: %w", err)
	}

	return &responseBody.Data.Patron, nil
}
