package gov

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/stake"
	abci "github.com/tendermint/tendermint/abci/types"
	"fmt"
)

func testTickExpiredDepositPeriod(t *testing.T) {
	mapp, keeper, _, addrs, _, _ := getMockApp(t, 10)
	mapp.BeginBlock(abci.RequestBeginBlock{})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})
	govHandler := NewHandler(keeper)

	require.Nil(t, keeper.InactiveProposalQueuePeek(ctx))
	require.False(t, shouldPopInactiveProposalQueue(ctx, keeper))

	newProposalMsg := NewMsgSubmitProposal("Test", "test", ProposalTypeText, addrs[0], sdk.Coins{sdk.NewInt64Coin("steak", 5)})

	res := govHandler(ctx, newProposalMsg)
	require.True(t, res.IsOK())

	EndBlocker(ctx, keeper)
	require.NotNil(t, keeper.InactiveProposalQueuePeek(ctx))
	require.False(t, shouldPopInactiveProposalQueue(ctx, keeper))

	ctx = ctx.WithBlockHeight(10)
	EndBlocker(ctx, keeper)
	require.NotNil(t, keeper.InactiveProposalQueuePeek(ctx))
	require.False(t, shouldPopInactiveProposalQueue(ctx, keeper))

	ctx = ctx.WithBlockHeight(250)
	require.NotNil(t, keeper.InactiveProposalQueuePeek(ctx))
	require.True(t, shouldPopInactiveProposalQueue(ctx, keeper))
	EndBlocker(ctx, keeper)
	require.Nil(t, keeper.InactiveProposalQueuePeek(ctx))
	require.False(t, shouldPopInactiveProposalQueue(ctx, keeper))
}

func testTickMultipleExpiredDepositPeriod(t *testing.T) {
	mapp, keeper, _, addrs, _, _ := getMockApp(t, 10)
	mapp.BeginBlock(abci.RequestBeginBlock{})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})
	govHandler := NewHandler(keeper)

	require.Nil(t, keeper.InactiveProposalQueuePeek(ctx))
	require.False(t, shouldPopInactiveProposalQueue(ctx, keeper))

	newProposalMsg := NewMsgSubmitProposal("Test", "test", ProposalTypeText, addrs[0], sdk.Coins{sdk.NewInt64Coin("steak", 5)})

	res := govHandler(ctx, newProposalMsg)
	require.True(t, res.IsOK())

	EndBlocker(ctx, keeper)
	require.NotNil(t, keeper.InactiveProposalQueuePeek(ctx))
	require.False(t, shouldPopInactiveProposalQueue(ctx, keeper))

	ctx = ctx.WithBlockHeight(10)
	EndBlocker(ctx, keeper)
	require.NotNil(t, keeper.InactiveProposalQueuePeek(ctx))
	require.False(t, shouldPopInactiveProposalQueue(ctx, keeper))

	newProposalMsg2 := NewMsgSubmitProposal("Test2", "test2", ProposalTypeText, addrs[1], sdk.Coins{sdk.NewInt64Coin("steak", 5)})
	res = govHandler(ctx, newProposalMsg2)
	require.True(t, res.IsOK())

	ctx = ctx.WithBlockHeight(205)
	require.NotNil(t, keeper.InactiveProposalQueuePeek(ctx))
	require.True(t, shouldPopInactiveProposalQueue(ctx, keeper))
	EndBlocker(ctx, keeper)
	require.NotNil(t, keeper.InactiveProposalQueuePeek(ctx))
	require.False(t, shouldPopInactiveProposalQueue(ctx, keeper))

	ctx = ctx.WithBlockHeight(215)
	require.NotNil(t, keeper.InactiveProposalQueuePeek(ctx))
	require.True(t, shouldPopInactiveProposalQueue(ctx, keeper))
	EndBlocker(ctx, keeper)
	require.Nil(t, keeper.InactiveProposalQueuePeek(ctx))
	require.False(t, shouldPopInactiveProposalQueue(ctx, keeper))
}

func testTickPassedDepositPeriod(t *testing.T) {
	mapp, keeper, _, addrs, _, _ := getMockApp(t, 10)
	mapp.BeginBlock(abci.RequestBeginBlock{})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})
	govHandler := NewHandler(keeper)

	require.Nil(t, keeper.InactiveProposalQueuePeek(ctx))
	require.False(t, shouldPopInactiveProposalQueue(ctx, keeper))
	require.Nil(t, keeper.ActiveProposalQueuePeek(ctx))
	require.False(t, shouldPopActiveProposalQueue(ctx, keeper))

	newProposalMsg := NewMsgSubmitProposal("Test", "test", ProposalTypeText, addrs[0], sdk.Coins{sdk.NewInt64Coin("steak", 5)})

	res := govHandler(ctx, newProposalMsg)
	require.True(t, res.IsOK())
	var proposalID int64
	keeper.cdc.UnmarshalBinaryBare(res.Data, &proposalID)

	EndBlocker(ctx, keeper)
	require.NotNil(t, keeper.InactiveProposalQueuePeek(ctx))
	require.False(t, shouldPopInactiveProposalQueue(ctx, keeper))

	ctx = ctx.WithBlockHeight(10)
	EndBlocker(ctx, keeper)
	require.NotNil(t, keeper.InactiveProposalQueuePeek(ctx))
	require.False(t, shouldPopInactiveProposalQueue(ctx, keeper))

	newDepositMsg := NewMsgDeposit(addrs[1], proposalID, sdk.Coins{sdk.NewInt64Coin("steak", 5)})
	res = govHandler(ctx, newDepositMsg)
	require.True(t, res.IsOK())

	require.NotNil(t, keeper.InactiveProposalQueuePeek(ctx))
	require.True(t, shouldPopInactiveProposalQueue(ctx, keeper))
	require.NotNil(t, keeper.ActiveProposalQueuePeek(ctx))

	EndBlocker(ctx, keeper)

	require.Nil(t, keeper.InactiveProposalQueuePeek(ctx))
	require.False(t, shouldPopInactiveProposalQueue(ctx, keeper))
	require.NotNil(t, keeper.ActiveProposalQueuePeek(ctx))
	require.False(t, shouldPopActiveProposalQueue(ctx, keeper))

}

func testTickPassedVotingPeriod(t *testing.T) {
	mapp, keeper, _, addrs, _, _ := getMockApp(t, 10)
	SortAddresses(addrs)
	mapp.BeginBlock(abci.RequestBeginBlock{})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})
	govHandler := NewHandler(keeper)

	require.Nil(t, keeper.InactiveProposalQueuePeek(ctx))
	require.False(t, shouldPopInactiveProposalQueue(ctx, keeper))
	require.Nil(t, keeper.ActiveProposalQueuePeek(ctx))
	require.False(t, shouldPopActiveProposalQueue(ctx, keeper))

	newProposalMsg := NewMsgSubmitProposal("Test", "test", ProposalTypeText, addrs[0], sdk.Coins{sdk.NewInt64Coin("steak", 5)})

	res := govHandler(ctx, newProposalMsg)
	require.True(t, res.IsOK())
	var proposalID int64
	keeper.cdc.UnmarshalBinaryBare(res.Data, &proposalID)

	ctx = ctx.WithBlockHeight(10)
	newDepositMsg := NewMsgDeposit(addrs[1], proposalID, sdk.Coins{sdk.NewInt64Coin("steak", 5)})
	res = govHandler(ctx, newDepositMsg)
	require.True(t, res.IsOK())

	EndBlocker(ctx, keeper)

	ctx = ctx.WithBlockHeight(215)
	require.True(t, shouldPopActiveProposalQueue(ctx, keeper))
	depositsIterator := keeper.GetDeposits(ctx, proposalID)
	require.True(t, depositsIterator.Valid())
	depositsIterator.Close()
	require.Equal(t, StatusVotingPeriod, keeper.GetProposal(ctx, proposalID).GetStatus())

	EndBlocker(ctx, keeper)

	require.Nil(t, keeper.ActiveProposalQueuePeek(ctx))
	depositsIterator = keeper.GetDeposits(ctx, proposalID)
	require.False(t, depositsIterator.Valid())
	depositsIterator.Close()
	require.Equal(t, StatusRejected, keeper.GetProposal(ctx, proposalID).GetStatus())
	require.True(t, keeper.GetProposal(ctx, proposalID).GetTallyResult().Equals(EmptyTallyResult()))
}

func testSlashing(t *testing.T) {
	mapp, keeper, sk, addrs, _, _ := getMockApp(t, 10)
	SortAddresses(addrs)
	mapp.BeginBlock(abci.RequestBeginBlock{})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})
	govHandler := NewHandler(keeper)
	stakeHandler := stake.NewHandler(sk)

	createValidators(t, stakeHandler, ctx, addrs[:3], []int64{25, 6, 7})

	initTotalPower := keeper.ds.GetValidatorSet().TotalPower(ctx)
	val0Initial := keeper.ds.GetValidatorSet().Validator(ctx, addrs[0]).GetPower().Quo(initTotalPower)
	val1Initial := keeper.ds.GetValidatorSet().Validator(ctx, addrs[1]).GetPower().Quo(initTotalPower)
	val2Initial := keeper.ds.GetValidatorSet().Validator(ctx, addrs[2]).GetPower().Quo(initTotalPower)
	fmt.Println(keeper.ds.GetValidatorSet().Validator(ctx, addrs[0]).GetPower(),keeper.ds.GetValidatorSet().Validator(ctx, addrs[1]).GetPower(),keeper.ds.GetValidatorSet().Validator(ctx, addrs[2]).GetPower())

	newProposalMsg := NewMsgSubmitProposal("Test", "test", ProposalTypeText, addrs[0], sdk.Coins{sdk.Coin{Denom: "steak", Amount: sdk.NewInt(int64(10)).Mul(Pow10(18))}})

	res := govHandler(ctx, newProposalMsg)
	require.True(t, res.IsOK())
	var proposalID int64
	keeper.cdc.UnmarshalBinaryBare(res.Data, &proposalID)

	ctx = ctx.WithBlockHeight(5)
	require.Equal(t, StatusVotingPeriod, keeper.GetProposal(ctx, proposalID).GetStatus())
    fmt.Println("----",proposalID)
	newVoteMsg := NewMsgVote(addrs[0], proposalID, OptionYes)
	res = govHandler(ctx, newVoteMsg)
	require.True(t, res.IsOK())

	EndBlocker(ctx, keeper)


	ctx = ctx.WithBlockHeight(12)
	require.Equal(t, StatusVotingPeriod, keeper.GetProposal(ctx, proposalID).GetStatus())

	EndBlocker(ctx, keeper)

	fmt.Println(keeper.GetProposal(ctx, proposalID).GetTallyResult(),EmptyTallyResult())
	require.False(t, keeper.GetProposal(ctx, proposalID).GetTallyResult().Equals(EmptyTallyResult()))

	endTotalPower := keeper.ds.GetValidatorSet().TotalPower(ctx)
	val0End := keeper.ds.GetValidatorSet().Validator(ctx, addrs[0]).GetPower().Quo(endTotalPower)
	val1End := keeper.ds.GetValidatorSet().Validator(ctx, addrs[1]).GetPower().Quo(endTotalPower)
	val2End := keeper.ds.GetValidatorSet().Validator(ctx, addrs[2]).GetPower().Quo(endTotalPower)

    fmt.Println(keeper.ds.GetValidatorSet().Validator(ctx, addrs[0]).GetPower(),keeper.ds.GetValidatorSet().Validator(ctx, addrs[1]).GetPower(),keeper.ds.GetValidatorSet().Validator(ctx, addrs[2]).GetPower())
	require.True(t, val0End.GTE(val0Initial))
	require.True(t, val1End.LT(val1Initial))
	require.True(t, val2End.LT(val2Initial))
}
