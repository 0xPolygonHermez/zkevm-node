package helper

import (
	"hash"
	"sync"

	"github.com/hermeznetwork/hermez-core/log"
	"github.com/umbracle/fastrlp"
	"golang.org/x/crypto/sha3"
)

// Keccak256 calculates the Keccak256
func Keccak256(v ...[]byte) []byte {
	h := sha3.NewLegacyKeccak256()
	for _, i := range v {
		_, err := h.Write(i)
		if err != nil {
			log.Error("error while calculating the keccak-256")
		}
	}
	return h.Sum(nil)
}

// Keccak256To hashes a src with keccak-256
func Keccak256To(dst, src []byte) []byte {
	h := DefaultKeccakPool.Get()
	_, err := h.Write(src)
	if err != nil {
		log.Error("error while hashing with keccak-256")
	}
	dst = h.Sum(dst)
	DefaultKeccakPool.Put(h)
	return dst
}

// DefaultKeccakPool is a default pool
var DefaultKeccakPool Pool

// Pool is a pool of keccaks
type Pool struct {
	pool sync.Pool
}

// Get returns the keccak
func (p *Pool) Get() *Keccak {
	v := p.pool.Get()
	if v == nil {
		return NewKeccak256()
	}
	return v.(*Keccak)
}

// Put releases the keccak
func (p *Pool) Put(k *Keccak) {
	k.Reset()
	p.pool.Put(k)
}

type hashImpl interface {
	hash.Hash
	Read(b []byte) (int, error)
}

// Keccak is the sha256 keccak hash
type Keccak struct {
	buf  []byte // buffer to store intermediate rlp marshal values
	tmp  []byte
	hash hashImpl
}

// WriteRlp writes an RLP value
func (k *Keccak) WriteRlp(dst []byte, v *fastrlp.Value) []byte {
	k.buf = v.MarshalTo(k.buf[:0])
	_, err := k.Write(k.buf)
	if err != nil {
		log.Error("error while writing RLP value")
	}
	return k.Sum(dst)
}

// Write implements the hash interface
func (k *Keccak) Write(b []byte) (int, error) {
	return k.hash.Write(b)
}

// Reset implements the hash interface
func (k *Keccak) Reset() {
	k.buf = k.buf[:0]
	k.hash.Reset()
}

// Read hashes the content and returns the intermediate buffer.
func (k *Keccak) Read() []byte {
	_, err := k.hash.Read(k.tmp)
	if err != nil {
		log.Error("error while reading the hash")
	}
	return k.tmp
}

// Sum implements the hash interface
func (k *Keccak) Sum(dst []byte) []byte {
	_, err := k.hash.Read(k.tmp)
	if err != nil {
		log.Error("error while reading Sum")
	}
	dst = append(dst, k.tmp[:]...)
	return dst
}

func newKeccak(hash hashImpl) *Keccak {
	return &Keccak{
		hash: hash,
		tmp:  make([]byte, hash.Size()),
	}
}

// NewKeccak256 returns a new keccak 256
func NewKeccak256() *Keccak {
	return newKeccak(sha3.NewLegacyKeccak256().(hashImpl))
}

// NewKeccak512 returns a new keccak 512
func NewKeccak512() *Keccak {
	return newKeccak(sha3.NewLegacyKeccak512().(hashImpl))
}
