package types

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/rlp"
)

type IncrementalMerkle struct {
	NodeCount   uint64       `json:"node_count"`
	ActiveNodes []rlp.Sha256 `json:"active_nodes"`
}

/**
 * return the current root of the incremental merkle
 *
 * @return
 */
func (m IncrementalMerkle) GetRoot() rlp.Sha256 {
	if m.NodeCount > 0 {
		return m.ActiveNodes[len(m.ActiveNodes)-1]
	} else {
		return rlp.Sha256{}
	}
}

/**
 * Add a node to the incremental tree and recalculate the _active_nodes so they
 * are prepared for the next append.
 *
 * The algorithm for this is to start at the new node and retreat through the tree
 * for any node that is the concatenation of a fully-realized node and a partially
 * realized node we must record the value of the fully-realized node in the new
 * _active_nodes so that the next append can fetch it.   Fully realized nodes and
 * Fully implied nodes do not have an effect on the _active_nodes.
 *
 * For convention _AND_ to allow appends when the _node_count is a power-of-2, the
 * current root of the incremental tree is always appended to the end of the new
 * _active_nodes.
 *
 * In practice, this can be done iteratively by recording any "left" value that
 * is to be combined with an implied node.
 *
 * If the appended node is a "left" node in its pair, it will immediately push itself
 * into the new active nodes list.
 *
 * If the new node is a "right" node it will begin collapsing upward in the tree,
 * reading and discarding the "left" node data from the old active nodes list, until
 * it becomes a "left" node.  It must then push the "top" of its current collapsed
 * sub-tree into the new active nodes list.
 *
 * Once any value has been added to the new active nodes, all remaining "left" nodes
 * should be present in the order they are needed in the previous active nodes as an
 * artifact of the previous append.  As they are read from the old active nodes, they
 * will need to be copied in to the new active nodes list as they are still needed
 * for future appends.
 *
 * As a result, if an append collapses the entire tree while always being the "right"
 * node, the new list of active nodes will be empty and by definition the tree contains
 * a power-of-2 number of nodes.
 *
 * Regardless of the contents of the new active nodes list, the top "collapsed" value
 * is appended.  If this tree is _not_ a power-of-2 number of nodes, this node will
 * not be used in the next append but still serves as a conventional place to access
 * the root of the current tree. If this _is_ a power-of-2 number of nodes, this node
 * will be needed during then collapse phase of the next append so, it serves double
 * duty as a legitimate active node and the conventional storage location of the root.
 *
 *
 * @param digest - the node to add
 * @return - the new root
 */
func (m *IncrementalMerkle) Append(digest rlp.Sha256) rlp.Sha256 {
	partial := false
	maxDepth := calculateMaxDepth(m.NodeCount + 1)
	currentDepth := maxDepth - 1
	index := m.NodeCount
	top := digest
	activeIter := 0
	updateActiveNodes := make([]rlp.Sha256, 0, maxDepth)

	for currentDepth > 0 {
		if index&0x1 == 0 {
			// we are collapsing from a "left" value and an implied "right" creating a partial node

			// we only need to append this node if it is fully-realized and by definition
			// if we have encountered a partial node during collapse this cannot be
			// fully-realized
			if !partial {
				updateActiveNodes = append(updateActiveNodes, top)
			}

			// calculate the partially realized node value by implying the "right" value is identical
			// to the "left" value
			top = rlp.Hash256(makeCanonicalPair(top, top))
			partial = true
		} else {
			// we are collapsing from a "right" value and an fully-realized "left"

			// pull a "left" value from the previous active nodes
			leftValue := m.ActiveNodes[activeIter]
			activeIter++

			// if the "right" value is a partial node we will need to copy the "left" as future appends still need it
			// otherwise, it can be dropped from the set of active nodes as we are collapsing a fully-realized node
			if partial {
				updateActiveNodes = append(updateActiveNodes, leftValue)
			}

			// calculate the node
			top = rlp.Hash256(makeCanonicalPair(leftValue, top))
		}

		// move up a level in the tree
		currentDepth--
		index = index >> 1
	}

	// append the top of the collapsed tree (aka the root of the merkle)
	updateActiveNodes = append(updateActiveNodes, top)

	// store the new active_nodes
	moveNodes(&m.ActiveNodes, &updateActiveNodes)

	// update the node count
	m.NodeCount++

	return m.ActiveNodes[len(m.ActiveNodes)-1]
}

func Merkle(ids []rlp.Sha256) rlp.Sha256 {
	if 0 == len(ids) {
		return rlp.Sha256{}
	}

	for len(ids) > 1 {
		if len(ids)%2 > 0 {
			ids = append(ids, ids[len(ids)-1])
		}

		for i := 0; i < len(ids)/2; i++ {
			ids[i] = rlp.Hash256(makeCanonicalPair(ids[2*i], ids[(2*i)+1]))
		}

		ids = ids[:len(ids)/2]
	}

	return ids[0]
}

func makeCanonicalLeft(val rlp.Sha256) rlp.Sha256 {
	canonicalL := val
	canonicalL.Hash[0] &= 0xFFFFFFFFFFFFFF7F
	return canonicalL
}
func makeCanonicalRight(val rlp.Sha256) rlp.Sha256 {
	canonicalR := val
	canonicalR.Hash[0] &= 0x0000000000000080
	return canonicalR
}

func makeCanonicalPair(l rlp.Sha256, r rlp.Sha256) common.Pair {
	return common.Pair{makeCanonicalLeft(l), makeCanonicalRight(r)}
}

func isCanonicalLeft(val rlp.Sha256) bool {
	return val.Hash[0]&0x0000000000000080 == 0
}

func isCanonicalRight(val rlp.Sha256) bool {
	return val.Hash[0]&0x0000000000000080 != 0
}

/**
 * Given a number of nodes return the depth required to store them
 * in a fully balanced binary tree.
 *
 * @param node_count - the number of nodes in the implied tree
 * @return the max depth of the minimal tree that stores them
 */
func calculateMaxDepth(nodeCount uint64) int {
	if nodeCount == 0 {
		return 0
	}
	impliedCount := nextPowerOf2(nodeCount)
	return clzPower2(impliedCount) + 1
}

/**
 * given an unsigned integral number return the smallest
 * power-of-2 which is greater than or equal to the given number
 *
 * @param value - an unsigned integral
 * @return - the minimum power-of-2 which is >= value
 */
func nextPowerOf2(value uint64) uint64 {
	value -= 1
	value |= value >> 1
	value |= value >> 2
	value |= value >> 4
	value |= value >> 8
	value |= value >> 16
	value |= value >> 32
	value += 1
	return value
}

/**
 * Given a power-of-2 (assumed correct) return the number of leading zeros
 *
 * This is a classic count-leading-zeros in parallel without the necessary
 * math to make it safe for anything that is not already a power-of-2
 *
 * @param value - and integral power-of-2
 * @return the number of leading zeros
 */
func clzPower2(value uint64) int {
	lz := 64

	if value > 0 {
		lz--
	}
	if value&0x00000000FFFFFFFF > 0 {
		lz -= 32
	}
	if value&0x0000FFFF0000FFFF > 0 {
		lz -= 16
	}
	if value&0x00FF00FF00FF00FF > 0 {
		lz -= 8
	}
	if value&0x0F0F0F0F0F0F0F0F > 0 {
		lz -= 4
	}
	if value&0x3333333333333333 > 0 {
		lz -= 2
	}
	if value&0x5555555555555555 > 0 {
		lz -= 1
	}

	return lz
}

func moveNodes(to *[]rlp.Sha256, from *[]rlp.Sha256) {
	*to = make([]rlp.Sha256, len(*from))
	copy(*to, *from)
}
