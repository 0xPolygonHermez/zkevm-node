package jsonrpc

import (
	"context"
	"math/big"
	"strconv"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
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

func argBytesPtr(b []byte) *argBytes {
	bb := argBytes(b)

	return &bb
}

type argBig big.Int

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

	return []byte("0x" + b.Text(hex.Base)), nil
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
	From     common.Address
	To       *common.Address
	Gas      *argUint64
	GasPrice *argUint64
	Value    *argBytes
	Data     *argBytes
}

// ToTransaction transforms txnArgs into a Transaction
func (args *txnArgs) ToUnsignedTransaction(ctx context.Context, st stateInterface, blockNumber uint64, cfg Config, dbTx pgx.Tx) (common.Address, *types.Transaction, error) {
	gas := cfg.MaxCumulativeGasUsed
	if args.Gas != nil && uint64(*args.Gas) > uint64(0) {
		gas = uint64(*args.Gas)
	}

	value := big.NewInt(0)
	if args.Value != nil {
		value.SetBytes(*args.Value)
	}

	data := []byte{}
	if args.Data != nil {
		data = *args.Data
	}

	sender := args.From
	nonce := uint64(0)
	gasPrice := big.NewInt(0)

	defaultSenderAddress := common.HexToAddress(cfg.DefaultSenderAddress)
	if sender == state.ZeroAddress {
		sender = defaultSenderAddress
	}

	if sender != defaultSenderAddress {
		if args.GasPrice != nil {
			gasPrice.SetUint64(uint64(*args.GasPrice))
		}

		n, err := st.GetNonce(ctx, sender, blockNumber, dbTx)
		if err != nil {
			return common.Address{}, nil, err
		}
		nonce = uint64(n)
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       args.To,
		Value:    value,
		Gas:      gas,
		GasPrice: gasPrice,
		Data:     data,
	})

	return sender, tx, nil
}

type rpcBlock struct {
	ParentHash      common.Hash            `json:"parentHash"`
	Sha3Uncles      common.Hash            `json:"sha3Uncles"`
	Miner           common.Address         `json:"miner"`
	StateRoot       common.Hash            `json:"stateRoot"`
	TxRoot          common.Hash            `json:"transactionsRoot"`
	ReceiptsRoot    common.Hash            `json:"receiptsRoot"`
	LogsBloom       types.Bloom            `json:"logsBloom"`
	Difficulty      argUint64              `json:"difficulty"`
	TotalDifficulty argUint64              `json:"totalDifficulty"`
	Size            argUint64              `json:"size"`
	Number          argUint64              `json:"number"`
	GasLimit        argUint64              `json:"gasLimit"`
	GasUsed         argUint64              `json:"gasUsed"`
	Timestamp       argUint64              `json:"timestamp"`
	ExtraData       argBytes               `json:"extraData"`
	MixHash         common.Hash            `json:"mixHash"`
	Nonce           argBytes               `json:"nonce"`
	Hash            common.Hash            `json:"hash"`
	Transactions    []rpcTransactionOrHash `json:"transactions"`
	Uncles          []common.Hash          `json:"uncles"`
}

func l2BlockToRPCBlock(b *types.Block, fullTx bool) *rpcBlock {
	h := b.Header()

	n := big.NewInt(0).SetUint64(h.Nonce.Uint64())
	nonce := common.LeftPadBytes(n.Bytes(), 8) //nolint:gomnd

	var difficulty uint64
	if h.Difficulty != nil {
		difficulty = h.Difficulty.Uint64()
	} else {
		difficulty = uint64(0)
	}

	res := &rpcBlock{
		ParentHash:      h.ParentHash,
		Sha3Uncles:      h.UncleHash,
		Miner:           h.Coinbase,
		StateRoot:       h.Root,
		TxRoot:          h.TxHash,
		ReceiptsRoot:    h.ReceiptHash,
		LogsBloom:       h.Bloom,
		Difficulty:      argUint64(difficulty),
		TotalDifficulty: argUint64(difficulty),
		Size:            argUint64(b.Size()),
		Number:          argUint64(b.Number().Uint64()),
		GasLimit:        argUint64(h.GasLimit),
		GasUsed:         argUint64(h.GasUsed),
		Timestamp:       argUint64(h.Time),
		ExtraData:       argBytes(h.Extra),
		MixHash:         h.MixDigest,
		Nonce:           nonce,
		Hash:            b.Hash(),
		Transactions:    []rpcTransactionOrHash{},
		Uncles:          []common.Hash{},
	}

	for idx, txn := range b.Transactions() {
		if fullTx {
			blockHash := b.Hash()
			txIndex := uint64(idx)
			tx := toRPCTransaction(txn, b.Number(), &blockHash, &txIndex)
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
type rpcTransactionOrHash interface {
	getHash() common.Hash
}

type rpcTransaction struct {
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

func (t rpcTransaction) getHash() common.Hash { return t.Hash }

// Redefine to implement getHash() of transactionOrHash
type transactionHash common.Hash

func (h transactionHash) getHash() common.Hash { return common.Hash(h) }

func (h transactionHash) MarshalText() ([]byte, error) {
	return []byte(common.Hash(h).String()), nil
}

func toRPCTransaction(
	t *types.Transaction,
	blockNumber *big.Int,
	blockHash *common.Hash,
	txIndex *uint64,
) *rpcTransaction {
	v, r, s := t.RawSignatureValues()

	from, _ := state.GetSender(*t)

	res := &rpcTransaction{
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
		bn := argUint64(blockNumber.Uint64())
		res.BlockNumber = &bn
	}

	res.BlockHash = blockHash

	if txIndex != nil {
		ti := argUint64(*txIndex)
		res.TxIndex = &ti
	}

	return res
}

type rpcReceipt struct {
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
	FromAddr          common.Address  `json:"from"`
	ToAddr            *common.Address `json:"to"`
	ContractAddress   *common.Address `json:"contractAddress"`
	Type              argUint64       `json:"type"`
}

func receiptToRPCReceipt(tx types.Transaction, r *types.Receipt) (rpcReceipt, error) {
	to := tx.To()
	logs := r.Logs
	if logs == nil {
		logs = []*types.Log{}
	}

	var contractAddress *common.Address
	if r.ContractAddress != state.ZeroAddress {
		ca := r.ContractAddress
		contractAddress = &ca
	}

	blockNumber := argUint64(0)
	if r.BlockNumber != nil {
		blockNumber = argUint64(r.BlockNumber.Uint64())
	}

	from, err := state.GetSender(tx)
	if err != nil {
		return rpcReceipt{}, err
	}

	return rpcReceipt{
		Root:              common.BytesToHash(r.PostState),
		CumulativeGasUsed: argUint64(r.CumulativeGasUsed),
		LogsBloom:         r.Bloom,
		Logs:              logs,
		Status:            argUint64(r.Status),
		TxHash:            r.TxHash,
		TxIndex:           argUint64(r.TransactionIndex),
		BlockHash:         r.BlockHash,
		BlockNumber:       blockNumber,
		GasUsed:           argUint64(r.GasUsed),
		ContractAddress:   contractAddress,
		FromAddr:          from,
		ToAddr:            to,
		Type:              argUint64(r.Type),
	}, nil
}

type rpcLog struct {
	Address     common.Address `json:"address"`
	Topics      []common.Hash  `json:"topics"`
	Data        argBytes       `json:"data"`
	BlockNumber argUint64      `json:"blockNumber"`
	TxHash      common.Hash    `json:"transactionHash"`
	TxIndex     argUint64      `json:"transactionIndex"`
	BlockHash   common.Hash    `json:"blockHash"`
	LogIndex    argUint64      `json:"logIndex"`
	Removed     bool           `json:"removed"`
}

func logToRPCLog(l types.Log) rpcLog {
	return rpcLog{
		Address:     l.Address,
		Topics:      l.Topics,
		Data:        l.Data,
		BlockNumber: argUint64(l.BlockNumber),
		TxHash:      l.TxHash,
		TxIndex:     argUint64(l.TxIndex),
		BlockHash:   l.BlockHash,
		LogIndex:    argUint64(l.Index),
		Removed:     l.Removed,
	}
}
