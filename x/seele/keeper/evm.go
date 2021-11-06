package keeper

import (
	"errors"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tharsis/ethermint/server/config"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/Seele-N/Seele/x/seele/types"
)

// CallEVM execute an evm message from native module
func (k Keeper) CallEVM(ctx sdk.Context, to *common.Address, data []byte, value *big.Int) (*ethtypes.Message, *evmtypes.MsgEthereumTxResponse, error) {
	k.evmKeeper.WithContext(ctx)

	nonce := k.evmKeeper.GetNonce(types.EVMModuleAddress)
	msg := ethtypes.NewMessage(
		types.EVMModuleAddress,
		to,
		nonce,
		value, // amount
		config.DefaultGasCap,
		big.NewInt(0), // gasPrice
		data,
		nil,   // accessList
		false, // checkNonce
	)

	params := k.evmKeeper.GetParams(ctx)
	// return error if contract creation or call are disabled through governance
	if !params.EnableCreate && to == nil {
		return nil, nil, errors.New("failed to create new contract")
	} else if !params.EnableCall && to != nil {
		return nil, nil, errors.New("failed to call contract")
	}
	ethCfg := params.ChainConfig.EthereumConfig(k.evmKeeper.ChainID())

	// get the coinbase address from the block proposer
	coinbase, err := k.evmKeeper.GetCoinbaseAddress(ctx)
	if err != nil {
		return nil, nil, errors.New("failed to obtain coinbase address")
	}
	evm := k.evmKeeper.NewEVM(msg, ethCfg, params, coinbase, types.NewDummyTracer())
	ret, err := k.evmKeeper.ApplyMessage(evm, msg, ethCfg, true)
	if err != nil {
		return nil, nil, err
	}
	k.evmKeeper.CommitCachedContexts()
	return &msg, ret, nil
}

// CallModuleSRC20 call a method of ModuleSRC20 contract
func (k Keeper) CallModuleSRC20(ctx sdk.Context, contract common.Address, method string, args ...interface{}) ([]byte, error) {
	data, err := types.ModuleSRC20Contract.ABI.Pack(method, args...)
	if err != nil {
		return nil, err
	}
	_, res, err := k.CallEVM(ctx, &contract, data, big.NewInt(0))
	if err != nil {
		return nil, err
	}
	if res.Failed() {
		return nil, fmt.Errorf("call contract failed: %s, %s, %s", contract.Hex(), method, res.Ret)
	}
	return res.Ret, nil
}

// DeploySnpDelegate deploy an embed snp delegate contract
func (k Keeper) DeploySnpDelegate(ctx sdk.Context) (common.Address, error) {
	ctor, err := types.SnpDelegateContract.ABI.Pack("")
	if err != nil {
		return common.Address{}, err
	}
	data := types.SnpDelegateContract.Bin
	data = append(data, ctor...)

	msg, res, err := k.CallEVM(ctx, nil, data, big.NewInt(0))
	if err != nil {
		return common.Address{}, err
	}

	if res.Failed() {
		return common.Address{}, fmt.Errorf("contract deploy failed: %s", res.Ret)
	}
	return crypto.CreateAddress(types.EVMModuleAddress, msg.Nonce()), nil
}

// DeployModuleSRC20 deploy an embed erc20 contract
func (k Keeper) DeployModuleSRC20(ctx sdk.Context, denom string) (common.Address, error) {
	ctor, err := types.ModuleSRC20Contract.ABI.Pack("", denom+" Token", denom, uint8(18))
	if err != nil {
		return common.Address{}, err
	}
	data := types.ModuleSRC20Contract.Bin
	data = append(data, ctor...)

	msg, res, err := k.CallEVM(ctx, nil, data, big.NewInt(0))
	if err != nil {
		return common.Address{}, err
	}

	if res.Failed() {
		return common.Address{}, fmt.Errorf("contract deploy failed: %s", res.Ret)
	}
	return crypto.CreateAddress(types.EVMModuleAddress, msg.Nonce()), nil
}

// ConvertCoinFromNativeToSRC20 convert native token to erc20 token
func (k Keeper) ConvertCoinFromNativeToSRC20(ctx sdk.Context, sender common.Address, coin sdk.Coin, autoDeploy bool) error {
	if !types.IsValidDenomToWrap(coin.Denom) {
		return fmt.Errorf("coin %s is not supported for wrapping", coin.Denom)
	}

	var err error
	// external contract is returned in preference to auto-deployed ones
	contract, found := k.GetContractByDenom(ctx, coin.Denom)
	if !found {
		if !autoDeploy {
			return fmt.Errorf("no contract found for the denom %s", coin.Denom)
		}
		contract, err = k.DeployModuleSRC20(ctx, coin.Denom)
		if err != nil {
			return err
		}
		k.SetAutoContractForDenom(ctx, coin.Denom, contract)

		k.Logger(ctx).Info(fmt.Sprintf("contract address %s created for coin denom %s", contract.String(), coin.Denom))

		if coin.Denom == "snp" {
			contractSnpDelegate, err := k.DeploySnpDelegate(ctx)
			if err != nil {
				return err
			}

			k.SetContractForContractName(ctx, types.SnpDelegateContract.ContractName, contractSnpDelegate)

			k.Logger(ctx).Info(fmt.Sprintf("contract address %s created name %s", contractSnpDelegate.String(), types.SnpDelegateContract.ContractName))

		}
	}
	err = k.bankKeeper.SendCoins(ctx, sdk.AccAddress(sender.Bytes()), sdk.AccAddress(contract.Bytes()), sdk.NewCoins(coin))
	if err != nil {
		return err
	}
	_, err = k.CallModuleSRC20(ctx, contract, "mint_by_seele_module", sender, coin.Amount.BigInt())
	if err != nil {
		return err
	}

	return nil
}

// ConvertCoinFromSRC20ToNative convert erc20 token to native token
func (k Keeper) ConvertCoinFromSRC20ToNative(ctx sdk.Context, contract common.Address, receiver common.Address, amount sdk.Int) error {
	denom, found := k.GetDenomByContract(ctx, contract)
	if !found {
		return fmt.Errorf("the contract address %s is not mapped to native token", contract.String())
	}

	err := k.bankKeeper.SendCoins(
		ctx,
		sdk.AccAddress(contract.Bytes()),
		sdk.AccAddress(receiver.Bytes()),
		sdk.NewCoins(sdk.NewCoin(denom, amount)),
	)
	if err != nil {
		return err
	}

	_, err = k.CallModuleSRC20(ctx, contract, "burn_by_seele_module", receiver, amount.BigInt())
	if err != nil {
		return err
	}

	return nil
}

// ConvertCoinsFromNativeToSRC20 convert native tokens to erc20 tokens
func (k Keeper) ConvertCoinsFromNativeToSRC20(ctx sdk.Context, sender common.Address, coins sdk.Coins, autoDeploy bool) error {
	for _, coin := range coins {
		if err := k.ConvertCoinFromNativeToSRC20(ctx, sender, coin, autoDeploy); err != nil {
			return err
		}
	}
	return nil
}
