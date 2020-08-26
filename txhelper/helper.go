package txhelper

import (
	"encoding/hex"
	"math"
	"strings"

	"decred.org/dcrwallet/wallet"
	"github.com/decred/dcrd/dcrutil/v3"
	"github.com/decred/dcrd/wire"
)

// MsgTxFromHex returns a wire.MsgTx struct built from the transaction hex string.
func MsgTxFromHex(txhex string) (*wire.MsgTx, error) {
	msgTx := wire.NewMsgTx()
	if err := msgTx.Deserialize(hex.NewDecoder(strings.NewReader(txhex))); err != nil {
		return nil, err
	}
	return msgTx, nil
}

// FeeRate computes the fee rate in atoms/kB for a transaction provided the
// total amount of the transaction's inputs, the total amount of the
// transaction's outputs, and the size of the transaction in bytes. Note that a
// kB refers to 1000 bytes, not a kiB. If the size is 0, the returned fee is -1.
func FeeRate(amtIn, amtOut, sizeBytes int64) int64 {
	if sizeBytes == 0 {
		return -1
	}
	return 1000 * (amtIn - amtOut) / sizeBytes
}

func MsgTxFeeSizeRate(transactionHex string) (msgTx *wire.MsgTx, fee dcrutil.Amount, size int, feeRate dcrutil.Amount, err error) {
	msgTx, err = MsgTxFromHex(transactionHex)
	if err != nil {
		return
	}

	size = msgTx.SerializeSize()
	var amtIn int64
	for iv := range msgTx.TxIn {
		amtIn += msgTx.TxIn[iv].ValueIn
	}
	var amtOut int64
	for iv := range msgTx.TxOut {
		amtOut += msgTx.TxOut[iv].Value
	}
	txSize := int64(msgTx.SerializeSize())
	fee = dcrutil.Amount(amtIn - amtOut)
	feeRate = dcrutil.Amount(FeeRate(amtIn, amtOut, txSize))
	return
}

func TransactionAmountAndDirection(inputTotal, outputTotal, fee int64) (amount int64, direction int32) {
	amountDifference := outputTotal - inputTotal

	if amountDifference < 0 && float64(fee) == math.Abs(float64(amountDifference)) {
		// transferred internally, the only real amount spent was transaction fee
		direction = TxDirectionTransferred
		amount = fee
	} else if amountDifference > 0 {
		// received
		direction = TxDirectionReceived
		amount = outputTotal
	} else {
		// sent
		direction = TxDirectionSent
		amount = inputTotal - outputTotal - fee
	}

	return
}

func FormatTransactionType(txType wallet.TransactionType) string {
	switch txType {
	case wallet.TransactionTypeCoinbase:
		return TxTypeCoinBase
	case wallet.TransactionTypeTicketPurchase:
		return TxTypeTicketPurchase
	case wallet.TransactionTypeVote:
		return TxTypeVote
	case wallet.TransactionTypeRevocation:
		return TxTypeRevocation
	default:
		return TxTypeRegular
	}
}
