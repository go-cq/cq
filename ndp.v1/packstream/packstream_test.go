package packstream

import (
	"bytes"
	"io/ioutil"
	"testing"

	. "gopkg.in/check.v1"
)

type PackstreamSuite struct{}

var _ = Suite(&PackstreamSuite{})

func Test(t *testing.T) {
	TestingT(t)
}

func (s *PackstreamSuite) TestEncodeBool(c *C) {
	buf := &bytes.Buffer{}
	encoder, err := NewEncoder(buf)
	if err != nil {
		c.Fatal(err)
	}
	encoder.Encode(true)
	res, err := ioutil.ReadAll(buf)
	c.Assert(res, DeepEquals, []byte{0xC3})
	encoder.Encode(false)
	res, err = ioutil.ReadAll(buf)
	c.Assert(res, DeepEquals, []byte{0xC2})
}

func (s *PackstreamSuite) TestDecodeBool(c *C) {
	buf := &bytes.Buffer{}
	decoder, err := NewDecoder(buf)
	if err != nil {
		c.Fatal(err)
	}

	n, err := buf.Write([]byte{TRUE})
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 1)
	v, err := decoder.Decode()
	c.Assert(err, IsNil)
	c.Assert(v, Equals, true)

	n, err = buf.Write([]byte{FALSE})
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 1)
	v, err = decoder.Decode()
	c.Assert(err, IsNil)
	c.Assert(v, Equals, false)
}

func (s *PackstreamSuite) TestEncodeFloat(c *C) {
	buf := &bytes.Buffer{}
	encoder, err := NewEncoder(buf)
	if err != nil {
		c.Fatal(err)
	}
	err = encoder.Encode(1.1)
	c.Assert(err, IsNil)
	res, err := ioutil.ReadAll(buf)
	c.Assert(err, IsNil)
	c.Assert(res, DeepEquals, []byte{0xC1, 0x3F, 0xF1, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9A})
	err = encoder.Encode(-1.1)
	c.Assert(err, IsNil)
	res, err = ioutil.ReadAll(buf)
	c.Assert(err, IsNil)
	c.Assert(res, DeepEquals, []byte{0xC1, 0xBF, 0xF1, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9A})
}

var inttests = []struct {
	i   int
	arr []byte
}{
	{0, []byte{0x0}},
	{42, []byte{0x2A}},
	{-12, []byte{0xf4}},
	{-120, []byte{0xC8, 0x88}},
	{-1299, []byte{0xC9, 0xfa, 0xed}},
	{1299, []byte{0xC9, 0x05, 0x13}},
	{1234567, []byte{0xCA, 0x00, 0x12, 0xd6, 0x87}},
	{12345678900, []byte{0xCB, 0x0, 0x0, 0x0, 0x2, 0xdf, 0xdc, 0x1c, 0x34}},
}

func (s *PackstreamSuite) TestEncodeInt(c *C) {
	buf := &bytes.Buffer{}
	encoder, err := NewEncoder(buf)
	if err != nil {
		c.Fatal(err)
	}

	for _, test := range inttests {
		err = encoder.Encode(test.i)
		c.Assert(err, IsNil)
		res, err := ioutil.ReadAll(buf)
		c.Assert(err, IsNil)
		c.Assert(res, DeepEquals, test.arr)
	}
}

func (s *PackstreamSuite) TestDecodeInt(c *C) {
	buf := &bytes.Buffer{}
	decoder, err := NewDecoder(buf)
	if err != nil {
		c.Fatal(err)
	}

	for _, test := range inttests {
		n, err := buf.Write(test.arr)
		c.Assert(err, IsNil)
		c.Assert(n, Equals, len(test.arr))
		v, err := decoder.Decode()
		c.Assert(err, IsNil)
		c.Assert(v, Equals, test.i)
	}
}

type Person struct {
	Name   string
	Age    int8
	Weight int16
}

/*
func (s *PackstreamSuite) TestEncodeStruct(c *C) {
	buf := &bytes.Buffer{}
	encoder, err := NewEncoder(buf)
	if err != nil {
		c.Fatal(err)
	}

	person := Person{"wes", 12, 1234}
	err = encoder.Encode(person)
	c.Assert(err, IsNil)
}
*/
