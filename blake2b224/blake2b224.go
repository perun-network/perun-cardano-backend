// Copyright 2022, 2023 - See NOTICE file for copyright holders.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package blake2b224

import (
	"fmt"
	"golang.org/x/crypto/blake2b"
)

// Size224 is the length of blake2b224 hashes in bytes.
const Size224 = 224 / 8

// Sum224 returns the blake2b224 hash of the data.
func Sum224(data []byte) ([Size224]byte, error) {
	blake2b224, err := blake2b.New(Size224, nil)
	if err != nil {
		return [Size224]byte{}, fmt.Errorf("unable to instantiate blake2b224: %w", err)
	}
	_, err = blake2b224.Write(data)
	if err != nil {
		return [Size224]byte{}, fmt.Errorf("unable to compute blake2b224 hash of data: %w", err)
	}
	var res [Size224]byte
	n := copy(res[:], blake2b224.Sum(nil))
	if n != Size224 {
		return res, fmt.Errorf(
			"resulting blake2b224 hash has wrong length. expected: %d bytes, actual: %d bytes",
			Size224,
			n,
		)
	}
	return res, nil
}
