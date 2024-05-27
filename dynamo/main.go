package dynamo

import (
	"vue-converter-backend/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type CreateCustomerInterface interface {
	CreateCustomerFunc(customer models.Customer) error
}

type GetCustomerInterface interface {
	GetCustomerFunc(id string) (models.Customer, error)
}

type TopUpTokenBalanceForCustomerInterface interface {
	TopUpTokenBalanceForCustomerFunc(customer models.Customer, tokenAmount int) error
}

type DeductTokensForCustomerInterface interface {
	DeductTokensForCustomerFunc(customer models.Customer, tokenAmount int) error
}

type CreateCustomer struct{}

type GetCustomer struct{}

type TopUpTokenBalanceForCustomer struct{}

type DeductTokenBalanceForCustomers struct{}

func (f *CreateCustomer) CreateCustomerFunc(customer models.Customer) error {
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{Region: aws.String("eu-central-1")})
	table := db.Table("Customer")
	error := table.Put(customer).Run()
	return error
}

func (f *GetCustomer) GetCustomerFunc(id string) (models.Customer, error) {
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{Region: aws.String("eu-central-1")})
	table := db.Table("Customers")
	var customer models.Customer
	error := table.Get("ID", id).One(customer)

	return customer, error
}

func (f *TopUpTokenBalanceForCustomer) TopUpTokenBalanceForCustomerFunc(customer models.Customer, tokenAmount int) error {
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{Region: aws.String("eu-central-1")})
	table := db.Table("Customers")
	error := table.Update("ID", customer.Id).Add("AiCredits", tokenAmount).Run()
	return error
}

func (f *DeductTokenBalanceForCustomers) DeductTokenBalanceForCustomersFunc(customer models.Customer, tokenAmount int) error {
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{Region: aws.String("eu-central-1")})
	table := db.Table("Customers")
	error := table.Update("ID", customer.Id).Add("AiCredits", -tokenAmount).Run()
	return error
}
