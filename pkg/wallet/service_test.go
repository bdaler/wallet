package wallet

import (
	"testing"
)

func TestService_FindAccountByID_success(t *testing.T) {
	var service Service
	service.RegisterAccount("9127660305")

	account, err := service.FindAccountByID(1)

	if err != nil {
		t.Errorf("account => %v", account)
	}

}
func TestService_FindAccountByID_notFound(t *testing.T) {
	var service Service
	service.RegisterAccount("9127660305")

	account, err := service.FindAccountByID(2)

	if err == nil {
		t.Errorf("method returned nil error, account => %v", account)
	}

}

func TestService_Reject_success_user(t *testing.T) {
	var service Service
	service.RegisterAccount("9127660305")
	account, err := service.FindAccountByID(1)

	if err != nil {
		t.Errorf("error => %v", err)
	}

	err = service.Deposit(account.ID, 100_00)
	if err != nil {
		t.Errorf("error => %v", err)
	}

	payment, err := service.Pay(account.ID, 10_00, "Food")

	if err != nil {
		t.Errorf("error => %v", err)
	}

	pay, err := service.FindPaymentByID(payment.ID)

	if err != nil {
		t.Errorf("error => %v", err)
	}

	err = service.Reject(pay.ID)

	if err != nil {
		t.Errorf("error => %v", err)
	}

}

func TestService_Reject_fail_user(t *testing.T) {
	var service Service
	service.RegisterAccount("9127660305")
	account, err := service.FindAccountByID(1)

	if err != nil {
		t.Errorf("account => %v", account)
	}

	err = service.Deposit(account.ID, 100_00)
	if err != nil {
		t.Errorf("error => %v", err)
	}

	payment, err := service.Pay(account.ID, 10_00, "Food")

	if err != nil {
		t.Errorf("account => %v", account)
	}

	pay, err := service.FindPaymentByID(payment.ID)

	if err != nil {
		t.Errorf("payment => %v", payment)
	}

	err = service.Reject(pay.ID + "uu")

	if err == nil {
		t.Errorf("pay => %v", pay)
	}

}
