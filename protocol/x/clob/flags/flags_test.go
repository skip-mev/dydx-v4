package flags_test

import (
	"fmt"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/flags"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddFlagsToCommand(t *testing.T) {
	cmd := cobra.Command{}

	flags.AddClobFlagsToCmd(&cmd)
	tests := map[string]struct {
		flagName string
	}{
		fmt.Sprintf("Has %s flag", flags.MaxLiquidationOrdersPerBlock): {
			flagName: flags.MaxLiquidationOrdersPerBlock,
		},
		fmt.Sprintf("Has %s flag", flags.MevTelemetryHosts): {
			flagName: flags.MevTelemetryHosts,
		},
		fmt.Sprintf("Has %s flag", flags.MevTelemetryIdentifier): {
			flagName: flags.MevTelemetryIdentifier,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Contains(t, cmd.Flags().FlagUsages(), tc.flagName)
		})
	}
}

func TestGetFlagValuesFromOptions(t *testing.T) {
	tests := map[string]struct {
		// Parameters.
		optsMap map[string]any

		// Expectations.
		expectedMaxLiquidationOrdersPerBlock uint32
		expectedMevTelemetryHosts            []string
		expectedMevTelemetryIdentifier       string
	}{
		"Sets to default if unset": {
			expectedMaxLiquidationOrdersPerBlock: flags.DefaultMaxLiquidationOrdersPerBlock,
			expectedMevTelemetryHosts:            flags.DefaultMevTelemetryHosts,
			expectedMevTelemetryIdentifier:       flags.DefaultMevTelemetryIdentifier,
		},
		"Sets values from options with one host": {
			optsMap: map[string]any{
				flags.MaxLiquidationOrdersPerBlock: uint32(50),
				flags.MevTelemetryHosts:      "https://localhost:13137",
				flags.MevTelemetryIdentifier:       "node-agent-01",
			},
			expectedMaxLiquidationOrdersPerBlock: uint32(50),
			expectedMevTelemetryHosts:      []string{"https://localhost:13137"},
			expectedMevTelemetryIdentifier:       "node-agent-01",
        },
		"Sets values from options with multiple hosts": {
			optsMap: map[string]any{
				flags.MaxLiquidationOrdersPerBlock: uint32(50),
				flags.MevTelemetryHosts:      "https://localhost:13137,https://localhost:13337,https://localtest:13537",
				flags.MevTelemetryIdentifier: "node-agent-01",
			},
            expectedMaxLiquidationOrdersPerBlock: uint32(50),
			expectedMevTelemetryHosts: []string{
				"https://localhost:13137", "https://localhost:13337",
				"https://localtest:13537",
			},
			expectedMevTelemetryIdentifier: "node-agent-01",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockOpts := mocks.AppOptions{}
			mockOpts.On("Get", mock.AnythingOfType("string")).
				Return(func(key string) interface{} {
					return tc.optsMap[key]
				})

			flags := flags.GetClobFlagValuesFromOptions(&mockOpts)
			require.Equal(
				t,
				tc.expectedMevTelemetryHosts,
				flags.MevTelemetryHosts,
			)
			require.Equal(
				t,
				tc.expectedMevTelemetryIdentifier,
				flags.MevTelemetryIdentifier,
			)
			require.Equal(
				t,
				tc.expectedMaxLiquidationOrdersPerBlock,
				flags.MaxLiquidationOrdersPerBlock,
			)
		})
	}
}
