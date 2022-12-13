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

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func FuzzStringHasher(f *testing.F) {
	f.Add("")
	f.Add("hello world")
	f.Add("github.com/dolthub/maphash")
	f.Fuzz(func(t *testing.T, key string) {
		testHasher(t, key)
	})
}

func FuzzRuneHasher(f *testing.F) {
	f.Add('a')
	f.Add('z')
	f.Add('A')
	f.Add('Z')
	f.Fuzz(func(t *testing.T, key rune) {
		testHasher(t, key)
	})
}

func FuzzIntHasher(f *testing.F) {
	f.Add(int(0))
	f.Add(int(math.MaxInt32))
	f.Fuzz(func(t *testing.T, key int) {
		testHasher(t, key)
	})
}

func FuzzInt8Hasher(f *testing.F) {
	f.Add(int8(0))
	f.Add(int8(math.MaxInt8))
	f.Fuzz(func(t *testing.T, key int8) {
		testHasher(t, key)
	})
}

func FuzzInt16Hasher(f *testing.F) {
	f.Add(int16(0))
	f.Add(int16(math.MaxInt16))
	f.Fuzz(func(t *testing.T, key int16) {
		testHasher(t, key)
	})
}

func FuzzInt32Hasher(f *testing.F) {
	f.Add(int32(0))
	f.Add(int32(math.MaxInt32))
	f.Fuzz(func(t *testing.T, key int32) {
		testHasher(t, key)
	})
}

func FuzzInt64Hasher(f *testing.F) {
	f.Add(int64(0))
	f.Add(int64(math.MaxInt64))
	f.Fuzz(func(t *testing.T, key int64) {
		testHasher(t, key)
	})
}

func FuzzUintHasher(f *testing.F) {
	f.Add(uint(0))
	f.Add(uint(math.MaxUint32))
	f.Fuzz(func(t *testing.T, key uint) {
		testHasher(t, key)
	})
}

func FuzzUint8Hasher(f *testing.F) {
	f.Add(uint8(0))
	f.Add(uint8(math.MaxUint8))
	f.Fuzz(func(t *testing.T, key uint8) {
		testHasher(t, key)
	})
}

func FuzzUint16Hasher(f *testing.F) {
	f.Add(uint16(0))
	f.Add(uint16(math.MaxUint16))
	f.Fuzz(func(t *testing.T, key uint16) {
		testHasher(t, key)
	})
}

func FuzzUint32Hasher(f *testing.F) {
	f.Add(uint32(0))
	f.Add(uint32(math.MaxUint32))
	f.Fuzz(func(t *testing.T, key uint32) {
		testHasher(t, key)
	})
}

func FuzzUint64Hasher(f *testing.F) {
	f.Add(uint64(0))
	f.Add(uint64(math.MaxUint64))
	f.Fuzz(func(t *testing.T, key uint64) {
		testHasher(t, key)
	})
}
func FuzzFloat32Hasher(f *testing.F) {
	f.Add(float32(0))
	f.Add(float32(math.Pi))
	f.Add(float32(math.E))
	f.Fuzz(func(t *testing.T, key float32) {
		testHasher(t, key)
	})
}

func FuzzFloat64Hasher(f *testing.F) {
	f.Add(float64(0))
	f.Add(float64(math.Pi))
	f.Add(float64(math.E))
	f.Fuzz(func(t *testing.T, key float64) {
		testHasher(t, key)
	})
}

func FuzzArrayHasher(f *testing.F) {
	f.Add(0, 0)
	f.Add(1, -1)
	f.Fuzz(func(t *testing.T, a, b int) {
		testHasher(t, [2]int{a, b})
	})
}

func FuzzStructHasher(f *testing.F) {
	type obj struct {
		i int
		f float32
		t time.Time
	}
	f.Add(int(0), float32(0), int64(0))
	f.Add(int(-1), float32(-1), int64(-1))
	f.Fuzz(func(t *testing.T, i int, f float32, m int64) {
		o := obj{i: i, f: f, t: time.UnixMicro(m)}
		testHasher(t, o)
	})
}

func FuzzStringPairHasher(f *testing.F) {
	type pair struct {
		a, b string
	}
	f.Add("", "")
	f.Add("a", "b")
	f.Add("hello", "world")
	f.Fuzz(func(t *testing.T, a, b string) {
		testHasher(t, pair{a, b})
	})
}

func testHasher[K comparable](t *testing.T, key K) {
	h1 := NewHasher[K]().Hash(key)
	h2 := NewHasher[K]().Hash(key)
	assert.NotEqual(t, h1, h2) // new seed
}

func TestRefAllocs(t *testing.T) {
	t.Run("*int", func(t *testing.T) {
		x := int(42)
		testNoAllocs(t, NewHasher[*int](), &x)
	})
	t.Run("*uint", func(t *testing.T) {
		x := uint(42)
		testNoAllocs(t, NewHasher[*uint](), &x)
	})
	t.Run("*float", func(t *testing.T) {
		x := float64(math.E)
		testNoAllocs(t, NewHasher[*float64](), &x)
	})
	t.Run("*string", func(t *testing.T) {
		x := string("asdf")
		testNoAllocs(t, NewHasher[*string](), &x)
	})
}

func testSomeAllocs[K comparable](t *testing.T, h Hasher[K], key K) {
	a := testing.AllocsPerRun(64, func() {
		h.Hash(key)
	})
	assert.True(t, a >= 1.0)
}

func TestNoValueAllocs(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		testNoAllocs(t, NewHasher[int](), 42)
	})
	t.Run("uint", func(t *testing.T) {
		testNoAllocs(t, NewHasher[uint](), 42)
	})
	t.Run("float", func(t *testing.T) {
		testNoAllocs(t, NewHasher[float64](), math.E)
	})
	t.Run("string", func(t *testing.T) {
		testNoAllocs(t, NewHasher[string](), "asdf")
	})
	type uuid [16]byte
	t.Run("uuid", func(t *testing.T) {
		testNoAllocs(t, NewHasher[uuid](), uuid{})
	})
	t.Run("time", func(t *testing.T) {
		testNoAllocs(t, NewHasher[time.Time](), time.Now())
	})
}

func testNoAllocs[K comparable](t *testing.T, h Hasher[K], key K) {
	a := testing.AllocsPerRun(64, func() {
		h.Hash(key)
	})
	assert.Equal(t, 0.0, a)
}
