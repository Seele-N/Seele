package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// AccountAddressPrefix Account Address Prefix
	AccountAddressPrefix = "seele"
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
	config.Seal()
}
