package types

import (
	// embed compiled smart contract
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// ByteString is a byte array that serializes to hex
type ByteString []byte

// MarshalJSON serializes ByteArray to hex
func (s ByteString) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(fmt.Sprintf("%x", string(s)))
	return bytes, err
}

// UnmarshalJSON deserializes ByteArray to hex
func (s *ByteString) UnmarshalJSON(data []byte) error {
	var x string
	err := json.Unmarshal(data, &x)
	if err == nil {
		str, e := hex.DecodeString(x)
		*s = str
		err = e
	}

	return err
}

// CompiledContract contains compiled bytecode and abi
type CompiledContract struct {
	ContractName string
	ABI          abi.ABI
	Bin          ByteString
}

const EVMModuleName = "seele-evm"

var (
	//go:embed contracts/ModuleSRC20.json
	seeleERC20JSON []byte

	// ModuleSRC20Contract is the compiled seele erc20 contract
	ModuleSRC20Contract CompiledContract

	//go:embed contracts/SnpDelegate.json
	snpDelegateJSON []byte

	// SnpDelegateContract is the compiled Snp Delegate contract
	SnpDelegateContract CompiledContract

	// EVMModuleAddress is the native module address for EVM
	EVMModuleAddress common.Address
)

func init() {
	EVMModuleAddress = common.BytesToAddress(authtypes.NewModuleAddress(EVMModuleName).Bytes())
	/*
		fmt.Printf("seele-evm module address:%s\n", EVMModuleAddress.String())
		add := crypto.CreateAddress(EVMModuleAddress, 0)
		fmt.Printf("noce=0;seele-evm module contract address:%s\n", add.String())
		add = crypto.CreateAddress(EVMModuleAddress, 1)
		fmt.Printf("noce=1;seele-evm module contract address:%s\n", add.String())
		add = crypto.CreateAddress(EVMModuleAddress, 2)
		fmt.Printf("noce=2;seele-evm module contract address:%s\n", add.String())
	*/
	err := json.Unmarshal(seeleERC20JSON, &ModuleSRC20Contract)
	if err != nil {
		panic(err)
	}

	if len(ModuleSRC20Contract.Bin) == 0 {
		panic("load src20 contract failed")
	}

	err = json.Unmarshal(snpDelegateJSON, &SnpDelegateContract)
	if err != nil {
		panic(err)
	}

	if len(SnpDelegateContract.Bin) == 0 {
		panic("load snp delegate contract failed")
	}
}
