package keeper_test

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
)

func (suite *KeeperTestSuite) TestDeployContract() {
	suite.SetupTest()
	keeper := suite.app.SeeleKeeper

	_, err := keeper.DeployModuleSRC20(suite.ctx, "test", uint8(18))
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) TestTokenConversion() {
	suite.SetupTest()
	keeper := suite.app.SeeleKeeper

	// generate test address
	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	address := common.BytesToAddress(priv.PubKey().Address().Bytes())
	cosmosAddress := sdk.AccAddress(address.Bytes())

	denom := "ibc/0000000000000000000000000000000000000000000000000000000000000000"
	amount := big.NewInt(100)
	coins := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewIntFromBigInt(amount)))

	// mint native tokens
	err = suite.MintCoins(sdk.AccAddress(address.Bytes()), coins)
	suite.Require().NoError(err)

	// send to erc20
	err = keeper.ConvertCoinsFromNativeToSRC20(suite.ctx, "", address, coins, true)
	suite.Require().NoError(err)

	// check erc20 balance
	contract, found := keeper.GetContractByDenom(suite.ctx, denom)
	suite.Require().True(found)

	ret, err := keeper.CallModuleSRC20(suite.ctx, contract, "balanceOf", address)
	suite.Require().NoError(err)
	suite.Require().Equal(amount, big.NewInt(0).SetBytes(ret))

	ret, err = keeper.CallModuleSRC20(suite.ctx, contract, "totalSupply")
	suite.Require().NoError(err)
	suite.Require().Equal(amount, big.NewInt(0).SetBytes(ret))

	// convert back to native
	err = keeper.ConvertCoinFromSRC20ToNative(suite.ctx, contract, address, coins[0].Amount)
	suite.Require().NoError(err)

	ret, err = keeper.CallModuleSRC20(suite.ctx, contract, "balanceOf", address)
	suite.Require().NoError(err)
	suite.Require().Equal(0, big.NewInt(0).Cmp(big.NewInt(0).SetBytes(ret)))

	ret, err = keeper.CallModuleSRC20(suite.ctx, contract, "totalSupply")
	suite.Require().NoError(err)
	suite.Require().Equal(0, big.NewInt(0).Cmp(big.NewInt(0).SetBytes(ret)))

	// native balance recovered
	coin := suite.app.BankKeeper.GetBalance(suite.ctx, cosmosAddress, denom)
	suite.Require().Equal(amount, coin.Amount.BigInt())
}
