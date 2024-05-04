package consensus

import (
	"fmt"
	"testing"
	"time"

	"github.com/cometbft/cometbft/abci/example/kvstore"
	"github.com/cometbft/cometbft/internal/test"
	"github.com/cometbft/cometbft/types"
	sm "github.com/cometbft/cometbft/state"
	"github.com/cometbft/cometbft/crypto"
)

func TestDecideGenesis1(t *testing.T) {
	// Most simple test:
	// All bookmarks point to the same sequence of blocks
	// Send votes from only slightly more than 1/2 of the validators
	// Ensure that the block at height 1 is selected to be committed
	// List of commits can be empty for now

	// Parameters:
	N := 7 // Number of validators
	F := 3 // Number of Byzantine validators (in this case, their bookmarks will just be missing)
	var newGenHeight int64 = 3
	
	logger := consensusLogger().With("test", "decide_genesis")

	state, cs, validators := createStateAndValidators(N) 
	validatorSet := types.NewValidatorSet(validators)
	
	cs.Height = 2
	logger.Debug("Consensus height at initialization", "height", cs.Height)

	somePubKey := validators[0].PubKey

	// create a block with empty or random data
	blocks := makeBlockSequence(newGenHeight, state, somePubKey)

	// create bookmark messages
	fullBookmarks := make([]*types.FullBookmark, N - F)
	for i := F; i < N; i++ {
		valIndex, _ := validatorSet.GetByAddress(validators[i].Address)
		fullBookmarks[i - F] = &types.FullBookmark{
			Height : newGenHeight,
			AllBlocks: blocks,
			AllBlockCommits: make([]*types.Commit, newGenHeight),
			ValidatorAddress: validators[i].Address,
			ValidatorIndex: valIndex,
		}
	}

	blockHashes := make([]string, len(blocks))
	for i := 0; i < len(blocks); i++ {
		blockHashes[i] = fmt.Sprintf("%X", blocks[i].Hash())
	}
	bookmarkHeights := make([]int64, len(fullBookmarks))
	for i := 0; i < len(fullBookmarks); i++ {
		bookmarkHeights[i] = fullBookmarks[i].Height
	}
	logger.Debug("Created blocks and bookmarks", "blocks", blockHashes, "bookmark_heights", bookmarkHeights)

	// Find out the new blocks that need to be committed
	ready, blocksToCommit, _ := cs.decideBlocksToCommit(fullBookmarks, validatorSet)
	if !ready {
		t.Fatalf("Blocks are not ready to be committed")
		return
	}
	logger.Debug("Decided blocks to commit", "num_blocks", len(blocksToCommit))
	for _, block := range blocksToCommit {
		// logger.Debug("Block to be committed", "height", block.Height, "hash", block.Hash())
		t.Logf("Block to be committed: height %d, hash %X", block.Height, block.Hash())
	}
}

func TestDecideGenesis2(t *testing.T) {
	// Honest validators have bookmarks at different heights
	// Adversarial bookmarks are missing
	N := 4
	F := 0
	var newGenHeight int64 = 3

	logger := consensusLogger().With("test", "decide_genesis")

	state, cs, validators := createStateAndValidators(N) 
	validatorSet := types.NewValidatorSet(validators)
	
	logger.Debug("Consensus height at initialization", "height", cs.Height)

	somePubKey := validators[0].PubKey

	// create a block with empty or random data
	blocks := makeBlockSequence(newGenHeight, state, somePubKey)

	// create bookmark messages
	fullBookmarks := make([]*types.FullBookmark, N - F)
	for i := F; i < N; i++ {
		valIndex, _ := validatorSet.GetByAddress(validators[i].Address)
		var bHeight int64
		if i <= N / 2 {
			bHeight = newGenHeight
		} else {
			bHeight = newGenHeight - 1
		}
		fullBookmarks[i - F] = &types.FullBookmark{
			Height : bHeight,
			AllBlocks: blocks[:bHeight],
			AllBlockCommits: make([]*types.Commit, bHeight),
			ValidatorAddress: validators[i].Address,
			ValidatorIndex: valIndex,
		}
	}

	blockHashes := make([]string, len(blocks))
	for i := 0; i < len(blocks); i++ {
		blockHashes[i] = fmt.Sprintf("%X", blocks[i].Hash())
	}
	bookmarkHeights := make([]int64, len(fullBookmarks))
	for i := 0; i < len(fullBookmarks); i++ {
		bookmarkHeights[i] = fullBookmarks[i].Height
	}
	logger.Debug("Created blocks and bookmarks", "blocks", blockHashes, "bookmark_heights", bookmarkHeights)


	// Find out the new blocks that need to be committed
	ready, blocksToCommit, _ := cs.decideBlocksToCommit(fullBookmarks, validatorSet)
	if !ready {
		t.Fatalf("Blocks are not ready to be committed")
		return
	}
	logger.Debug("Decided blocks to commit", "num_blocks", len(blocksToCommit))
	for _, block := range blocksToCommit {
		// logger.Debug("Block to be committed", "height", block.Height, "hash", block.Hash())
		t.Logf("Block to be committed: height %d, hash %X", block.Height, block.Hash())
	}
}

//////////////////////////////////////////////////

func createStateAndValidators(N int) (*sm.State, *State, []*types.Validator) {
	// Create N validators
	// We need a validator set and consensus state object
	// Copying some code from common_test.go
	c := test.ConsensusParams()
	state, privVals := randGenesisState(N, false, 10, c)
	cs := newState(state, privVals[0], kvstore.NewInMemoryApplication()) // I think it doesn't matter which privVal we use here
	validators := make([]*types.Validator, N)
	for i := 0; i < N; i++ {
		pubKey, _ := privVals[i].GetPubKey()
		validators[i] = types.NewValidator(pubKey, 10)
	}
	return &state, cs, validators
}

func makeBlockSequence(newGenHeight int64, state *sm.State, somePubKey crypto.PubKey) []*types.Block {
	blocks := make([]*types.Block, newGenHeight)
	for i := 0; int64(i) < newGenHeight; i++ {
		blocks[i] = types.MakeBlock(int64(i+1), []types.Tx{types.Tx(fmt.Sprintf("Hello World %d", i+1))}, &types.Commit{
			Height: int64(i),
			Round: 0,
			BlockID: types.BlockID{},
			Signatures: make([]types.CommitSig, 0),
		}, nil)
		blocks[i].Header.Populate(
			state.Version.Consensus, state.ChainID,
			time.Now(), state.LastBlockID,
			state.Validators.Hash(), state.NextValidators.Hash(),
			state.ConsensusParams.Hash(), state.AppHash, state.LastResultsHash,
			somePubKey.Address(),
		)
	}
	return blocks
}