package keeper_test

import (
	"errors"

	"github.com/Seele-N/Seele/app"
	seelemodulekeeper "github.com/Seele-N/Seele/x/seele/keeper"
	keepertest "github.com/Seele-N/Seele/x/seele/keeper/mock"
	"github.com/Seele-N/Seele/x/seele/types"
)

func (suite *KeeperTestSuite) TestGetSourceChannelID() {

	testCases := []struct {
		name          string
		ibcDenom      string
		expectedError error
		postCheck     func(channelID string)
	}{
		{
			"wrong ibc denom",
			"test",
			errors.New("test is invalid: ibc cro denom is invalid"),
			func(channelID string) {},
		},
		{
			"correct ibc denom",
			types.IbcCroDenomDefaultValue,
			nil,
			func(channelID string) {
				suite.Require().Equal(channelID, "channel-0")
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset
			// Create Seele Keeper with mock transfer keeper
			seeleKeeper := *seelemodulekeeper.NewKeeper(
				app.MakeEncodingConfig().Marshaler,
				suite.app.GetKey(types.StoreKey),
				suite.app.GetKey(types.MemStoreKey),
				suite.app.GetSubspace(types.ModuleName),
				suite.app.BankKeeper,
				keepertest.IbcKeeperMock{},
				suite.app.EvmKeeper,
			)
			suite.app.SeeleKeeper = seeleKeeper

			channelId, err := suite.app.SeeleKeeper.GetSourceChannelID(suite.ctx, tc.ibcDenom)
			if tc.expectedError != nil {
				suite.Require().EqualError(err, tc.expectedError.Error())
			} else {
				suite.Require().NoError(err)
				tc.postCheck(channelId)
			}
		})
	}
}
