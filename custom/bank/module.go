package bank

import (
	"cosmossdk.io/core/address"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankmodule "github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/bank/types"

	custombankkeeper "github.com/terra-money/alliance/custom/bank/keeper"
)

// AppModule wraps around the bank module and the bank keeper to return the right total supply ignoring bonded tokens
// that the alliance module minted to rebalance the voting power
// It modifies the TotalSupply and SupplyOf GRPC queries
type AppModule struct {
	bankmodule.AppModule
	keeper       custombankkeeper.Keeper
	addressCodec address.Codec
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper custombankkeeper.Keeper, accountKeeper types.AccountKeeper) AppModule {
	bankModule := bankmodule.NewAppModule(cdc, keeper, accountKeeper)
	return AppModule{
		AppModule:    bankModule,
		keeper:       keeper,
		addressCodec: accountKeeper.AddressCodec(),
	}
}

// RegisterServices registers module services.
// NOTE: Overriding this method as not doing so will cause a panic
// when trying to force this custom keeper into a bankkeeper.BaseKeeper
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), custombankkeeper.NewMsgServerImpl(am.keeper, am.addressCodec))
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}

// RegisterStoreDecoder registers a decoder for supply module's types
func (am AppModule) RegisterStoreDecoder(sdr simtypes.StoreDecoderRegistry) {
	sdr[types.StoreKey] = simtypes.NewStoreDecoderFuncFromCollectionsSchema(am.keeper.BaseKeeper.Schema)
}
