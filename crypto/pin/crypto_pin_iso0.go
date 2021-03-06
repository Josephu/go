package pin

import (
	"bytes"
	_ "crypto/rand"
	"encoding/hex"
	"fmt"
	_ "log"
	"strconv"
)

type PinBlock_Iso0 struct {
	PinBlocker
}

func (pinBlock *PinBlock_Iso0) Encrypt(pan string, clearPin string, key []byte) ([]byte, error) {

	if len(clearPin) > 12 {
		panic("pin length > 12")
	}

	buf := bytes.NewBufferString(fmt.Sprintf("0%X%s", len(clearPin), clearPin))
	for buf.Len() < 16 {
		buf.WriteString("F")
	}

	pinBlockDataA, err := hex.DecodeString(buf.String())
	if err != nil {
		return nil, err
	}

	pan12digits := pan
	if len(pan) != 12 {
		pan12digits = pan[len(pan)-13 : len(pan)-1]
	}
	pinBlockDataB, err := hex.DecodeString("0000" + pan12digits)
	if err != nil {
		return nil, err
	}

	fmt.Println(pinBlockDataA, pinBlockDataB, pan12digits)

	for i, v := range pinBlockDataB {
		pinBlockDataA[i] = pinBlockDataA[i] ^ v
	}
	//log.Printf(" xor'ed pin block =", hex.EncodeToString(pin_block_data_a))

	encPinBlock, err := EncryptPinBlock(pinBlockDataA, key)
	return encPinBlock, err

}

func (pinBlock *PinBlock_Iso0) GetPin(pan string, pin_block_data []byte, key []byte) (res string, err error) {

	clearPinBlock, err := DecryptPinBlock(pin_block_data, key)

	pan12Digits := pan
	if len(pan12Digits) != 12 {
		//assume we have a full pan
		pan12Digits = pan[len(pan)-13 : len(pan)-1]
	}

	pinBlockDataB, _ := hex.DecodeString("0000" + pan12Digits)
	//log.Printf(" pin block (b) =", hex.EncodeToString(pin_block_data_b))

	for i, v := range pinBlockDataB {
		clearPinBlock[i] = clearPinBlock[i] ^ v
	}

	pinBlockStr := hex.EncodeToString(clearPinBlock)
	//log.Printf(" clear pin block (b) =", pin_block_str)

	nPinDigits, _ := strconv.ParseInt(pinBlockStr[1:2], 16, 16)
	res = pinBlockStr[2:(2 + nPinDigits)]

	return

}
