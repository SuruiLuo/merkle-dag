package merkledag

import "hash"

// HashPool 是一个哈希对象池的接口，规定了从池中获取哈希对象的方法
type HashPool interface {
	Get() hash.Hash // Get 方法用于获取一个哈希对象
}
