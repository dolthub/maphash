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

//go:build go1.19
// +build go1.19

package maphash

import (
	gohash "hash/maphash"
	"math/rand"
	"testing"
)

func BenchmarkUint64Hasher(b *testing.B) {
	h := NewHasher[uint64]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Hash(uint64(i))
	}
	b.ReportAllocs()
}

func BenchmarkCompareStringHasher(b *testing.B) {
	h := NewHasher[string]()
	seed := gohash.MakeSeed()
	const cnt uint64 = 4096
	const mod uint64 = 4096 - 1
	data := genStringData(cnt, 16)
	b.ResetTimer()

	b.Run("string hasher", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			h.Hash(data[uint64(i)&mod])
		}
		b.ReportAllocs()
	})
	b.Run("std string hasher", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gohash.String(seed, data[uint64(i)&mod])
		}
		b.ReportAllocs()
	})
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func genStringData(cnt, ln uint64) (data []string) {
	str := func(n uint64) string {
		s := make([]rune, n)
		for i := range s {
			s[i] = letters[rand.Intn(52)]
		}
		return string(s)
	}
	data = make([]string, cnt)
	for i := range data {
		data[i] = str(ln)
	}
	return
}
