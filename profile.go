package main

import "fmt"

const getProfileQuery = `
query {
  patron {
    ... on Patron {
      patronName
      cardNumbers
      loans {
        physicalLoans {
          loanId
          catalogueRecordId
          loanDate
          loanDueDate
          loanDateDiff
          remainingRenewals
          loanPerProduct
          isRenewable
          nonRenewableMessage
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
      }
      reservations {
        id
        branch {
          name
          slug
        }
        catalogueRecordId
        created
        editable
        pickupExpire
        pickUpNumber
        nofCopies
        queueNumber
        reservationStatus
        status
        validFrom
        validTo
        media {
          key
          title
          author
          image
          mediaTypeDisplay
          mediaSubTypeDisplay
          language
          targetGroup
          publishedDate
        }
        note
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
