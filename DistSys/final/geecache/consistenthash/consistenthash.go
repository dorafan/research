package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps bytes to uint32
type Hash func(data []byte) uint32

// Map constains all hashed keys
type Map struct {
	hash     Hash
	//Virtual node multiple
	replicas int
	//hash ring
	keys     []int // Sorted
	//The virtual node and the real node mapping table hashMap,
	//the key is the hash value of the virtual node,
	//and the value is the name of the real node.
	hashMap  map[int]string
}

// New creates a Map instance
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add adds some keys to the hash.
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		//For each real node key, m.replicas virtual nodes are created correspondingly,
		for i := 0; i < m.replicas; i++ {
			//the name of the virtual node is: strconv.Itoa(i) + key, that is,
			//different virtual nodes are distinguished by adding numbers.
			//Use m.hash() to calculate the hash value of the virtual node,
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			//and use append(m.keys, hash) to add to the ring.
			m.keys = append(m.keys, hash)
			//Add the mapping relationship between virtual nodes and real nodes in hashMap.
			m.hashMap[hash] = key
		}
	}
	//The last step is to sort the hash values on the ring.
	sort.Ints(m.keys)
}

// Get gets the closest item in the hash to the provided key.
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	//Calculate the hash value of key.
	hash := int(m.hash([]byte(key)))
	// Binary search for appropriate replica.
	//Find the index idx of the first matching virtual node clockwise,
	//and obtain the corresponding hash value from m.keys
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	//Get the real node through hashMap mapping
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
