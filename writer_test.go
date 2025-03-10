// Copyright 2021 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bstream

import (
	"bytes"
	"io"
	"testing"
	"time"

	pbbstream "github.com/streamingfast/pbgo/sf/bstream/v1"
	"github.com/stretchr/testify/require"
)

func TestBlockWriter(t *testing.T) {
	writerFactory := BlockWriterFactoryFunc(func(writer io.Writer) (BlockWriter, error) { return NewDBinBlockWriter(writer, "tst", 1) })

	buffer := bytes.NewBuffer([]byte{})
	blockWriter, err := writerFactory.New(buffer)
	require.NoError(t, err)

	block1Payload := []byte{0x0a, 0x0b, 0x0c}

	blk1 := &Block{
		Id:             "0a",
		Number:         1,
		PreviousId:     "0b",
		Timestamp:      time.Date(1970, time.December, 31, 19, 0, 0, 0, time.UTC),
		LibNum:         0,
		PayloadKind:    pbbstream.Protocol_ETH,
		PayloadVersion: 1,
		Payload:        &MemoryBlockPayload{data: block1Payload},
	}

	err = blockWriter.Write(blk1)
	require.NoError(t, err)

	// Reader part (to validate the data)

	var readerFactory BlockReaderFactory = BlockReaderFactoryFunc(func(reader io.Reader) (BlockReader, error) { return NewDBinBlockReader(reader, nil) })
	blockReader, err := readerFactory.New(buffer)
	require.NoError(t, err)

	readBlk1, err := blockReader.Read()
	require.Equal(t, blk1, readBlk1)
	require.NoError(t, err)
}
