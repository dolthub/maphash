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
	"hash/fnv"
	gohash "hash/maphash"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeebo/xxh3"
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
	const cnt uint64 = 4096
	const mod uint64 = 4096 - 1
	data := genStringData(cnt, 16)
	b.ResetTimer()

	b.Run("dolthub/maphash", func(b *testing.B) {
		h := NewHasher[string]()
		var x uint64
		for i := 0; i < b.N; i++ {
			x = h.Hash(data[uint64(i)&mod])
		}
		assert.NotNil(b, x)
		b.ReportAllocs()
	})
	b.Run("hash/maphash", func(b *testing.B) {
		seed := gohash.MakeSeed()
		var x uint64
		for i := 0; i < b.N; i++ {
			x = gohash.String(seed, data[uint64(i)&mod])
		}
		assert.NotNil(b, x)
		b.ReportAllocs()
	})
	b.Run("xxHash3", func(b *testing.B) {
		var x uint64
		for i := 0; i < b.N; i++ {
			x = xxh3.HashStringSeed(data[uint64(i)&mod], 1)
		}
		assert.NotNil(b, x)
		b.ReportAllocs()
	})
	b.Run("hash/fnv", func(b *testing.B) {
		h := fnv.New64()
		for i := 0; i < b.N; i++ {
			_, _ = h.Write([]byte(data[uint64(i)&mod]))
		}
		assert.NotNil(b, h.Sum64())
		b.ReportAllocs()
	})
}

func BenchmarkStringPairHasher(b *testing.B) {
	type pair struct {
		a, b string
	}
	h := NewHasher[pair]()
	const cnt uint64 = 4096
	const mod uint64 = 4096 - 1
	dataA := genStringData(cnt, 16)
	dataB := genStringData(cnt, 16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Hash(pair{
			a: dataA[uint64(i)&mod],
			b: dataB[uint64(i)&mod],
		})
	}
	b.ReportAllocs()
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
