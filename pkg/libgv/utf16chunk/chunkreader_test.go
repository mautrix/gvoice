// Copyright (c) 2024 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package utf16chunk_test

import (
	"io"
	"testing"
	"unicode/utf16"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.mau.fi/mautrix-gvoice/pkg/libgv/utf16chunk"
)

type fakeReader struct {
	chunks [][]byte
}

func (fr *fakeReader) Read(p []byte) (n int, err error) {
	if len(fr.chunks) == 0 {
		return 0, io.EOF
	}
	n = copy(p, fr.chunks[0])
	if n == len(fr.chunks[0]) {
		fr.chunks = fr.chunks[1:]
	} else {
		fr.chunks[0] = fr.chunks[0][n:]
	}
	return
}

type testCase struct {
	name     string
	chunks   []string
	expected []string
}

var testCases = []testCase{
	{
		name:     "SingleChunk",
		chunks:   []string{"5\nhello"},
		expected: []string{"hello"},
	},
	{
		name:     "SingleChunkEmoji",
		chunks:   []string{"3\n🐈️"},
		expected: []string{"🐈️"},
	},
	{
		name:     "MultipleChunksInOne",
		chunks:   []string{"5\nhello5\nworld"},
		expected: []string{"hello", "world"},
	},
	{
		name:     "MultipleChunksInOneWithEmoji",
		chunks:   []string{"5\nhello3\n🐈️"},
		expected: []string{"hello", "🐈️"},
	},
	{
		name:     "OneChunkInMultipleParts",
		chunks:   []string{"15\nhello 🐈️ world"},
		expected: []string{"hello 🐈️ world"},
	},
	{
		name:     "ChunkWithNewline",
		chunks:   []string{"15\nhello\n🐈️\nworld"},
		expected: []string{"hello\n🐈️\nworld"},
	},
}

func TestReader_ReadChunk(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fr := &fakeReader{chunks: make([][]byte, len(tc.chunks))}
			for i, chunk := range tc.chunks {
				fr.chunks[i] = []byte(chunk)
			}
			cbr := utf16chunk.NewReader(fr)
			for _, expected := range tc.expected {
				chunk, err := cbr.ReadChunk()
				require.NoError(t, err)
				require.Equal(t, expected, string(chunk))
			}
			_, err := cbr.ReadChunk()
			require.ErrorIs(t, err, io.EOF)
		})
	}
}

func FuzzUTF16Length(f *testing.F) {
	fuzzTests := []string{
		"meow",
		"🐈️",
		"åäö",
		"🐈️åäö",
		"🐈️🐈️🐈️🐈️🐈️🐈️🐈️🐈️🐈️🐈️🐈️🐈️🐈️🐈️🐈️",
		"世界",
		"🤔",
		"\U00010000",
	}
	for _, val := range fuzzTests {
		f.Add(val)
	}
	f.Fuzz(func(t *testing.T, input string) {
		expected := utf16.Encode([]rune(input))
		got, _ := utf16chunk.UTF16Length([]byte(input), 1<<31)
		assert.Equal(t, len(expected), got, "input: %q", input)
	})
}
