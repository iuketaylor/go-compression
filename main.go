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

func openFile(filename string) (*os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file '%s': %w", filename, err)
	}
	return file, nil
}

func readFileContents(file *os.File) ([]byte, error) {
	contents, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("Failed to read file': %w", err)
	}
	return contents, nil
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

func main() {
	if len(os.Args) < 2 {
		log.Fatal("You need to provide a file to compress")
		os.Exit(1)
	}

	file, err := openFile(os.Args[1])
	if err != nil {
		log.Fatal("Error opening file", err)
		os.Exit(1)
	}
	defer file.Close()

	originalBytes, err := readFileContents(file)
	if err != nil {
		log.Fatal("Error reading file contents", err)
		os.Exit(1)
	}

	compressors := []Compressor{
		{"Luke", LukeCompression},
		{"Gzip", GzipCompression},
		{"Zlib", ZlibCompression},
		{"Flate", FlateCompression},
		{"LZW", LzwCompression},
	}

	minifiedBytes, err := minifyJson(originalBytes)
	if err != nil {
		log.Fatal("Failed to minify json", err)
	}

	results := []CompressionResult{}
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

	fmt.Printf("Original Size: %d bytes\n", len(originalBytes))

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
