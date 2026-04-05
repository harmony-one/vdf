package main

import (
	"fmt"
	"io"
	"log"
	"testing"

	"github.com/harmony-one/vdf/src/vdf_go"
)

const benchmarkIntSizeBits = 2048

var benchmarkTargets = []int{500, 1000, 10000}

var benchmarkSeed = []byte{
	0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
	0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
	0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
	0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
}

func silenceVDFLogs(b *testing.B) {
	b.Helper()

	originalWriter := log.Writer()
	log.SetOutput(io.Discard)
	b.Cleanup(func() {
		log.SetOutput(originalWriter)
	})
}

func BenchmarkGenerateVDFByTarget(b *testing.B) {
	silenceVDFLogs(b)

	for _, target := range benchmarkTargets {
		target := target
		b.Run(fmt.Sprintf("target_%d", target), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				yBuf, proofBuf := vdf_go.GenerateVDF(benchmarkSeed, target, benchmarkIntSizeBits)
				if yBuf == nil || proofBuf == nil {
					b.Fatalf("GenerateVDF returned nil output for target=%d", target)
				}
			}
		})
	}
}

func BenchmarkVerifyVDFByTarget(b *testing.B) {
	silenceVDFLogs(b)

	for _, target := range benchmarkTargets {
		target := target
		b.Run(fmt.Sprintf("target_%d", target), func(b *testing.B) {
			yBuf, proofBuf := vdf_go.GenerateVDF(benchmarkSeed, target, benchmarkIntSizeBits)
			if yBuf == nil || proofBuf == nil {
				b.Fatalf("GenerateVDF returned nil output for target=%d", target)
			}

			proofBlob := append(yBuf, proofBuf...)
			if !vdf_go.VerifyVDF(benchmarkSeed, proofBlob, target, benchmarkIntSizeBits) {
				b.Fatalf("precomputed proof is invalid for target=%d", target)
			}

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				if !vdf_go.VerifyVDF(benchmarkSeed, proofBlob, target, benchmarkIntSizeBits) {
					b.Fatalf("VerifyVDF failed for target=%d", target)
				}
			}
		})
	}
}
