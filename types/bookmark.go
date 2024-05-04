package types

import (
	"bytes"
)

type FullBookmark struct {
	Height int64
	AllBlocks []*Block
	AllBlockCommits []*Commit
	ValidatorAddress Address
	ValidatorIndex int32
}

func (b1 *FullBookmark) IsPrefixOf(b2 *FullBookmark) bool {
	if len(b1.AllBlocks) > len(b2.AllBlocks) {
		return false
	}
	for i, block := range b1.AllBlocks {
		if !bytes.Equal(block.Hash(), b2.AllBlocks[i].Hash()) {
			return false
		}
	}
	return true
}