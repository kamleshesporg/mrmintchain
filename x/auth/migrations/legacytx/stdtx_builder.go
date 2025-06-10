package legacytx

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
)

// StdTxBuilder wraps StdTx to implement client.TxBuilder interface.
// This is maintained for backwards compatibility with amino StdTx.
type StdTxBuilder struct {
	StdTx
	cdc *codec.LegacyAmino
}

var _ client.TxBuilder = &StdTxBuilder{}

func (s *StdTxBuilder) GetTx() authsigning.Tx {
	return s.StdTx
}

func (s *StdTxBuilder) SetMsgs(msgs ...sdk.Msg) error {
	s.Msgs = msgs
	return nil
}

func (s *StdTxBuilder) SetSignatures(signatures ...signing.SignatureV2) error {
	sigs := make([]StdSignature, len(signatures))
	var err error
	for i, sig := range signatures {
		sigs[i], err = SignatureV2ToStdSignature(s.cdc, sig)
		if err != nil {
			return err
		}
	}
	s.Signatures = sigs
	return nil
}

func (s *StdTxBuilder) SetFeeAmount(amount sdk.Coins) {
	s.Fee.Amount = amount
}

func (s *StdTxBuilder) SetGasLimit(limit uint64) {
	s.Fee.Gas = limit
}

func (s *StdTxBuilder) SetTip(_ *tx.Tip) {
	panic("StdTxBuilder does not support tips")
}

func (s *StdTxBuilder) SetMemo(memo string) {
	s.Memo = memo
}

func (s *StdTxBuilder) SetTimeoutHeight(height uint64) {
	s.TimeoutHeight = height
}

func (s *StdTxBuilder) SetFeeGranter(_ sdk.AccAddress) {}

func (s *StdTxBuilder) SetFeePayer(_ sdk.AccAddress) {}

func (s *StdTxBuilder) AddAuxSignerData(_ tx.AuxSignerData) error {
	return sdkerrors.ErrLogic.Wrap("cannot use AuxSignerData with StdTxBuilder")
}

// StdTxConfig implements client.TxConfig for StdTx (amino).
type StdTxConfig struct {
	Cdc *codec.LegacyAmino
}

var _ client.TxConfig = StdTxConfig{}

func (s StdTxConfig) NewTxBuilder() client.TxBuilder {
	return &StdTxBuilder{
		StdTx: StdTx{},
		cdc:   s.Cdc,
	}
}

func (s StdTxConfig) WrapTxBuilder(tx sdk.Tx) (client.TxBuilder, error) {
	stdTx, ok := tx.(StdTx)
	if !ok {
		return nil, fmt.Errorf("wrong type, expected %T, got %T", stdTx, tx)
	}
	return &StdTxBuilder{StdTx: stdTx, cdc: s.Cdc}, nil
}

func (s StdTxConfig) TxEncoder() sdk.TxEncoder {
	return DefaultTxEncoder(s.Cdc)
}

func (s StdTxConfig) TxDecoder() sdk.TxDecoder {
	return mkDecoder(s.Cdc.Unmarshal)
}

func (s StdTxConfig) TxJSONEncoder() sdk.TxEncoder {
	return func(tx sdk.Tx) ([]byte, error) {
		return s.Cdc.MarshalJSON(tx)
	}
}

func (s StdTxConfig) TxJSONDecoder() sdk.TxDecoder {
	return mkDecoder(s.Cdc.UnmarshalJSON)
}

func (s StdTxConfig) MarshalSignatureJSON(sigs []signing.SignatureV2) ([]byte, error) {
	stdSigs := make([]StdSignature, len(sigs))
	for i, sig := range sigs {
		stdSig, err := SignatureV2ToStdSignature(s.Cdc, sig)
		if err != nil {
			return nil, err
		}
		stdSigs[i] = stdSig
	}
	return s.Cdc.MarshalJSON(stdSigs)
}

func (s StdTxConfig) UnmarshalSignatureJSON(bz []byte) ([]signing.SignatureV2, error) {
	var stdSigs []StdSignature
	err := s.Cdc.UnmarshalJSON(bz, &stdSigs)
	if err != nil {
		return nil, err
	}

	sigs := make([]signing.SignatureV2, len(stdSigs))
	for i, stdSig := range stdSigs {
		sig, err := StdSignatureToSignatureV2(s.Cdc, stdSig)
		if err != nil {
			return nil, err
		}
		sigs[i] = sig
	}
	return sigs, nil
}

func (s StdTxConfig) SignModeHandler() authsigning.SignModeHandler {
	return &stdTxSignModeHandler{} // ✅ Return pointer type
}

// Deprecated: SignatureV2ToStdSignature is legacy amino-only.
func SignatureV2ToStdSignature(cdc *codec.LegacyAmino, sig signing.SignatureV2) (StdSignature, error) {
	var (
		sigBz []byte
		err   error
	)

	if sig.Data != nil {
		sigBz, err = SignatureDataToAminoSignature(cdc, sig.Data)
		if err != nil {
			return StdSignature{}, err
		}
	}

	return StdSignature{
		PubKey:    sig.PubKey,
		Signature: sigBz,
	}, nil
}

// Unmarshaler is a generic type for Unmarshal functions
type Unmarshaler func(bytes []byte, ptr interface{}) error

func mkDecoder(unmarshaler Unmarshaler) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, error) {
		if len(txBytes) == 0 {
			return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "tx bytes are empty")
		}
		var tx StdTx
		err := unmarshaler(txBytes, &tx)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
		}
		return tx, nil
	}
}

func DefaultTxEncoder(cdc *codec.LegacyAmino) sdk.TxEncoder {
	return func(tx sdk.Tx) ([]byte, error) {
		return cdc.Marshal(tx)
	}
}
