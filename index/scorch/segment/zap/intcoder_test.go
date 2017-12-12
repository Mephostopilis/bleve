package zap

import (
	"bytes"
	"reflect"
	"testing"
)

func TestChunkIntCoder(t *testing.T) {
	tests := []struct {
		maxDocNum uint64
		chunkSize uint64
		docNums   []uint64
		vals      [][]uint64
		expected  []byte
	}{
		{
			maxDocNum: 0,
			chunkSize: 1,
			docNums:   []uint64{0},
			vals: [][]uint64{
				[]uint64{3},
			},
			// 1 chunk, chunk-0 length 1, value 3
			expected: []byte{0x1, 0x1, 0x3},
		},
		{
			maxDocNum: 1,
			chunkSize: 1,
			docNums:   []uint64{0, 1},
			vals: [][]uint64{
				[]uint64{3},
				[]uint64{7},
			},
			// 2 chunks, chunk-0 length 1, chunk-1 length 1, value 3, value 7
			expected: []byte{0x2, 0x1, 0x1, 0x3, 0x7},
		},
	}

	for _, test := range tests {

		cic := newChunkedIntCoder(test.chunkSize, test.maxDocNum)
		for i, docNum := range test.docNums {
			err := cic.Add(docNum, test.vals[i]...)
			if err != nil {
				t.Fatalf("error adding to intcoder: %v", err)
			}
		}
		cic.Close()
		var actual bytes.Buffer
		_, err := cic.Write(&actual)
		if err != nil {
			t.Fatalf("error writing: %v", err)
		}
		if !reflect.DeepEqual(test.expected, actual.Bytes()) {
			t.Errorf("got % x, expected % x", actual.Bytes(), test.expected)
		}
	}
}
