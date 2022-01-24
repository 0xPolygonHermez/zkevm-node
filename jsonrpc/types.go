package jsonrpc

import (
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/state/helper"
)

type argUint64 uint64

// MarshalText marshals into text
func (b argUint64) MarshalText() ([]byte, error) {
	buf := make([]byte, 2, encoding.Base10) //nolint:gomnd
	copy(buf, `0x`)
	buf = strconv.AppendUint(buf, uint64(b), hex.Base)
	return buf, nil
}

// UnmarshalText unmarshals from text
func (b *argUint64) UnmarshalText(input []byte) error {
	str := strings.TrimPrefix(string(input), "0x")
	num, err := strconv.ParseUint(str, hex.Base, encoding.BitSize64)
	if err != nil {
		return err
	}
	*b = argUint64(num)
	return nil
}

type argBytes []byte

// MarshalText marshals into text
func (b argBytes) MarshalText() ([]byte, error) {
	return encodeToHex(b), nil
}

// UnmarshalText unmarshals from text
func (b *argBytes) UnmarshalText(input []byte) error {
	hh, err := decodeToHex(input)
	if err != nil {
		return nil
	}
	aux := make([]byte, len(hh))
	copy(aux[:], hh[:])
	*b = aux
	return nil
}

type argBig big.Int

func argBigPtr(b *big.Int) *argBig {
	v := argBig(*b)

	return &v
}

func (a *argBig) UnmarshalText(input []byte) error {
	buf, err := decodeToHex(input)
	if err != nil {
		return err
	}

	b := new(big.Int)
	b.SetBytes(buf)
	*a = argBig(*b)

	return nil
}

func (a argBig) MarshalText() ([]byte, error) {
	b := (*big.Int)(&a)

	return []byte("0x" + b.Text(16)), nil
}

func decodeToHex(b []byte) ([]byte, error) {
	str := string(b)
	str = strings.TrimPrefix(str, "0x")
	if len(str)%2 != 0 {
		str = "0" + str
	}
	return hex.DecodeString(str)
}

func encodeToHex(b []byte) []byte {
	str := hex.EncodeToString(b)
	if len(str)%2 != 0 {
		str = "0" + str
	}
	return []byte("0x" + str)
}

// txnArgs is the transaction argument for the rpc endpoints
type txnArgs struct {
	From     *common.Address
	To       *common.Address
	Gas      *argUint64
	GasPrice *argBytes
	Value    *argBytes
	Input    *argBytes
	Data     *argBytes
	Nonce    *argUint64
}

// ToTransaction transforms txnArgs into a Transaction
func (arg *txnArgs) ToTransaction() *types.Transaction {
	nonce := uint64(0)
	if arg.Nonce != nil {
		nonce = uint64(*arg.Nonce)
	}

	gas := uint64(0)
	if arg.Gas != nil {
		gas = uint64(*arg.Gas)
	}

	gasPrice := hex.DecodeHexToBig(string(*arg.GasPrice))

	value := big.NewInt(0)
	if arg.Value != nil {
		value = hex.DecodeHexToBig(string(*arg.Value))
	}

	data := []byte{}
	if arg.Data != nil {
		data = *arg.Data
	}

	tx := types.NewTransaction(nonce, *arg.To, value, gas, gasPrice, data)

	return tx
}

type block struct {
	ParentHash      common.Hash         `json:"parentHash"`
	Sha3Uncles      common.Hash         `json:"sha3Uncles"`
	Miner           common.Address      `json:"miner"`
	StateRoot       common.Hash         `json:"stateRoot"`
	TxRoot          common.Hash         `json:"transactionsRoot"`
	ReceiptsRoot    common.Hash         `json:"receiptsRoot"`
	LogsBloom       types.Bloom         `json:"logsBloom"`
	Difficulty      argUint64           `json:"difficulty"`
	TotalDifficulty argUint64           `json:"totalDifficulty"`
	Size            argUint64           `json:"size"`
	Number          argUint64           `json:"number"`
	GasLimit        argUint64           `json:"gasLimit"`
	GasUsed         argUint64           `json:"gasUsed"`
	Timestamp       argUint64           `json:"timestamp"`
	ExtraData       argBytes            `json:"extraData"`
	MixHash         common.Hash         `json:"mixHash"`
	Nonce           uint64              `json:"nonce"`
	Hash            common.Hash         `json:"hash"`
	Transactions    []transactionOrHash `json:"transactions"`
	Uncles          []common.Hash       `json:"uncles"`
}

func toBlock(b *types.Block, fullTx bool) *block {
	h := b.Header()
	res := &block{
		ParentHash:      h.ParentHash,
		Sha3Uncles:      h.UncleHash,
		Miner:           h.Coinbase,
		StateRoot:       h.Root,
		TxRoot:          h.TxHash,
		ReceiptsRoot:    h.ReceiptHash,
		LogsBloom:       h.Bloom,
		Difficulty:      argUint64(h.Difficulty.Uint64()),
		TotalDifficulty: argUint64(h.Difficulty.Uint64()), // not needed for POS
		Size:            argUint64(b.Size()),
		Number:          argUint64(h.Number.Uint64()),
		GasLimit:        argUint64(h.GasLimit),
		GasUsed:         argUint64(h.GasUsed),
		Timestamp:       argUint64(h.Time),
		ExtraData:       argBytes(h.Extra),
		MixHash:         h.MixDigest,
		Nonce:           h.Nonce.Uint64(),
		Hash:            h.Hash(),
		Transactions:    []transactionOrHash{},
		Uncles:          []common.Hash{},
	}

	for idx, txn := range b.Transactions() {
		if fullTx {
			number := argUint64(b.Number().Uint64())
			hash := b.Hash()

			tx := toTransaction(
				txn,
				&number,
				&hash,
				&idx,
			)
			res.Transactions = append(
				res.Transactions,
				tx,
			)
		} else {
			res.Transactions = append(
				res.Transactions,
				transactionHash(txn.Hash()),
			)
		}
	}

	for _, uncle := range b.Uncles() {
		res.Uncles = append(res.Uncles, uncle.Hash())
	}

	return res
}

// For union type of transaction and types.Hash
type transactionOrHash interface {
	getHash() common.Hash
}

type transaction struct {
	Nonce       argUint64       `json:"nonce"`
	GasPrice    argBig          `json:"gasPrice"`
	Gas         argUint64       `json:"gas"`
	To          *common.Address `json:"to"`
	Value       argBig          `json:"value"`
	Input       argBytes        `json:"input"`
	V           argBig          `json:"v"`
	R           argBig          `json:"r"`
	S           argBig          `json:"s"`
	Hash        common.Hash     `json:"hash"`
	From        common.Address  `json:"from"`
	BlockHash   *common.Hash    `json:"blockHash"`
	BlockNumber *argUint64      `json:"blockNumber"`
	TxIndex     *argUint64      `json:"transactionIndex"`
}

func (t transaction) getHash() common.Hash { return t.Hash }

// Redefine to implement getHash() of transactionOrHash
type transactionHash common.Hash

func (h transactionHash) getHash() common.Hash { return common.Hash(h) }

func (h transactionHash) MarshalText() ([]byte, error) {
	return []byte(common.Hash(h).String()), nil
}

func toTransaction(
	t *types.Transaction,
	blockNumber *argUint64,
	blockHash *common.Hash,
	txIndex *int,
) *transaction {
	v, r, s := t.RawSignatureValues()

	from, _ := helper.GetSender(t)

	res := &transaction{
		Nonce:    argUint64(t.Nonce()),
		GasPrice: argBig(*t.GasPrice()),
		Gas:      argUint64(t.Gas()),
		To:       t.To(),
		Value:    argBig(*t.Value()),
		Input:    t.Data(),
		V:        argBig(*v),
		R:        argBig(*r),
		S:        argBig(*s),
		Hash:     t.Hash(),
		From:     from,
	}

	if blockNumber != nil {
		res.BlockNumber = blockNumber
	}

	if blockHash != nil {
		res.BlockHash = blockHash
	}

	if txIndex != nil {
		i := argUint64(uint64(*txIndex))
		res.TxIndex = &i
	}

	return res
}

type receipt struct {
	Root              common.Hash     `json:"root"`
	CumulativeGasUsed argUint64       `json:"cumulativeGasUsed"`
	LogsBloom         types.Bloom     `json:"logsBloom"`
	Logs              []*types.Log    `json:"logs"`
	Status            argUint64       `json:"status"`
	TxHash            common.Hash     `json:"transactionHash"`
	TxIndex           argUint64       `json:"transactionIndex"`
	BlockHash         common.Hash     `json:"blockHash"`
	BlockNumber       argUint64       `json:"blockNumber"`
	GasUsed           argUint64       `json:"gasUsed"`
	ContractAddress   common.Address  `json:"contractAddress"`
	FromAddr          common.Address  `json:"from"`
	ToAddr            *common.Address `json:"to"`
}
