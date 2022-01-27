package types

const (
	// ModuleName defines the module name
	ModuleName = "seele"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_seele"

	// this line is used by starport scaffolding # ibc/keys/name
)

// prefix bytes for the seele persistent store
const (
	prefixDenomToExternalContract = iota + 1
	prefixDenomToAutoContract
	prefixContractToDenom
	prefixContractNameToContractAddress
	prefixExternalContractToDenom
)

// KVStore key prefixes
var (
	KeyPrefixDenomToExternalContract       = []byte{prefixDenomToExternalContract}
	KeyPrefixDenomToAutoContract           = []byte{prefixDenomToAutoContract}
	KeyPrefixContractToDenom               = []byte{prefixContractToDenom}
	KeyPrefixContractNameToContractAddress = []byte{prefixContractNameToContractAddress}
	KeyprefixExternalContractToDenom       = []byte{prefixExternalContractToDenom}
)

// this line is used by starport scaffolding # ibc/keys/port

// DenomToExternalContractKey defines the store key for denom to contract mapping
func DenomToExternalContractKey(denom string) []byte {
	return append(KeyPrefixDenomToExternalContract, denom...)
}

// ExternalContractToDenomKey defines the store key for External contract to denom mapping
func ExternalContractToDenomKey(contract string) []byte {
	return append(KeyprefixExternalContractToDenom, contract...)
}

// DenomToAutoContractKey defines the store key for denom to auto contract mapping
func DenomToAutoContractKey(denom string) []byte {
	return append(KeyPrefixDenomToAutoContract, denom...)
}

// ContractToDenomKey defines the store key for contract to denom reverse index
func ContractToDenomKey(contract []byte) []byte {
	return append(KeyPrefixContractToDenom, contract...)
}

// ContractNameToContractAddressKey defines the store key for contractname to auto contract address mapping
func ContractNameToContractAddressKey(contractname string) []byte {
	return append(KeyPrefixContractNameToContractAddress, contractname...)
}
