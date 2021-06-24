package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

const (
	singlelinefile = "singlecursor_test.log"
	largefile      = "largecursor_test.log"
	cursor         = "Zv-BAwEBE2VsYXN0aWNzZWFyY2hDdXJzb3IB_4IAAQUBBUxpbWl0AQQAAQlTdGFydFRpbWUB_4QAAQdFbmRUaW1lAf-EAAELU2VhcmNoQWZ0ZXIBBAABClRpZUJyZWFrZXIBDAAAAAr_gwUBAv-GAAAAKf-CAQID-gLwJpluIAEaRTdWNlYzR01TVkFNSk5LU0UyUkFSSzJVWDQA"
)

/////////////////////////// Writing Benchmarks ////////////////////////////////

// Use the open truncate command to clear the last cursor: 140149 ns/op  |   120 B/op   |    3 allocs/op
func BenchmarkSaveSingleCursorWithOpenTruncate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		singleCursorOpenWithTruncate()
	}
}

func singleCursorOpenWithTruncate() {
	file, err := os.OpenFile(singlelinefile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprintln(file, cursor)
	if err != nil {
		panic(err)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}
}

// Call truncate method on file, reset offset, and write: 20017 ns/op   |   0 B/op   |   0 allocs/op
func BenchmarkSaveSingleCursorWithTruncate(b *testing.B) {
	file, err := os.OpenFile(singlelinefile, os.O_RDWR|os.O_CREATE, 0664)
	b.ResetTimer()
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		singleCursorWithCallingTruncate(file)
	}
}

func singleCursorWithCallingTruncate(file *os.File) {
	err := file.Truncate(0)
	if err != nil {
		panic(err)
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprintln(file, cursor)
	if err != nil {
		panic(err)
	}
}

// Writing the entire history: 4440 ns/op   |   0 B/op   |    0 allocs/op
func BenchmarkSaveCursorHistory(b *testing.B) {
	store, err := OpenStore(largefile)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.SaveCursor(cursor)
	}
}

/////////////////////////// Reading Benchmarks //////////////////////////////////

// Get last line of a small file: 397149 ns/op   |   30341 B/op   |  626 allocs/op
func BenchmarkForReadingFromLargeCursorHistory(b *testing.B) {
	store, err := OpenStore(largefile)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.GetCursor()
	}
}

// Get last line of large file: 399772 ns/op   |   30345 B/op   |   626 allocs/op
func BenchmarkForReadingFromShortCursorHistory(b *testing.B) {
	store, err := OpenStore(singlelinefile)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.GetCursor()
	}
}

// Get last line of large file: 13551 ns/op   |   1096 B/op   |   5 allocs/op
func BenchmarkForReadingFromSingleCursorFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		readEntireSingleFile()
	}
}

func readEntireSingleFile() []byte {
	bytes, err := ioutil.ReadFile(singlelinefile)
	if err != nil {
		panic(err)
	}
	return bytes
}
