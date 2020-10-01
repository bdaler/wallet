package wallet

import (
	"reflect"
	"testing"
	"wallet/pkg/types"
)

func TestService_Deposit(t *testing.T) {
	type fields struct {
		nextAccountID int64
		accounts      []*types.Account
		payments      []*types.Payment
	}
	type args struct {
		accountID int64
		amount    types.Money
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				nextAccountID: tt.fields.nextAccountID,
				accounts:      tt.fields.accounts,
				payments:      tt.fields.payments,
			}
			if err := s.Deposit(tt.args.accountID, tt.args.amount); (err != nil) != tt.wantErr {
				t.Errorf("Deposit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_Pay(t *testing.T) {
	type fields struct {
		nextAccountID int64
		accounts      []*types.Account
		payments      []*types.Payment
	}
	type args struct {
		accountID int64
		amount    types.Money
		category  types.PaymentCategory
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Payment
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				nextAccountID: tt.fields.nextAccountID,
				accounts:      tt.fields.accounts,
				payments:      tt.fields.payments,
			}
			got, err := s.Pay(tt.args.accountID, tt.args.amount, tt.args.category)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pay() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pay() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_RegisterAccount(t *testing.T) {
	type fields struct {
		nextAccountID int64
		accounts      []*types.Account
		payments      []*types.Payment
	}
	type args struct {
		phone types.Phone
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Account
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				nextAccountID: tt.fields.nextAccountID,
				accounts:      tt.fields.accounts,
				payments:      tt.fields.payments,
			}
			got, err := s.RegisterAccount(tt.args.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegisterAccount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_FindAccountByID(t *testing.T) {
	//var accs *[]types.Account
	//accs = &[]types.Account{
	//types.Account{
	//	ID:      10,
	//	Phone:   "9127660305",
	//	Balance: 0,
	//},
	//types.Account{
	//	ID:      10,
	//	Phone:   "9127660305",
	//	Balance: 0,
	//},
	//}
	type fields struct {
		nextAccountID int64
		accounts      []*types.Account
		payments      []*types.Payment
	}
	type args struct {
		accountID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Account
		wantErr bool
	}{
		//{name: testing.CoverMode(), fields: fields{}, args: args{}, want: &types.Account{ID:      10, Phone:   "9127660305", Balance: 0}, wantErr: true},
		//{name: testing.CoverMode(), fields: fields{nextAccountID: 0, payments:      nil}, args: args{
		//	accountID: 10,
		//}, want: &types.Account{
		//	ID:      10,
		//	Phone:   "9127660305",
		//	Balance: 0,
		//}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				nextAccountID: tt.fields.nextAccountID,
				accounts:      tt.fields.accounts,
				payments:      tt.fields.payments,
			}
			got, err := s.FindAccountByID(tt.args.accountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAccountByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindAccountByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
