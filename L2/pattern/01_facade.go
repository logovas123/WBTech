package main

import "fmt"

/*
Если наша система состоит из множества подсистем, то она может иметь большую сложность и в ней можно запутаться.
Паттерн Фасад позволяет работать со всеми компонентами системы, имея при этом один интерфейс.

В примере ниже у нас система, которая, например, кладёт деньги на карту. Её логика состоит из подсистем:
- проверка имени аккаунта (подсистема Acoount)
- проверка введеного пин кода карты (подсистема SecurityCode)
- баланс карты (добавить количество денег к балансу, подсистема Wallet)

Все подсистемы мы обьединяем под общей структурой WalletFacade. И теперь можем управлять всеми подсистемами из этого общего "фасада",
и не задумываться о внутренней реализации каждой подсистемы
*/

type WalletFacade struct {
	account      *Account
	securityCode *SecurityCode
	wallet       *Wallet
}

func NewWalletFacade(accountName string, code int) *WalletFacade {
	fmt.Println("Start create account...")

	wf := &WalletFacade{
		account:      NewAccount(accountName),
		securityCode: NewSecurityCode(code),
		wallet:       NewWallet(),
	}

	fmt.Println("Account create")

	return wf
}

func (wf *WalletFacade) addMonyeToWallet(acccountName string, code int, amount int) error {
	fmt.Println("Process of add money to wallet...")

	err := wf.account.CheckAccountName(acccountName)
	if err != nil {
		return err
	}

	err = wf.securityCode.CheckCode(code)
	if err != nil {
		return err
	}

	wf.wallet.AddMoney(amount)

	return nil
}

// -----
type Account struct {
	name string
}

func NewAccount(accauntName string) *Account {
	return &Account{name: accauntName}
}

func (a *Account) CheckAccountName(accountName string) error {
	if a.name != accountName {
		return fmt.Errorf("Account name incorrect")
	}

	fmt.Println("Success check account name: name correct")

	return nil
}

// -----
type SecurityCode struct {
	code int
}

func NewSecurityCode(code int) *SecurityCode {
	return &SecurityCode{code: code}
}

func (s *SecurityCode) CheckCode(code int) error {
	if s.code != code {
		return fmt.Errorf("Code incorrect")
	}

	fmt.Println("Success check code: code correct")

	return nil
}

// -----
type Wallet struct {
	balance int
}

func NewWallet() *Wallet {
	return &Wallet{
		balance: 0,
	}
}

func (w *Wallet) AddMoney(amountMoney int) {
	w.balance += amountMoney
	fmt.Println("Success add money to wallet")
}

// -----
/*
func main() {
	fmt.Println()
	walletFacade := NewWalletFacade("john", 1234)
	fmt.Println()

	err := walletFacade.addMonyeToWallet("john", 1234, 10) // успешно добавляем деньги
	if err != nil {
		fmt.Println("Error add money to wallet:", err)
	}

	fmt.Println()

	err = walletFacade.addMonyeToWallet("john", 134, 10) // ошибка добавления - неверный пин-код
	if err != nil {
		fmt.Println("Error add money to wallet:", err)
	}

	fmt.Println()
}
*/
