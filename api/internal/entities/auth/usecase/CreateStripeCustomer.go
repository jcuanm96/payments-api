package service

import (
	"github.com/stripe/stripe-go/v72"
)

func (svc *usecase) CreateStripeCustomer(stripeParams *stripe.CustomerParams) (string, error) {
	custmr, createStripeCustomerErr := svc.stripeClient.Customers.New(stripeParams)
	if createStripeCustomerErr != nil {
		return "", createStripeCustomerErr
	}

	return custmr.ID, nil
}
