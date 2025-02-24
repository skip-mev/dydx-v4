package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestCompleteBridge(t *testing.T) {
	tests := map[string]struct {
		// Initial balance of bridge module account.
		initialBalance sdk.Coin
		// Bridge event to complete.
		bridgeEvent types.BridgeEvent
		// Whether bridging is disabled.
		bridgingDisabled bool

		// Expected error, if any.
		expectedError string
		// Expected balance of bridge module account after bridge completion.
		expectedBalance sdk.Coin
	}{
		"Success": {
			initialBalance:  sdk.NewCoin("adv4tnt", sdkmath.NewInt(1_000)),
			bridgeEvent:     constants.BridgeEvent_Id0_Height0,           // bridges 888 tokens.
			expectedBalance: sdk.NewCoin("adv4tnt", sdkmath.NewInt(112)), // 1000 - 888.
		},
		"Failure: coin amount is 0": {
			initialBalance: sdk.NewCoin("adv4tnt", sdkmath.NewInt(1_000)),
			bridgeEvent: types.BridgeEvent{
				Id:             7,
				Address:        constants.BobAccAddress.String(),
				Coin:           sdk.NewCoin("adv4tnt", sdkmath.ZeroInt()),
				EthBlockHeight: 3,
			},
			expectedError:   "invalid coin",
			expectedBalance: sdk.NewCoin("adv4tnt", sdkmath.NewInt(1_000)),
		},
		"Failure: coin amount is negative": {
			initialBalance: sdk.NewCoin("adv4tnt", sdkmath.NewInt(1_000)),
			bridgeEvent: types.BridgeEvent{
				Id:      7,
				Address: constants.BobAccAddress.String(),
				Coin: sdk.Coin{
					Denom:  "adv4tnt",
					Amount: sdkmath.NewInt(-1),
				},
				EthBlockHeight: 3,
			},
			expectedError:   "invalid coin",
			expectedBalance: sdk.NewCoin("adv4tnt", sdkmath.NewInt(1_000)),
		},
		"Failure: invalid address string": {
			initialBalance: sdk.NewCoin("adv4tnt", sdkmath.NewInt(1_000)),
			bridgeEvent: types.BridgeEvent{
				Id:             4,
				Address:        "not an address string",
				Coin:           sdk.NewCoin("adv4tnt", sdkmath.NewInt(1)),
				EthBlockHeight: 2,
			},
			expectedError:   "decoding bech32 failed",
			expectedBalance: sdk.NewCoin("adv4tnt", sdkmath.NewInt(1_000)),
		},
		"Failure: bridge module account has insufficient balance": {
			initialBalance:  sdk.NewCoin("adv4tnt", sdkmath.NewInt(500)),
			bridgeEvent:     constants.BridgeEvent_Id0_Height0, // bridges 888 tokens.
			expectedError:   "insufficient funds",
			expectedBalance: sdk.NewCoin("adv4tnt", sdkmath.NewInt(500)),
		},
		"Failure: bridging is disabled": {
			initialBalance:   sdk.NewCoin("adv4tnt", sdkmath.NewInt(1_000)),
			bridgeEvent:      constants.BridgeEvent_Id0_Height0, // bridges 888 tokens.
			bridgingDisabled: true,
			expectedError:    types.ErrBridgingDisabled.Error(),
			expectedBalance:  sdk.NewCoin("adv4tnt", sdkmath.NewInt(1_000)), // same as initial.
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize context and keeper.
			ctx, bridgeKeeper, _, _, _, bankKeeper, _ := keepertest.BridgeKeepers(t)
			err := bridgeKeeper.UpdateSafetyParams(ctx, types.SafetyParams{
				IsDisabled:  tc.bridgingDisabled,
				DelayBlocks: bridgeKeeper.GetSafetyParams(ctx).DelayBlocks,
			})
			require.NoError(t, err)
			// Fund bridge module account with enought balance.
			err = bankKeeper.MintCoins(
				ctx,
				types.ModuleName,
				sdk.NewCoins(tc.initialBalance),
			)
			require.NoError(t, err)

			// Complete bridge.
			err = bridgeKeeper.CompleteBridge(ctx, tc.bridgeEvent)

			// Assert expectations.
			if tc.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError)
			} else {
				require.NoError(t, err)

				// Assert that target account's balance of bridged token is as expected.
				balance := bankKeeper.GetBalance(
					ctx,
					sdk.MustAccAddressFromBech32(tc.bridgeEvent.Address),
					tc.bridgeEvent.Coin.Denom,
				)
				require.Equal(t, tc.bridgeEvent.Coin.Denom, balance.Denom)
				require.Equal(t, tc.bridgeEvent.Coin.Amount, balance.Amount)
			}
			// Assert that bridge module account's balance is as expected.
			modAccBalance := bankKeeper.GetBalance(
				ctx,
				authtypes.NewModuleAddress(types.ModuleName),
				tc.bridgeEvent.Coin.Denom,
			)
			require.Equal(t, tc.expectedBalance.Denom, modAccBalance.Denom)
			require.Equal(t, tc.expectedBalance.Amount, modAccBalance.Amount)
		})
	}
}
