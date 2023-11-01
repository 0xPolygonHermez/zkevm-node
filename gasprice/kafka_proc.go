package gasprice

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	kafka "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

const (
	okbcoinId      = 7184
	defaultTime    = 10
	defaultMaxData = 10e6 // 10M
)

// MsgInfo msg info
type MsgInfo struct {
	Topic string `json:"topic"`
	Data  *Body  `json:"data"`
}

// Body msg body
type Body struct {
	Id        string   `json:"id"`
	PriceList []*Price `json:"priceList"`
}

// Price coin price
type Price struct {
	CoinId                   int     `json:"coinId"`
	Symbol                   string  `json:"symbol"`
	FullName                 string  `json:"fullName"`
	Price                    float64 `json:"price"`
	PriceStatus              int     `json:"priceStatus"`
	MaxPrice24H              float64 `json:"maxPrice24H"`
	MinPrice24H              float64 `json:"minPrice24H"`
	MarketCap                float64 `json:"marketCap"`
	Timestamp                int64   `json:"timestamp"`
	Vol24H                   float64 `json:"vol24h"`
	CirculatingSupply        float64 `json:"circulatingSupply"`
	MaxSupply                float64 `json:"maxSupply"`
	TotalSupply              float64 `json:"totalSupply"`
	PriceChange24H           float64 `json:"priceChange24H"`
	PriceChangeRate24H       float64 `json:"priceChangeRate24H"`
	CirculatingMarketCap     float64 `json:"circulatingMarketCap"`
	PriceChange7D            float64 `json:"priceChange7D"`
	PriceChangeRate7D        float64 `json:"priceChangeRate7D"`
	PriceChange30D           float64 `json:"priceChange30D"`
	PriceChangeRate30D       float64 `json:"priceChangeRate30D"`
	PriceChangeYearStart     float64 `json:"priceChangeYearStart"`
	PriceChangeRateYearStart float64 `json:"priceChangeRateYearStart"`
	ExceptionStatus          int     `json:"exceptionStatus"`
	Source                   int     `json:"source"`
	Type                     string  `json:"type"`
	Id                       string  `json:"id"`
}

// KafkaProcessor kafka processor
type KafkaProcessor struct {
	kreader  *kafka.Reader
	L2Price  float64
	ctx      context.Context
	rwLock   sync.RWMutex
	l2CoinId int
}

func newKafkaProcessor(cfg Config, ctx context.Context) *KafkaProcessor {
	rp := &KafkaProcessor{
		kreader:  getKafkaReader(cfg),
		L2Price:  cfg.DefaultL2CoinPrice,
		ctx:      ctx,
		l2CoinId: okbcoinId,
	}
	if cfg.L2CoinId != 0 {
		rp.l2CoinId = cfg.L2CoinId
	}

	go rp.processor()
	return rp
}

func getKafkaReader(cfg Config) *kafka.Reader {
	brokers := strings.Split(cfg.KafkaURL, ",")

	var dialer *kafka.Dialer
	if cfg.Password != "" && cfg.Username != "" && cfg.RootCAPath != "" {
		rootCA, err := os.ReadFile(cfg.RootCAPath)
		if err != nil {
			panic("kafka read root ca fail")
		}
		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(rootCA); !ok {
			panic("caCertPool.AppendCertsFromPEM")
		}
		dialer = &kafka.Dialer{
			Timeout:       defaultTime * time.Second,
			DualStack:     true,
			SASLMechanism: plain.Mechanism{Username: cfg.Username, Password: cfg.Password},
		}
		{ // #nosec G402
			dialer.TLS = &tls.Config{RootCAs: caCertPool, InsecureSkipVerify: true}
		}
	}

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     cfg.GroupID,
		Topic:       cfg.Topic,
		MinBytes:    1, // 1
		MaxBytes:    defaultMaxData,
		Dialer:      dialer,
		StartOffset: kafka.LastOffset, // read data from new message
	})
}

func (rp *KafkaProcessor) processor() {
	log.Info("kafka processor start processor ")
	defer rp.kreader.Close()
	for {
		select {
		case <-rp.ctx.Done():
			return
		default:
			value, err := rp.ReadAndCalc(rp.ctx)
			if err != nil {
				log.Warn("get the destion data fail ", err)
				time.Sleep(time.Second * defaultTime)
				continue
			}
			rp.updateL2CoinPrice(value)
		}
	}
}

// ReadAndCalc read and calc
func (rp *KafkaProcessor) ReadAndCalc(ctx context.Context) (float64, error) {
	m, err := rp.kreader.ReadMessage(ctx)
	if err != nil {
		return 0, err
	}
	return rp.parseL2CoinPrice(m.Value)
}

func (rp *KafkaProcessor) updateL2CoinPrice(price float64) {
	rp.rwLock.Lock()
	defer rp.rwLock.Unlock()
	rp.L2Price = price
}

// GetL2CoinPrice get L2 coin price
func (rp *KafkaProcessor) GetL2CoinPrice() float64 {
	rp.rwLock.RLock()
	defer rp.rwLock.RUnlock()
	return rp.L2Price
}

func (rp *KafkaProcessor) parseL2CoinPrice(value []byte) (float64, error) {
	msgI := &MsgInfo{}
	err := json.Unmarshal(value, &msgI)
	if err != nil {
		return 0, err
	}
	if msgI.Data == nil || len(msgI.Data.PriceList) == 0 {
		return 0, fmt.Errorf("the data PriceList is empty")
	}
	for _, price := range msgI.Data.PriceList {
		if price.CoinId == rp.l2CoinId {
			return price.Price, nil
		}
	}
	return 0, fmt.Errorf("not find a correct coin price coinId=%v", rp.l2CoinId)
}
