package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/request"
	"github.com/VamaSingapore/vama-api/internal/entities/_shared/models/response"
	"github.com/VamaSingapore/vama-api/internal/utils"
	vlog "github.com/VamaSingapore/vama-api/internal/vamaLogger"
	"github.com/VamaSingapore/vama-api/pkg/httperr"
	"github.com/nyaruka/phonenumbers"
)

const uploadContactsErr = "Something went wrong when uploading contacts."

func (svc *usecase) UploadContacts(ctx context.Context, req request.UploadContacts) error {
	user, userErr := svc.user.GetCurrentUser(ctx)
	if userErr != nil {
		return httperr.NewCtx(
			ctx,
			500,
			http.StatusInternalServerError,
			uploadContactsErr,
			fmt.Sprintf("Could not find user in the current context. Err: %v", userErr),
		)
	}
	if user == nil {
		return httperr.NewCtx(
			ctx,
			404,
			http.StatusNotFound,
			uploadContactsErr,
			"user came back nil in UploadContacts",
		)
	}

	var wg sync.WaitGroup
	wg.Add(len(req.Contacts))
	for _, contact := range req.Contacts {
		go func(contact request.Contact) {
			defer wg.Done()

			if len(contact.FirstName) > constants.MAX_UPLOAD_CONTACT_NAME_LENGTH {
				contact.FirstName = contact.FirstName[:constants.MAX_UPLOAD_CONTACT_NAME_LENGTH]
			}

			if len(contact.LastName) > constants.MAX_UPLOAD_CONTACT_NAME_LENGTH {
				contact.LastName = contact.LastName[:constants.MAX_UPLOAD_CONTACT_NAME_LENGTH]
			}

			svc.uploadContact(ctx, svc.repo.MasterNode(), user, contact)

		}(contact)
	}
	wg.Wait()

	return nil
}

func (svc *usecase) uploadContact(ctx context.Context, runnable utils.Runnable, user *response.User, contact request.Contact) error {
	number, countryCode, parseErr := parsePhoneNumberWithRetry(ctx, contact.Phone, user.CountryCode)
	if parseErr != nil {
		return parseErr
	}

	otherUser, getUserByPhoneErr := svc.user.GetUserByPhone(ctx, runnable, number)
	if getUserByPhoneErr != nil {
		vlog.Errorf(ctx, "Error getting user by number %s: %v", number, getUserByPhoneErr)
		return getUserByPhoneErr
	}
	if otherUser != nil {
		_, createContactErr := svc.CreateContact(ctx, otherUser.ID)
		return createContactErr
	}

	phoneStruct := request.Phone{
		Number:      number,
		CountryCode: countryCode,
	}

	addPendingContactErr := svc.repo.InsertPendingContact(ctx, user.ID, phoneStruct, contact.FirstName, contact.LastName)
	if addPendingContactErr != nil {
		vlog.Errorf(ctx, "Error adding %d's pending contact %s: %v", user.ID, number, addPendingContactErr)
		return addPendingContactErr
	}

	return nil
}

func parsePhoneNumberWithRetry(ctx context.Context, number string, retryCountryCode string) (string, string, error) {
	phoneStruct, parseErr := phonenumbers.Parse(number, "")
	if parseErr != nil {
		retryPhoneStruct, retryParseErr := phonenumbers.Parse(number, retryCountryCode)
		if retryParseErr != nil {
			vlog.Errorf(ctx, "Could not parse phone number %s with retry country code %s: %v", number, retryCountryCode, retryParseErr)
			return "", "", retryParseErr
		}
		phoneStruct = retryPhoneStruct
	}
	formattedNumber := phonenumbers.Format(phoneStruct, phonenumbers.E164)
	countryCodeString := phonenumbers.GetRegionCodeForCountryCode(int(*phoneStruct.CountryCode))
	return formattedNumber, countryCodeString, nil

}
