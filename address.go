package dcrlibwallet

import (
	w "decred.org/dcrwallet/wallet"
	"github.com/decred/dcrd/dcrutil/v3"
	"github.com/decred/dcrwallet/errors/v2"
)

// AddressInfo holds information about an address
// If the address belongs to the querying wallet, IsMine will be true and the AccountNumber and AccountName values will be populated
type AddressInfo struct {
	Address       string
	IsMine        bool
	AccountNumber uint32
	AccountName   string
}

func (mw *MultiWallet) IsAddressValid(address string) bool {
	_, err := dcrutil.DecodeAddress(address, mw.chainParams)
	return err == nil
}

func (wallet *Wallet) HaveAddress(address string) bool {
	addr, err := dcrutil.DecodeAddress(address, wallet.chainParams)
	if err != nil {
		return false
	}

	have, err := wallet.internal.HaveAddress(wallet.shutdownContext(), addr)
	if err != nil {
		return false
	}

	return have
}

func (wallet *Wallet) AddressInfo(address string) (*AddressInfo, error) {
	addr, err := dcrutil.DecodeAddress(address, wallet.chainParams)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	addressInfo := new(AddressInfo)
	info, err := wallet.internal.KnownAddress(wallet.shutdownContext(), addr)
	if info != nil {
		n, err := wallet.AccountNumber(info.AccountName())
		if err != nil {
			return nil, err
		}
		addressInfo.Address = address
		addressInfo.IsMine = true
		addressInfo.AccountNumber = n
		addressInfo.AccountName = info.AccountName()
	}

	return addressInfo, err
}

func (wallet *Wallet) CurrentAddress(account int32) (string, error) {
	if wallet.IsRestored && !wallet.HasDiscoveredAccounts {
		return "", errors.E(ErrAddressDiscoveryNotDone)
	}

	addr, err := wallet.internal.CurrentAddress(uint32(account))
	if err != nil {
		log.Error(err)
		return "", err
	}
	return addr.Address(), nil
}

func (wallet *Wallet) NextAddress(account int32) (string, error) {
	if wallet.IsRestored && !wallet.HasDiscoveredAccounts {
		return "", errors.E(ErrAddressDiscoveryNotDone)
	}

	addr, err := wallet.internal.NewExternalAddress(wallet.shutdownContext(), uint32(account), w.WithGapPolicyWrap())
	if err != nil {
		log.Error(err)
		return "", err
	}
	return addr.Address(), nil
}

//AddressPubKey decodes an address to a public key.
func (wallet *Wallet) AddressPubKey(address string) (string, error) {
	addr, err := dcrutil.DecodeAddress(address, wallet.chainParams)
	if err != nil {
		return "", err
	}
	switch addr := addr.(type) {
	case *dcrutil.AddressSecpPubKey:
		return addr.String(), nil
	}

	a, err := wallet.internal.KnownAddress(wallet.shutdownContext(), addr)
	if err != nil {
		if errors.Is(err, errors.NotExist) {
			return "", translateError(err)
		}
		return "", err
	}
	var pubKey []byte
	switch a := a.(type) {
	case w.PubKeyHashAddress:
		pubKey = a.PubKey()
	default:
		err = errors.New("address has no associated public key")
		return "", err
	}
	pubKeyAddr, err := dcrutil.NewAddressSecpPubKey(pubKey, wallet.chainParams)
	if err != nil {
		return "", err
	}
	return pubKeyAddr.String(), nil
}
