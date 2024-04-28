package consensus

import (
	"testing"
	"context"
	"sync"
	"time"

	"github.com/cometbft/cometbft/p2p"
	"github.com/stretchr/testify/require"
	"github.com/cometbft/cometbft/types"
	cfg "github.com/cometbft/cometbft/config"
	cmtcons "github.com/cometbft/cometbft/proto/tendermint/consensus"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
)

func TestSafetyViolation(t *testing.T) {
	// boilerplate for generating a testnet
	N := 7
	F := 5 // number of Byzantine validators
	EnableFreezing := true
	logger := consensusLogger().With("test", "safety_violation")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := newKVStore
	css, cleanup := randConsensusNet(t, N, "safety_violation_test", newMockTickerFunc(false), app, func(c *cfg.Config) {
		c.Consensus.EnableFreezingGadget = EnableFreezing
		c.Consensus.WaitBeforeCommit = 2 * time.Second
	})
	defer cleanup()

	// give the byzantine validator(s) a normal ticker
	for i := 0; i < F; i++ {
		ticker := NewTimeoutTicker()
		ticker.SetLogger(css[i].Logger)
		css[i].SetTimeoutTicker(ticker)
	}

	// For now, copy following part from byzantine_test. Maybe reuse startConsensusNet later

	switches := make([]*p2p.Switch, N)
	p2pLogger := logger.With("module", "p2p")
	for i := 0; i < N; i++ {
		switches[i] = p2p.MakeSwitch(
			config.P2P,
			i,
			func(i int, sw *p2p.Switch) *p2p.Switch {
				return sw
			})
		switches[i].SetLogger(p2pLogger.With("validator", i))
	}

	// switches[5] and switches[6] correspond to the honest validators. 

	blocksSubs := make([]types.Subscription, N)
	reactors := make([]p2p.Reactor, N)

	// Define variables to store the two conflicting blocks that the first byzantine validator will propose. These will be required for the other Byzantine validators to vote upon.
	var byzBlockHash1 []byte
	var byzBlockHash2 []byte
	var byzBlockPartsHeader1 types.PartSetHeader
	var byzBlockPartsHeader2 types.PartSetHeader
	var byzBlockLock = new(sync.Mutex)

	byzantineDecideProposalFunc2 := func(ctx context.Context, t *testing.T, height int64, round int32, cs *State, sw *p2p.Switch) {
		// byzantine user should create two proposals and try to split the vote.
		// Avoid sending on internalMsgQueue and running consensus state.
	
		// Create a new proposal block from state/txs from the mempool.
		block1, err := cs.createProposalBlock(ctx)
		require.NoError(t, err)
		blockParts1, err := block1.MakePartSet(types.BlockPartSizeBytes)
		require.NoError(t, err)
		polRound, propBlockID := cs.ValidRound, types.BlockID{Hash: block1.Hash(), PartSetHeader: blockParts1.Header()}
		proposal1 := types.NewProposal(height, round, polRound, propBlockID)
		p1 := proposal1.ToProto()
		if err := cs.privValidator.SignProposal(cs.state.ChainID, p1); err != nil {
			t.Error(err)
		}
	
		proposal1.Signature = p1.Signature
	
		// some new transactions come in (this ensures that the proposals are different)
		deliverTxsRange(t, cs, 0, 1)
	
		// Create a new proposal block from state/txs from the mempool.
		block2, err := cs.createProposalBlock(ctx)
		require.NoError(t, err)
		blockParts2, err := block2.MakePartSet(types.BlockPartSizeBytes)
		require.NoError(t, err)
		polRound, propBlockID = cs.ValidRound, types.BlockID{Hash: block2.Hash(), PartSetHeader: blockParts2.Header()}
		proposal2 := types.NewProposal(height, round, polRound, propBlockID)
		p2 := proposal2.ToProto()
		if err := cs.privValidator.SignProposal(cs.state.ChainID, p2); err != nil {
			t.Error(err)
		}
	
		proposal2.Signature = p2.Signature
	
		byzBlockLock.Lock()
		byzBlockHash1 = block1.Hash()
		byzBlockHash2 = block2.Hash()
		byzBlockPartsHeader1 = blockParts1.Header()
		byzBlockPartsHeader2 = blockParts2.Header()
		byzBlockLock.Unlock()
	
		// broadcast conflicting proposals/block parts to peers
		peers := sw.Peers().List()
		t.Logf("Byzantine: broadcasting conflicting proposals to %d peers", len(peers))
		t.Logf("Block hash 1: %X", byzBlockHash1)
		t.Logf("Block hash 2: %X", byzBlockHash2)
		for _, peer := range peers {
			// Send one block to peer #5 and another block to peer #6.
			// Send both blocks to all other peers.
			if getSwitchIndex(switches, peer) != N - 1 { // send to everyone except 6
				go sendProposalAndParts3(height, round, peer, proposal1, blockParts1)
			}
			if getSwitchIndex(switches, peer) != N - 2 { // send to everyone except 5
				go sendProposalAndParts3(height, round, peer, proposal2, blockParts2)
			}
		}
	}

	byzantineDoPrevoteFunc2 := func(cs *State, sw *p2p.Switch, valIndex int) {
		// broadcast prevotes and precommits for the conflicting proposals
		// to all peers
		byzBlockLock.Lock()
		blockHash1 := byzBlockHash1
		blockHash2 := byzBlockHash2
		blockPartsHeader1 := byzBlockPartsHeader1
		blockPartsHeader2 := byzBlockPartsHeader2
		byzBlockLock.Unlock()
		peers := sw.Peers().List()
		t.Logf("Byzantine validator %d: broadcasting conflicting votes to %d peers", valIndex, len(peers))
		for _, peer := range peers {
			// Send one block to peer #5 and another block to peer #6.
			// Send both blocks to all other peers.
			if getSwitchIndex(switches, peer) != N - 1 {
				go sendVotes3(blockHash1, blockPartsHeader1, cs, peer)
			}
			if getSwitchIndex(switches, peer) != N - 2 {
				go sendVotes3(blockHash2, blockPartsHeader2, cs, peer)
			}
		}
	}

	for i := 0; i < N; i++ {

		// enable txs so we can create different proposals
		assertMempool(css[i].txNotifier).EnableTxsAvailable()
		//make first F val(s) byzantine but on ly the first one will propose conflicting blocks
		if i < F {
			// NOTE: Now, test validators are MockPV, which by default doesn't
			// do any safety checks.
			// Therefore, DisableChecks() does nothing.
			css[i].privValidator.(types.MockPV).DisableChecks()
			j := i
			// Set the doPrevote function to send prevote and precommits for both blocks proposed by byzantine validator #0.
			css[i].doPrevote = func(height int64, round int32) {
				byzantineDoPrevoteFunc2(css[j], switches[j], j)
			}
		}
		if i == 0 {
			// Override the decideProposal function for the byzantine validator
			j := i
			css[i].decideProposal = func(height int64, round int32) {
				// The decideProposal function in byzantine_test sends one proposal to a few nodes and another proposal to the rest. Here, we will send conflicting proposals to all nodes.
				byzantineDecideProposalFunc2(ctx, t, height, round, css[j], switches[j])
			}
		}

		eventBus := css[i].eventBus
		eventBus.SetLogger(logger.With("module", "events", "validator", i))

		var err error
		blocksSubs[i], err = eventBus.Subscribe(context.Background(), testSubscriber, types.EventQueryNewBlock)
		require.NoError(t, err)

		conR := NewReactor(css[i], true) // so we don't start the consensus states
		conR.SetLogger(logger.With("validator", i))
		conR.SetEventBus(eventBus)

		var conRI p2p.Reactor = conR

		// make first F vals byzantine
		// ByzantineReactor only changes the AddPeer function of Reactor.
		// Right now, I don't understand what the difference is and why we need it.
		if i < F {
			conRI = NewByzantineReactor(conR)
		}

		reactors[i] = conRI
		err = css[i].blockExec.Store().Save(css[i].state) // for save height 1's validators info
		require.NoError(t, err)
	}

	defer func() {
		for _, r := range reactors {
			if rr, ok := r.(*ByzantineReactor); ok {
				err := rr.reactor.Switch.Stop()
				require.NoError(t, err)
			} else {
				err := r.(*Reactor).Switch.Stop()
				require.NoError(t, err)
			}
		}
	}()

	p2p.MakeConnectedSwitches(config.P2P, N, func(i int, s *p2p.Switch) *p2p.Switch {
		// ignore new switch s, we already made ours
		switches[i].AddReactor("CONSENSUS", reactors[i])
		return switches[i]
	}, func(sws []*p2p.Switch, i, j int) {
		// In byzantine_test, there was a network partition. Here, we will connect all nodes.
		// Note that our simulations for Better Safe than Sorry must be in the synchronous network model.
		p2p.Connect2Switches(sws, i, j)
	})

	// start the non-byz state machines.
	// note these must be started before the byz
	for i := F; i < N; i++ {
		cr := reactors[i].(*Reactor)
		cr.SwitchToConsensus(cr.conS.GetState(), false)
	}

	// start the byzantine state machines
	for i := 0; i < F; i++ {
		byzR := reactors[i].(*ByzantineReactor)
		s := byzR.reactor.conS.GetState()
		byzR.reactor.SwitchToConsensus(s, false)
	}

	// Instead of waiting for all validators to commit the first block, we will just wait for a certain amount of time and then observe which blocks have been committed by all validators.
	<-time.After(10 * time.Second)
}

func sendProposalAndParts3(
	height int64,
	round int32,
	peer p2p.Peer,
	proposal *types.Proposal,
	parts *types.PartSet,
) {
	// proposal
	peer.Send(p2p.Envelope{
		ChannelID: DataChannel,
		Message:   &cmtcons.Proposal{Proposal: *proposal.ToProto()},
	})

	// parts
	for i := 0; i < int(parts.Total()); i++ {
		part := parts.GetPart(i)
		pp, err := part.ToProto()
		if err != nil {
			panic(err) // TODO: wbanfield better error handling
		}
		peer.Send(p2p.Envelope{
			ChannelID: DataChannel,
			Message: &cmtcons.BlockPart{
				Height: height, // This tells peer that this part applies to us.
				Round:  round,  // This tells peer that this part applies to us.
				Part:   *pp,
			},
		})
	}
}

func sendVotes3(
	blockHash []byte,
	partsHeader types.PartSetHeader,
	cs *State,
	peer p2p.Peer,
) {
	// votes
	cs.mtx.Lock()
	prevote, _ := cs.signVote(cmtproto.PrevoteType, blockHash, partsHeader, nil)
	precommit, _ := cs.signVote(cmtproto.PrecommitType, blockHash, partsHeader, nil)
	cs.mtx.Unlock()
	peer.Send(p2p.Envelope{
		ChannelID: VoteChannel,
		Message:   &cmtcons.Vote{Vote: prevote.ToProto()},
	})
	peer.Send(p2p.Envelope{
		ChannelID: VoteChannel,
		Message:   &cmtcons.Vote{Vote: precommit.ToProto()},
	})
}