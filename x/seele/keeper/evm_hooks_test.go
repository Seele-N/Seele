package keeper_test

import (
	"fmt"
	"math/big"

	"github.com/Seele-N/Seele/app"
	keepertest "github.com/Seele-N/Seele/x/seele/keeper/mock"
	"github.com/Seele-N/Seele/x/seele/types"

	seelemodulekeeper "github.com/Seele-N/Seele/x/seele/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/Seele-N/Seele/x/seele/keeper"
)

func (suite *KeeperTestSuite) TestEvmHooks() {
	suite.SetupTest()

	contract := common.BigToAddress(big.NewInt(1))
	txHash := common.BigToHash(big.NewInt(2))
	recipient := common.BigToAddress(big.NewInt(3))
	sender := common.BigToAddress(big.NewInt(4))
	denom := "testdenom"

	testCases := []struct {
		msg      string
		malleate func()
	}{
		{
			"invalid log data, but still success",
			func() {
				logs := []*ethtypes.Log{
					{
						Address: contract,
						Topics:  []common.Hash{keeper.SendToAccountEvent.ID},
					},
				}
				err := suite.app.EvmKeeper.PostTxProcessing(txHash, logs)
				suite.Require().NoError(err)
			},
		},
		{
			"not enough balance, expect fail",
			func() {
				data, err := keeper.SendToAccountEvent.Inputs.Pack(
					recipient,
					big.NewInt(100),
				)
				suite.Require().NoError(err)
				logs := []*ethtypes.Log{
					{
						Address: contract,
						Topics:  []common.Hash{keeper.SendToAccountEvent.ID},
						Data:    data,
					},
				}
				err = suite.app.EvmKeeper.PostTxProcessing(txHash, logs)
				suite.Require().Error(err)
			},
		},
		{
			"success send to account",
			func() {
				suite.app.SeeleKeeper.SetExternalContractForDenom(suite.ctx, denom, contract)
				coin := sdk.NewCoin(denom, sdk.NewInt(100))
				err := suite.MintCoins(sdk.AccAddress(contract.Bytes()), sdk.NewCoins(coin))
				suite.Require().NoError(err)

				balance := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.AccAddress(contract.Bytes()), denom)
				suite.Require().Equal(coin, balance)

				data, err := keeper.SendToAccountEvent.Inputs.Pack(
					recipient,
					coin.Amount.BigInt(),
				)
				suite.Require().NoError(err)
				logs := []*ethtypes.Log{
					{
						Address: contract,
						Topics:  []common.Hash{keeper.SendToAccountEvent.ID},
						Data:    data,
					},
				}
				err = suite.app.EvmKeeper.PostTxProcessing(txHash, logs)
				suite.Require().NoError(err)

				balance = suite.app.BankKeeper.GetBalance(suite.ctx, sdk.AccAddress(contract.Bytes()), denom)
				suite.Require().Equal(sdk.NewCoin(denom, sdk.NewInt(0)), balance)
				balance = suite.app.BankKeeper.GetBalance(suite.ctx, sdk.AccAddress(recipient.Bytes()), denom)
				suite.Require().Equal(coin, balance)
			},
		},
		{
			"failed send to ethereum, invalid gravity denom",
			func() {
				suite.SetupTest()

				suite.app.SeeleKeeper.SetExternalContractForDenom(suite.ctx, denom, contract)
				coin := sdk.NewCoin(denom, sdk.NewInt(100))
				err := suite.MintCoins(sdk.AccAddress(contract.Bytes()), sdk.NewCoins(coin))
				suite.Require().NoError(err)

				balance := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.AccAddress(contract.Bytes()), denom)
				suite.Require().Equal(coin, balance)

				data, err := keeper.SendToEthereumEvent.Inputs.Pack(
					recipient,
					coin.Amount.BigInt(),
					big.NewInt(0),
				)
				suite.Require().NoError(err)
				logs := []*ethtypes.Log{
					{
						Address: contract,
						Topics:  []common.Hash{keeper.SendToEthereumEvent.ID},
						Data:    data,
					},
				}
				err = suite.app.EvmKeeper.PostTxProcessing(txHash, logs)
				// should fail, because of not gravity denom name
				suite.Require().Error(err)
			},
		},
		{
			"fail send to ethereum", // gravity feature is removed
			func() {
				suite.SetupTest()
				denom := "gravity0x0000000000000000000000000000000000000000"

				suite.app.SeeleKeeper.SetExternalContractForDenom(suite.ctx, denom, contract)
				coin := sdk.NewCoin(denom, sdk.NewInt(100))
				err := suite.MintCoins(sdk.AccAddress(contract.Bytes()), sdk.NewCoins(coin))
				suite.Require().NoError(err)

				balance := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.AccAddress(contract.Bytes()), denom)
				suite.Require().Equal(coin, balance)

				data, err := keeper.SendToEthereumEvent.Inputs.Pack(
					recipient,
					coin.Amount.BigInt(),
					big.NewInt(0),
				)
				suite.Require().NoError(err)
				logs := []*ethtypes.Log{
					{
						Address: contract,
						Topics:  []common.Hash{keeper.SendToEthereumEvent.ID},
						Data:    data,
					},
				}
				err = suite.app.EvmKeeper.PostTxProcessing(txHash, logs)
				suite.Require().Error(err)
			},
		},
		{
			"failed send to ibc, invalid ibc denom",
			func() {
				suite.SetupTest()
				// Create Seele Keeper with mock transfer keeper
				seeleKeeper := *seelemodulekeeper.NewKeeper(
					app.MakeEncodingConfig().Marshaler,
					suite.app.GetKey(types.StoreKey),
					suite.app.GetKey(types.MemStoreKey),
					suite.app.GetSubspace(types.ModuleName),
					suite.app.BankKeeper,
					keepertest.IbcKeeperMock{},
					suite.app.GravityKeeper,
					suite.app.EvmKeeper,
					suite.app.StakingKeeper,
				)
				suite.app.SeeleKeeper = seeleKeeper

				suite.app.SeeleKeeper.SetExternalContractForDenom(suite.ctx, denom, contract)
				coin := sdk.NewCoin(denom, sdk.NewInt(100))
				err := suite.MintCoins(sdk.AccAddress(contract.Bytes()), sdk.NewCoins(coin))
				suite.Require().NoError(err)

				balance := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.AccAddress(contract.Bytes()), denom)
				suite.Require().Equal(coin, balance)

				data, err := keeper.SendToIbcEvent.Inputs.Pack(
					sender,
					"recipient",
					coin.Amount.BigInt(),
				)
				suite.Require().NoError(err)
				logs := []*ethtypes.Log{
					{
						Address: contract,
						Topics:  []common.Hash{keeper.SendToIbcEvent.ID},
						Data:    data,
					},
				}
				err = suite.app.EvmKeeper.PostTxProcessing(txHash, logs)
				// should fail, because of not ibc denom name
				suite.Require().Error(err)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
		})
	}
}
