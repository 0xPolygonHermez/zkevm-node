package ethtxmanager

import (
	"context"
	"math/big"
	"testing"
	"time"

	zktypes "github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
)

const (
	domain                 = "http://asset-onchain.base-defi.svc.test.local:7001"
	seqAddr                = "1a13bddcc02d363366e04d4aa588d3c125b0ff6f"
	aggAddr                = "66e39a1e507af777e8c385e2d91559e20e306303"
	contractAddr           = "8947dc90862f386968966b22dfe5edf96435bc2f"
	contractAddrAgg        = "1d5298ee11f7cd56fb842b7894346bfb2e47a95f"
	l1ChainID       uint64 = 11155111
	AccessKey              = ""
	SecretKey              = ""
	// domain       = "http://127.0.0.1:7001"
	// seqAddr      = "f39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	// aggAddr      = "70997970C51812dc3A010C7d01b50e0d17dc79C8"
	// contractAddr = "812cB73e48841a6736bB94c65c56341817cE6304"
)

func TestClientPostSignRequestAndWaitResultSeqFork8(t *testing.T) {
	client := &Client{
		etherman: mockEtherman{},
		cfg: Config{
			CustodialAssets: CustodialAssetsConfig{
				Enable:            false,
				URL:               domain,
				Symbol:            2882,
				SequencerAddr:     common.HexToAddress(seqAddr),
				AggregatorAddr:    common.HexToAddress(aggAddr),
				WaitResultTimeout: zktypes.NewDuration(4 * time.Minute),
				OperateTypeSeq:    3,
				OperateTypeAgg:    4,
				ProjectSymbol:     3011,
				OperateSymbol:     2,
				SysFrom:           3,
				UserID:            0,
				OperateAmount:     0,
				RequestSignURI:    "/priapi/v1/assetonchain/ecology/ecologyOperate",
				QuerySignURI:      "/priapi/v1/assetonchain/ecology/querySignDataByOrderNo",
				AccessKey:         AccessKey,
				SecretKey:         SecretKey,
			},
		},
	}
	ctx := context.WithValue(context.Background(), traceID, uuid.New().String())
	txInput, _ := hex.DecodeHex("0xdb5b0ed700000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000065f7d95600000000000000000000000000000000000000000000000000000000000500ae0000000000000000000000005e7b89ab3b2de21f0f35da4920b9d7310ccbe25900000000000000000000000000000000000000000000000000000000000003400000000000000000000000000000000000000000000000000000000000000005e4c5b9632f3e9d42c0ee4150e8ff8e1abb75dfa09b136cde6eaf45bd4ff5e9d500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000038fa368fd50abe3a986733aa13036b89aa4f5a0ef0a4b411ec1ce940a9646b0800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000086e4978ee17260754d1cdcc92b1fabb0f5169e289e215a764cb5b3bfab9eeded000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000032108d0820b532d572a962eb69ec35ff3bf1f4149f3f175282c63ea818b1907000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000db851df016929d0c4a75bb942a171d0d95bf8c107098a0292f34031a60a38e6b00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000691991e3774b9aabe2594ac3eeea8b010a82e0d6147bf6fd5fb3fe6ea731a17f1643d91b4e1925f88b8f85eb4dd01b03df735c08e0811b217dc6c340eb08535eef1bb1a3f77364476988d978bc67b8d0db8743325048fb8587d5a4b4b3c31f3c9844af08d16ec07bcb880000000000000000000000000000000000000000000000")

	tx := types.NewTransaction(0, common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"), big.NewInt(10), 50000, big.NewInt(10), txInput)

	seqReq, _ := client.unpackSequenceBatchesTx(tx)
	to := common.HexToAddress(contractAddr)
	mTx := monitoredTx{
		from:      common.HexToAddress(seqAddr),
		to:        &to,
		gasPrice:  big.NewInt(12345678912),
		gas:       2000000,
		gasOffset: 100,
	}
	ret, _ := seqReq.marshal(common.HexToAddress(contractAddr), mTx)

	req := client.newSignRequest(client.cfg.CustodialAssets.OperateTypeSeq, client.cfg.CustodialAssets.SequencerAddr, ret)

	_, err := client.postSignRequestAndWaitResult(ctx, mTx, req)
	if err != nil {
		t.Log(err)
	}
}

func TestClientPostSignRequestAndWaitResultAggFork8(t *testing.T) {
	client := &Client{
		etherman: mockEtherman{},
		cfg: Config{
			CustodialAssets: CustodialAssetsConfig{
				Enable:            false,
				URL:               domain,
				Symbol:            2882,
				SequencerAddr:     common.HexToAddress(seqAddr),
				AggregatorAddr:    common.HexToAddress(aggAddr),
				WaitResultTimeout: zktypes.NewDuration(4 * time.Minute),
				OperateTypeSeq:    3,
				OperateTypeAgg:    4,
				ProjectSymbol:     3011,
				OperateSymbol:     2,
				SysFrom:           3,
				UserID:            0,
				OperateAmount:     0,
				RequestSignURI:    "/priapi/v1/assetonchain/ecology/ecologyOperate",
				QuerySignURI:      "/priapi/v1/assetonchain/ecology/querySignDataByOrderNo",
				AccessKey:         AccessKey,
				SecretKey:         SecretKey,
			},
		},
	}
	ctx := context.WithValue(context.Background(), traceID, uuid.New().String())
	txInput, _ := hex.DecodeHex("0x1489ed100000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000500ae00000000000000000000000000000000000000000000000000000000000500b33d2b755195521747e8b32d1f28b033f0d93082872cd253ee42e3bcfaf2cd8a3129e24287e603ff699b4322778e511f2bf9f4d305f199403c24db686e44c5b1cf000000000000000000000000a57d7641fba20d916a7cb61de5af71e9743e4b571666fc7a9b894634feba222cebcf0678a56b2b6348e31af2eddea7b98fc9b36a017778ab8f75848be64ffeaf6974eb09369e01a96c6c443d611289e3b18ebcae02fa881020960c0cd3fc9819024a60e1869baead40db9422af24947c1ec825d31f3643021dd597b4b6799b1669e50d5ccc29dcab571308b4afd9f69856b6f1cf1560e3fdedf10b3b9b0fa9ebb9a09af91ebd61c3722274eef718948620bcc84a22ba517411823319b4fad2bac5ca0d30f72261ca3a291fe11ee3f2d65f21947c2cda3b57c03ea3d576cef736b45185e6a4f7e2fae3d2e1442245ac6cf906032d032cb2ce822d9e16521b862e46fa28bc62ce12f0804fa057e8ecef7bff12d72423095607eb51333715cdc3cd39d316e02ef9d92ce24ec75518353670b87b48a7178bb86cb37a5c2b08b8705977b3d046a8f7777dbcb988f6a8508d9ebe3c3ce42f33985c054be78209809df36aa39ebecc34c44bbdf116ba4226613b5f2562b41139c17ccf95c2791327defc83fc7e611bc8e08aa72a9bca9138652f8ca9742f25eb11f74ffd3737e147acb9dd4cbffc298860cfab785c4e9b5cceb92aafe344262278e991e0d715bd9a66aafe8153b723b9ca7ac050205970cfc85a8dca1cad2a22a2b4642d0217277893d12d1ede382ca42eed04b686a9229eb3fed45d838b27474a3242ee89ded32881b4abc11fd072d7b718af15171364c7f5b285f5da8d303245456ec4e9beb1c2930cb7db35b4ab991a7ea48a2373a4c2ef201ba4151c15e12c35f10e806ab5fd807fb6cf20255ca28e54d5952dba961b510af29296c7240e73818dfef35e5919998ce4885fc9671c0379f63479f2d61197725916b1a814b67a71ce84f0148fd0205ff1bc3e20df8a0e4a1180cecc57cdaa4f690c07bb2010b9e0c3d25daa958dbdbaad240e5900a97e8d0e75197f983b7870188775471344ab8ee51610daf6f634de0407c9d24ca762a38a6f2ba44d4114749c616d972bfbfd7cb939ababbf9e6d1d569b9a3ff396ab3475b3bfe3332b3ac218c2f6e516780ca06e7ec5a98c542686c92c2c2ff09cff58101df7824b01aa6b0b85e206")

	tx := types.NewTransaction(0, common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"), big.NewInt(10), 50000, big.NewInt(10), txInput)

	seqReq, _ := client.unpackVerifyBatchesTrustedAggregatorTx(tx)

	to := common.HexToAddress(contractAddrAgg)
	mTx := monitoredTx{
		from:      common.HexToAddress(aggAddr),
		to:        &to,
		nonce:     0,
		gas:       2000000,
		gasOffset: 200,
		gasPrice:  big.NewInt(12),
	}
	ret, _ := seqReq.marshal(common.HexToAddress(contractAddrAgg), mTx)

	req := client.newSignRequest(client.cfg.CustodialAssets.OperateTypeAgg, client.cfg.CustodialAssets.AggregatorAddr, ret)

	_, err := client.postSignRequestAndWaitResult(ctx, mTx, req)
	if err != nil {
		t.Log(err)
	}
}

type mockEtherman struct{}

func (m mockEtherman) GetTx(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockEtherman) GetTxReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockEtherman) WaitTxToBeMined(ctx context.Context, tx *types.Transaction, timeout time.Duration) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockEtherman) SendTx(ctx context.Context, tx *types.Transaction) error {
	//TODO implement me
	panic("implement me")
}

func (m mockEtherman) CurrentNonce(ctx context.Context, account common.Address) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockEtherman) SuggestedGasPrice(ctx context.Context) (*big.Int, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockEtherman) EstimateGas(ctx context.Context, from common.Address, to *common.Address, value *big.Int, data []byte) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockEtherman) CheckTxWasMined(ctx context.Context, txHash common.Hash) (bool, *types.Receipt, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockEtherman) SignTx(ctx context.Context, sender common.Address, tx *types.Transaction) (*types.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockEtherman) GetRevertMessage(ctx context.Context, tx *types.Transaction) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockEtherman) GetZkEVMAddressAndL1ChainID() (common.Address, common.Address, uint64, error) {
	return common.HexToAddress(contractAddr), common.HexToAddress(contractAddrAgg), l1ChainID, nil
}
