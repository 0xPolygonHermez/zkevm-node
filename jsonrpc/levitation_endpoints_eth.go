package jsonrpc

import (
	"context"

	"net/http"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
)

// LEVITATION_BEGIN

// SendRawTransaction has two different ways to handle new transactions:
// - for Sequencer nodes it tries to add the tx to the pool
// - for Non-Sequencer nodes it relays the Tx to the Sequencer node
func (e *EthEndpoints) SendRawTransaction(httpRequest *http.Request, input string) (interface{}, types.Error) {
	if e.cfg.SequencerNodeURI != "" {
		return e.relayTxToSequencerNode(input)
	} else {
		ip := ""
		ips := httpRequest.Header.Get("X-Forwarded-For")

		// TODO: this is temporary patch remove this log
		realIp := httpRequest.Header.Get("X-Real-IP")
		log.Infof("X-Forwarded-For: %s, X-Real-IP: %s", ips, realIp)

		if ips != "" {
			ip = strings.Split(ips, ",")[0]
		}

		hash, err := e.tryVerifyTx(input, ip)
		_ = hash
		if err != nil {
			return nil, err
		}

		return e.SendRawTransactionFromSequencer(httpRequest, input)
	}
}

func (e *EthEndpoints) SendRawTransactionFromSequencer(httpRequest *http.Request, input string) (interface{}, types.Error) {
	if e.cfg.SequencerNodeURI != "" {
		return e.relayTxToSequencerNode(input)
	} else {
		ip := ""
		ips := httpRequest.Header.Get("X-Forwarded-For")

		// TODO: this is temporary patch remove this log
		realIp := httpRequest.Header.Get("X-Real-IP")
		log.Infof("X-Forwarded-For: %s, X-Real-IP: %s", ips, realIp)

		if ips != "" {
			ip = strings.Split(ips, ",")[0]
		}
		return e.tryToAddTxToPool(input, ip)
	}
}

func (e *EthEndpoints) tryVerifyTx(input, ip string) (interface{}, types.Error) {
	tx, err := hexToTx(input)
	if err != nil {
		return RPCErrorResponse(types.InvalidParamsErrorCode, "invalid tx input", err)
	}

	log.Infof("verifying TX: %v", tx.Hash().Hex())
	if err := e.pool.VerifyTx(context.Background(), *tx, ip); err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, err.Error(), nil)
	}
	log.Infof("TX verified: %v", tx.Hash().Hex())

	return tx.Hash().Hex(), nil
}
