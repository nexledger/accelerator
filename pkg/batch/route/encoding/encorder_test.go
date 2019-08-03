/*
 *    Copyright 2019 Samsung SDS
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package encoding

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncoderOnCreateFailure(t *testing.T) {
	_, err := New("wrong_encoder")
	assert.Error(t, err)
}

func TestGobEncoder(t *testing.T) {
	encoder, err := New("gob")
	assert.NoError(t, err)

	args := [][][]byte{
		{[]byte("apple"), []byte("banana"), []byte("cherry")},
		{[]byte("alice"), []byte("bob"), []byte("charlie")},
		{[]byte("A"), []byte("B"), []byte("C")},
	}
	req, err := encoder.EncodeRequest(args)
	assert.NoError(t, err)

	decodedMsg := make([][][]byte, 0)
	err = gobDecode(req[0], &decodedMsg)
	assert.NoError(t, err)
	assert.Equal(t, args, decodedMsg)

	encodedMsg, err := gobEncode(args[0])
	resp, err := encoder.DecodeResponse(encodedMsg)
	assert.NoError(t, err)
	assert.Equal(t, args[0], resp)
}

func TestJsonEncoder(t *testing.T) {
	encoder, err := New("json")
	assert.NoError(t, err)

	args := [][][]byte{
		{[]byte("apple"), []byte("banana"), []byte("cherry")},
		{[]byte("alice"), []byte("bob"), []byte("charlie")},
		{[]byte("A"), []byte("B"), []byte("C")},
	}
	req, err := encoder.EncodeRequest(args)
	assert.NoError(t, err)

	decodedMsg := make([][][]byte, 0)
	err = jsonDecode(req[0], &decodedMsg)
	assert.NoError(t, err)
	assert.Equal(t, args, decodedMsg)

	encodedMsg, err := jsonEncode(args[0])
	resp, err := encoder.DecodeResponse(encodedMsg)
	assert.NoError(t, err)
	assert.Equal(t, args[0], resp)
}

func gobEncode(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func gobDecode(d []byte, v interface{}) error {
	buf := bytes.NewBuffer(d)
	if err := gob.NewDecoder(buf).Decode(v); err != nil {
		return err
	}
	return nil
}

func jsonEncode(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func jsonDecode(d []byte, v interface{}) error {
	buf := bytes.NewBuffer(d)
	if err := json.NewDecoder(buf).Decode(v); err != nil {
		return err
	}
	return nil
}
