// Copyright (c) 2024 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package utf16chunk

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"unicode/utf8"
)

type Reader struct {
	r   io.Reader
	buf []byte
	ptr int
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		r:   r,
		buf: make([]byte, 32*1024),
	}
}

const utf16MaxSinglePair = 0x10000 // matches surrSelf in encoding/utf16/utf16.go

func UTF16Length(utf8Bytes []byte, maxLength int) (length, i int) {
	if maxLength < 0 {
		maxLength = len(utf8Bytes)
	}
	for i < len(utf8Bytes) {
		r, size := utf8.DecodeRune(utf8Bytes[i:])
		if r >= utf16MaxSinglePair {
			length += 2
		} else {
			length += 1
		}
		i += size
		if length >= maxLength {
			break
		}
	}
	return
}

const LengthDataDelimiter = '\n'

func (cr *Reader) ReadChunk() ([]byte, error) {
	idx := -1
	n := cr.ptr
	if cr.ptr != 0 {
		idx = bytes.IndexByte(cr.buf[:n], LengthDataDelimiter)
	}
	for idx == -1 && n < 10 {
		newRead, err := cr.r.Read(cr.buf[cr.ptr:])
		if err != nil {
			return nil, err
		}
		n += newRead
		idx = bytes.IndexByte(cr.buf[:n], LengthDataDelimiter)
	}
	if idx == -1 {
		return nil, fmt.Errorf("no newline found in chunk")
	}
	expectedLength, err := strconv.Atoi(string(cr.buf[:idx]))
	if err != nil {
		return nil, fmt.Errorf("failed to parse chunk length: %w", err)
	}
	receivedLength, i := UTF16Length(cr.buf[idx+1:n], expectedLength)
	output := make([]byte, i)
	copy(output, cr.buf[idx+1:idx+1+i])
	if receivedLength >= expectedLength {
		cr.ptr = copy(cr.buf, cr.buf[idx+1+i:n])
		return output, nil
	}
	expectedLength -= receivedLength
	data := bytes.NewBuffer(output)
	cr.ptr = 0
	for expectedLength > 0 {
		n, err = cr.r.Read(cr.buf[cr.ptr:])
		if err != nil {
			return nil, err
		}
		receivedLength, i = UTF16Length(cr.buf[:n], expectedLength)
		expectedLength -= receivedLength
		data.Write(cr.buf[:i])
		if i < n {
			cr.ptr = copy(cr.buf, cr.buf[i:n])
		}
	}
	return data.Bytes(), nil
}
