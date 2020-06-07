package KeyVault

import (
	"fmt"
	core "github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/wallet_hd"
	"github.com/google/uuid"
)

func (portfolio *KeyVault) ID() uuid.UUID {
	return portfolio.id
}

// CreateAccount creates a new account in the wallet.
// This will error if an account with the name already exists.
// Will push to the new wallet the lock policy
func (portfolio *KeyVault) CreateWallet(name string) (core.Wallet, error) {
	// create wallet
	id := len(portfolio.indexMapper)
	path := fmt.Sprintf("/%d",id)
	key,err := portfolio.key.Derive(path)
	if err != nil {
		return nil,err
	}
	retWallet := wallet_hd.NewHDWallet(name,
		key,
		path,
		portfolio.context,
	)

	// register new wallet and save portfolio + wallet
	reset := func() {
		delete(portfolio.indexMapper,name)
	}
	portfolio.indexMapper[name] = retWallet.ID()
	err = portfolio.context.Storage.SaveWallet(retWallet)
	if err != nil {
		reset()
		return nil,err
	}
	err = portfolio.context.Storage.SavePortfolio(portfolio)
	if err != nil {
		reset()
		return nil,err
	}

	return retWallet,nil
}

// Accounts provides all accounts in the wallet.
func (portfolio *KeyVault) Wallets() <-chan core.Wallet {
	ch := make (chan core.Wallet,1024) // TODO - handle more? change from chan?
	go func() {
		for name := range portfolio.indexMapper {
			id := portfolio.indexMapper[name]
			wallet,err := portfolio.WalletByID(id)
			if err != nil {
				continue
			}
			ch <- wallet
		}
	}()

	return ch
}

// AccountByID provides a single account from the wallet given its ID.
// This will error if the account is not found.
func (portfolio *KeyVault) WalletByID(id uuid.UUID) (core.Wallet, error) {
	return portfolio.context.Storage.OpenWallet(id)
}

// AccountByName provides a single account from the wallet given its name.
// This will error if the account is not found.
func (portfolio *KeyVault) WalletByName(name string) (core.Wallet, error) {
	id,exists := portfolio.indexMapper[name]
	if !exists {
		return nil, fmt.Errorf("no wallet found")
	}

	return portfolio.WalletByID(id)
}