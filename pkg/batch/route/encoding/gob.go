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
)

type Gob struct{}

func (g Gob) EncodeRequest(v [][][]byte) ([][]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(v); err != nil {
		return nil, err
	}
	return [][]byte{buf.Bytes()}, nil
}

func (g Gob) DecodeResponse(payload []byte) ([][]byte, error) {
	results := make([][]byte, 0)
	if err := gob.NewDecoder(bytes.NewBuffer(payload)).Decode(&results); err != nil {
		return nil, err
	}
	return results, nil
}
