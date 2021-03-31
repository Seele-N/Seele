package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// AccountAddressPrefix Account Address Prefix
	AccountAddressPrefix = "seele"

	// CoinType seele-n in https://github.com/satoshilabs/slips/blob/master/slip-0044.md
	CoinType = 513

	// FullFundraiserPath BIP44Prefix is the parts of the BIP44 HD path that are fixed by
	// what we used during the fundraiser.
	FullFundraiserPath = "44'/513'/0'/0/0"
)

var (
	// AccountPubKeyPrefix Account PubKey Prefix
	AccountPubKeyPrefix = AccountAddressPrefix + "pub"
	// ValidatorAddressPrefix Validator Address Prefix
	ValidatorAddressPrefix = AccountAddressPrefix + "valoper"
	// ValidatorPubKeyPrefix Validator PubKey Prefix
	ValidatorPubKeyPrefix = AccountAddressPrefix + "valoperpub"
	// ConsNodeAddressPrefix Node Address Prefix
	ConsNodeAddressPrefix = AccountAddressPrefix + "valcons"
	// ConsNodePubKeyPrefix Node PubKey Prefix
	ConsNodePubKeyPrefix = AccountAddressPrefix + "valconspub"
)

// SetConfig set chain config
func SetConfig() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddressPrefix, ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsNodeAddressPrefix, ConsNodePubKeyPrefix)
	config.SetCoinType(CoinType)
	config.SetFullFundraiserPath(FullFundraiserPath)
	config.Seal()
}
