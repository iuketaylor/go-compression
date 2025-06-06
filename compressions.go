package main

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"compress/zlib"
	"fmt"
)

func GzipCompression(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(data); err != nil {
		return nil, fmt.Errorf("failed to write to gzip writer: %w", err)
	}
	if err := gz.Close(); err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %w", err)
	}
	return b.Bytes(), nil
}

func ZlibCompression(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	if _, err := w.Write(data); err != nil {
		return nil, fmt.Errorf("failed to write to zlib writer: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed to close zlib writer: %w", err)
	}
	return b.Bytes(), nil
}

func FlateCompression(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w, err := flate.NewWriter(&b, flate.DefaultCompression)
	if err != nil {
		return nil, fmt.Errorf("failed to create flate writer: %w", err)
	}
	if _, err := w.Write(data); err != nil {
		return nil, fmt.Errorf("failed to write to flate writer: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed to close flate writer: %w", err)
	}
	return b.Bytes(), nil
}

func LzwCompression(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := lzw.NewWriter(&b, lzw.LSB, 8)
	if _, err := w.Write(data); err != nil {
		return nil, fmt.Errorf("failed to write to lzw writer: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed to close lzw writer: %w", err)
	}
	return b.Bytes(), nil
}

func LukeCompression(data []byte) ([]byte, error) {
	return data, nil
}
