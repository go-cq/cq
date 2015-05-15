package packstream

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
)

const (
	TINY_TEXT_START = 0x80
	TINY_TEXT_END   = 0x8F
	TEXT_SIZE_8     = 0xD0
	TEXT_SIZE_16    = 0xD1
	TEXT_SIZE_32    = 0xD2

	TINY_LIST_START = 0x90
	TINY_LIST_END   = 0x9F
	LIST_SIZE_8     = 0xD4
	LIST_SIZE_16    = 0xD5
	LIST_SIZE_32    = 0xD6

	TINY_MAP_START = 0xA0
	TINY_MAP_END   = 0xAF
	MAP_SIZE_8     = 0xD8
	MAP_SIZE_16    = 0xD9
	MAP_SIZE_32    = 0xDA

	TINY_STRUCT_START = 0xB0
	TINY_STRUCT_END   = 0xBF
	STRUCT_SIZE_8     = 0xDC
	STRUCT_SIZE_16    = 0xDD

	NULL     = 0xC0
	FLOAT_64 = 0xC1
	FALSE    = 0xC2
	TRUE     = 0xC3

	MIN_TINY_INT = -16
	MAX_TINY_INT = 127
	INT_8        = 0xC8
	INT_16       = 0xC9
	INT_32       = 0xCA
	INT_64       = 0xCB
)

type Decoder struct {
	br *bufio.Reader
}

func NewDecoder(reader io.Reader) (*Decoder, error) {
	dec := Decoder{bufio.NewReader(reader)}
	return &dec, nil
}

type Encoder struct {
	bw *bufio.Writer
}

func NewEncoder(writer io.Writer) (*Encoder, error) {
	enc := Encoder{bufio.NewWriter(writer)}
	return &enc, nil
}

func (dec Decoder) Decode() (interface{}, error) {
	b, err := dec.br.ReadByte()
	if err != nil {
		return nil, err
	}
	if int8(b) >= MIN_TINY_INT && int8(b) <= MAX_TINY_INT {
		return int(int8(b)), nil
	}
	switch b {
	case FALSE:
		return false, nil
	case TRUE:
		return true, nil
	case INT_8:
		i, err := dec.br.ReadByte()
		return int(int8(i)), err
	case INT_16:
		barr := [2]byte{}
		n, err := dec.br.Read(barr[:])
		if err != nil {
			return nil, err
		}
		if n != 2 {
			return nil, errors.New("not enough bytes read")
		}
		v := int16(0)
		v |= int16(barr[0]) << 8
		v |= int16(barr[1])
		return int(int16(v)), nil
	case INT_32:
		barr := [4]byte{}
		n, err := dec.br.Read(barr[:])
		if err != nil {
			return nil, err
		}
		if n != 4 {
			return nil, errors.New("not enough bytes read")
		}
		v := int32(0)
		v |= int32(barr[0]) << 24
		v |= int32(barr[1]) << 16
		v |= int32(barr[2]) << 8
		v |= int32(barr[3])
		return int(int32(v)), nil
	case INT_64:
		barr := [8]byte{}
		n, err := dec.br.Read(barr[:])
		if err != nil {
			return nil, err
		}
		if n != 8 {
			return nil, errors.New("not enough bytes read")
		}
		v := int64(0)
		v |= int64(barr[0]) << 56
		v |= int64(barr[1]) << 48
		v |= int64(barr[2]) << 40
		v |= int64(barr[3]) << 32
		v |= int64(barr[4]) << 24
		v |= int64(barr[5]) << 16
		v |= int64(barr[6]) << 8
		v |= int64(barr[7])
		return int(int64(v)), nil
	}
	return nil, errors.New(fmt.Sprintf("decode: unsupported type: %x", b))
}

func (enc Encoder) Encode(v interface{}) error {
	switch v.(type) {
	case bool:
		return enc.encodeBool(v.(bool))
	case float64:
		return enc.encodeFloat64(v.(float64))
	case int64:
		return enc.encodeInt64(v.(int64))
	case int:
		return enc.encodeInt64(int64(v.(int)))
	}
	return errors.New("unsupported type")
}

func (enc Encoder) encodeBool(b bool) error {
	var n int
	var err error
	if b {
		n, err = enc.bw.Write([]byte{TRUE})
	} else {
		n, err = enc.bw.Write([]byte{FALSE})
	}
	if err != nil {
		return err

	}
	if n == 0 {
		return errors.New("failed to encodebool")
	}
	err = enc.bw.Flush()
	return err
}

func (enc Encoder) encodeFloat64(f float64) error {
	v := math.Float64bits(f)
	b := [9]byte{}
	b[0] = FLOAT_64
	b[1] = byte(v >> 56)
	b[2] = byte(v >> 48)
	b[3] = byte(v >> 40)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 24)
	b[6] = byte(v >> 16)
	b[7] = byte(v >> 8)
	b[8] = byte(v)
	n, err := enc.bw.Write(b[:])
	if err != nil {
		return err
	}
	if n != 9 {
		return errors.New("not enough bytes written")
	}
	enc.bw.Flush()
	return err
}

func (enc Encoder) encodeInt64(i int64) error {
	switch {
	case (i >= -9223372036854775808 && i <= -2147483649) || (i >= 2147483648 && i <= 9223372036854775807): // INT_64
		b := [9]byte{}
		b[0] = INT_64
		b[1] = byte(i >> 56)
		b[2] = byte(i >> 48)
		b[3] = byte(i >> 40)
		b[4] = byte(i >> 32)
		b[5] = byte(i >> 24)
		b[6] = byte(i >> 16)
		b[7] = byte(i >> 8)
		b[8] = byte(i)
		n, err := enc.bw.Write(b[:])
		if err != nil {
			return err
		}
		if n != 9 {
			return errors.New("not enough bytes written")
		}
		enc.bw.Flush()
		return err
	case (i >= -2147483648 && i <= -32769) || (i >= 32768 && i <= 2147483647): // INT_32
		return enc.encodeInt32(int(i))
	case (i >= -32768 && i <= -129) || (i >= 128 && i <= 32767): // INT_16
		return enc.encodeInt16(int(i))
	case i >= -128 && i <= -17: // INT_8
		return enc.encodeInt8(int(i))
	case i >= MIN_TINY_INT && i <= MAX_TINY_INT: // TINY_INT
		return enc.encodeTinyInt(int(i))
	}
	return errors.New("invalid int64: this should not happen")
}

func (enc Encoder) encodeTinyInt(i int) error {
	if i < MIN_TINY_INT || i > MAX_TINY_INT {
		return errors.New("encode tinyint: out of range")
	}
	b := byte(i)
	n, err := enc.bw.Write([]byte{b})
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("not enough bytes written")
	}
	enc.bw.Flush()
	return err
}

func (enc Encoder) encodeInt8(i int) error {
	if i < -128 || i > 127 {
		return errors.New("encode int8: out of range")
	}
	n, err := enc.bw.Write([]byte{INT_8, byte(i)})
	if err != nil {
		return err
	}
	if n != 2 {
		return errors.New("not enough bytes written")
	}
	enc.bw.Flush()
	return err
}

func (enc Encoder) encodeInt16(i int) error {
	if i < -32768 || i > 32767 {
		return errors.New("encode int16: out of range")
	}
	b := [3]byte{}
	b[0] = INT_16
	b[1] = byte(i >> 8)
	b[2] = byte(i)
	n, err := enc.bw.Write(b[:])
	if err != nil {
		return err
	}
	if n != 3 {
		return errors.New("not enough bytes written")
	}
	enc.bw.Flush()
	return err
}

func (enc Encoder) encodeInt32(i int) error {
	if i < -2147483648 || i > 2147483647 {
		return errors.New("encode int32: out of range")
	}
	b := [5]byte{}
	b[0] = INT_32
	b[1] = byte(i >> 24)
	b[2] = byte(i >> 16)
	b[3] = byte(i >> 8)
	b[4] = byte(i)
	n, err := enc.bw.Write(b[:])
	if err != nil {
		return err
	}
	if n != 5 {
		return errors.New("not enough bytes written")
	}
	enc.bw.Flush()
	return err
}
