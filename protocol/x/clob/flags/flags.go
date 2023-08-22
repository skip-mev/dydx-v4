package flags

import (
	"fmt"
	"strings"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cobra"
)

// A struct containing the values of all flags.
type ClobFlags struct {
	MaxLiquidationOrdersPerBlock uint32

	MevTelemetryEnabled    bool
	MevTelemetryHosts      []string
	MevTelemetryIdentifier string
}

// List of CLI flags.
const (
	// Liquidations.
	MaxLiquidationOrdersPerBlock = "max-liquidation-orders-per-block"

	// Mev.
	MevTelemetryEnabled    = "mev-telemetry-enabled"
	MevTelemetryHosts      = "mev-telemetry-hosts"
	MevTelemetryIdentifier = "mev-telemetry-identifier"
)

// Default values.

const (
	DefaultMaxLiquidationOrdersPerBlock = 35

	DefaultMevTelemetryEnabled     = false
    DefaultMevTelemetryHostsFlag   = ""
	DefaultMevTelemetryIdentifier  = ""
)

var DefaultMevTelemetryHosts = []string{}

// AddFlagsToCmd adds flags to app initialization.
// These flags should be applied to the `start` command of the V4 Cosmos application.
// E.g. `dydxprotocold start --non-validating-full-node true`.
func AddClobFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().Uint32(
		MaxLiquidationOrdersPerBlock,
		DefaultMaxLiquidationOrdersPerBlock,
		fmt.Sprintf(
			"Sets the maximum number of liquidation orders to process per block. Default = %d",
			DefaultMaxLiquidationOrdersPerBlock,
		),
	)
	cmd.Flags().Bool(
		MevTelemetryEnabled,
		DefaultMevTelemetryEnabled,
		"Runs the MEV Telemetry collection agent if true.",
	)
	cmd.Flags().String(
		MevTelemetryHosts,
		DefaultMevTelemetryHostsFlag,
		"Sets the addresses (comma-delimited) to connect to the MEV Telemetry collection agents.",
	)
	cmd.Flags().String(
		MevTelemetryIdentifier,
		DefaultMevTelemetryIdentifier,
		"Sets the identifier to use for MEV Telemetry collection agents.",
	)
}

func GetDefaultClobFlags() ClobFlags {
	return ClobFlags{
		MaxLiquidationOrdersPerBlock: DefaultMaxLiquidationOrdersPerBlock,
		MevTelemetryEnabled:          DefaultMevTelemetryEnabled,
		MevTelemetryHosts:            DefaultMevTelemetryHosts,
		MevTelemetryIdentifier:       DefaultMevTelemetryIdentifier,
	}
}

// GetFlagValuesFromOptions gets values from the `AppOptions` struct which contains values
// from the command-line flags.
func GetClobFlagValuesFromOptions(
	appOpts servertypes.AppOptions,
) ClobFlags {
	// Create default result.
	result := GetDefaultClobFlags()

	// Populate the flags if they exist.
	if v, ok := appOpts.Get(MevTelemetryEnabled).(bool); ok {
		result.MevTelemetryEnabled = v
	}

	if v, ok := appOpts.Get(MevTelemetryHosts).(string); ok {
		result.MevTelemetryHosts = strings.Split(v, ",")
	}

	if v, ok := appOpts.Get(MevTelemetryIdentifier).(string); ok {
		result.MevTelemetryIdentifier = v
	}

	if v, ok := appOpts.Get(MaxLiquidationOrdersPerBlock).(uint32); ok {
		result.MaxLiquidationOrdersPerBlock = v
	}

	return result
}
