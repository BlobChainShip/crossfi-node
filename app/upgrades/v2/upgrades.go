package v2

import (
	"errors"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	evmkeeper "github.com/evmos/evmos/v12/x/evm/keeper"
	evmtypes "github.com/evmos/evmos/v12/x/evm/types"
	feemarketkeeper "github.com/evmos/evmos/v12/x/feemarket/keeper"
	feemarkettypes "github.com/evmos/evmos/v12/x/feemarket/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	ek evmkeeper.Keeper,
	fk feemarketkeeper.Keeper,
	pk paramskeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)

		if err := ek.SetParams(ctx, evmtypes.DefaultParams()); err != nil {
			return nil, err
		}

		if err := fk.SetParams(ctx, feemarkettypes.DefaultParams()); err != nil {
			return nil, err
		}

		paramsSubspace, found := pk.GetSubspace(baseapp.Paramspace)
		if !found {
			return nil, errors.New("params subspace not found")
		}

		currentBlockParams := &abci.BlockParams{}
		paramsSubspace.Get(ctx, baseapp.ParamStoreKeyBlockParams, currentBlockParams)
		currentBlockParams.MaxGas = 20000000
		paramsSubspace.Set(ctx, baseapp.ParamStoreKeyBlockParams, currentBlockParams)

		logger.Debug("running module migrations ...")
		return mm.RunMigrations(ctx, configurator, vm)
	}
}