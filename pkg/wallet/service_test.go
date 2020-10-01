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
