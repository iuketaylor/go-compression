package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"text/tabwriter"
)

type Compressor struct {
	Name string
	Func func([]byte) ([]byte, error)
}

type CompressionResult struct {
	Name           string
	OriginalSize   int
	CompressedSize int
	Saving         float64
}

func percentageDelta(b1, b2 []byte) float64 {
	len1 := float64(len(b1))
	len2 := float64(len(b2))

	return ((len1 - len2) / len1) * 100.0
}

func minifyJson(data []byte) ([]byte, error) {
	minifiedJson := new(bytes.Buffer)
	if err := json.Compact(minifiedJson, data); err != nil {
		return nil, fmt.Errorf("Failed to minify json: %w", err)
	}
	return minifiedJson.Bytes(), nil
}

func runCompressions(originalBytes, minifiedBytes []byte) []CompressionResult {
	results := []CompressionResult{}

	compressors := []Compressor{
		{"Luke", LukeCompression},
		{"Gzip", GzipCompression},
		{"Zlib", ZlibCompression},
		{"Flate", FlateCompression},
		{"LZW", LzwCompression},
	}

	for _, compressor := range compressors {
		compressedBytes, err := compressor.Func(minifiedBytes)
		if err != nil {
			log.Fatalf("Error compressing with %s: %v", compressor.Name, err)
		}

		results = append(results, CompressionResult{
			Name:           compressor.Name,
			OriginalSize:   len(originalBytes),
			CompressedSize: len(compressedBytes),
			Saving:         percentageDelta(originalBytes, compressedBytes),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Saving > results[j].Saving
	})

	return results
}

func outputResults(results []CompressionResult) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight)
	defer w.Flush()

	fmt.Fprintln(w, "Algorithm\tOriginal Size (bytes)\tCompressed Size (bytes)\tSaving (%)\t")
	fmt.Fprintln(w, "---------\t---------------------\t-----------------------\t----------\t")

	for _, result := range results {
		fmt.Fprintf(
			w,
			"%s\t%d\t%d\t%.2f%%\t\n",
			result.Name,
			result.OriginalSize,
			result.CompressedSize,
			result.Saving,
		)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("You need to provide a file to compress")
		os.Exit(1)
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error opening file", err)
	}
	defer file.Close()

	originalBytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Error reading file contents", err)
	}

	minifiedBytes, err := minifyJson(originalBytes)
	if err != nil {
		log.Fatal("Failed to minify json", err)
	}

	results := runCompressions(originalBytes, minifiedBytes)
	outputResults(results)
}
