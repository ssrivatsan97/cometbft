package types

// Added all the fields that seem useful to verify the bookmark and later finalize the bookmarked block
// This is not optimal because there is repeated communication wherein every validator who bookmarks the same block sends the full block and block parts. But this avoids having to implement an interactive pull-based download of block parts.
// TODO: Skipping BlockParts for now! Will have to send each block part separately while sending the bookmark.
// TODO: Need to modify if more than one block needs to be confirmed during recovery
type Bookmark struct {
	Height int64
	BlockID BlockID
	Block *Block
	ValidatorAddress Address
	ValidatorIndex int32
	Signature []byte
}

func NewUnsignedBookmark(height int64, blockID BlockID, block *Block, validatorAddress Address, validatorIndex int32) *Bookmark {
	return &Bookmark{
		Height: height,
		BlockID: blockID,
		Block: block,
		ValidatorAddress: validatorAddress,
		ValidatorIndex: validatorIndex,
	}
}

// SignBookmark signs the bookmark with the given validator's private key.
// This function modifies the bookmark in place.
// func SignBookmark(bookmark *Bookmark, val PrivValidator, chainID string) (bool, error) {

// }