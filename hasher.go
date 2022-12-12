// Copyright 2022 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package maphash

import "unsafe"

type Hasher[K comparable] struct {
	m map[K]struct{}
}

func NewHasher[K comparable]() Hasher[K] {
	return Hasher[K]{m: make(map[K]struct{})}
}

func (h Hasher[K]) Hash(key K) uint64 {
	p := noescape(unsafe.Pointer(&key))
	return uint64(runtimeHash(p, h.m))
}
