// Copyright 2015-2016 Zack Scholl. All rights reserved.
// Use of this source code is governed by a AGPL
// license that can be found in the LICENSE file.

// utils.go is a collection of generic functions that are not specific to FIND.

package main

import (
	"bytes"
	"compress/flate"
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"log"
)

var (
	// Trace is a logging handler
	Trace *log.Logger
	// Info is a logging handler
	Info *log.Logger
	// Warning is a logging handler
	Warning *log.Logger
	// Debug is a logging handler
	Debug *log.Logger
	// Error is a logging handler
	Error *log.Logger
)

// Init function for generating the logging handlers
func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	debugHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE : ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO : ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Debug = log.New(debugHandle,
		"DEBUG: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARN : ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERR  : ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func init() {
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	// Trace.Println("I have something standard to say")
	// Info.Println("Special Information")
	// Warning.Println("There is something you need to know about")
	// Error.Println("Something has failed")
}

// GetLocalIP returns the local ip address
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "localhost"
	}
	bestIP := "localhost"
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && (strings.Contains(ipnet.IP.String(), "192.168.1") || strings.Contains(ipnet.IP.String(), "192.168")) {
				return ipnet.IP.String()
			}
		}
	}
	return bestIP
}

// stringInSlice returns boolean of whether a string is in a slice.
func stringInSlice(s string, strings []string) bool {
	for _, k := range strings {
		if s == k {
			return true
		}
	}
	return false
}

// timeTrack can be defered to provide function timing.
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	Debug.Println(name, " took ", elapsed)
}

// getMD5Hash returns a md5 hash of string.
func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// average64 computes the average of a float64 slice.
func average64(vals []float64) float64 {
	sum := float64(0)
	for _, val := range vals {
		sum += float64(val)
	}
	return sum / float64(len(vals))
}

// standardDeviation64 computes the standard deviation of a float64 slice.
func standardDeviation64(vals []float64) float64 {
	meanVal := average64(vals)

	sum := float64(0)
	for _, val := range vals {
		sum += math.Pow(float64(val)-meanVal, 2)
	}
	sum = sum / (float64(len(vals)) - 1)
	sd := math.Sqrt(sum)

	return float64(sd)
}

// standardDeviation comptues the standard deviation of a float32 slice.
func standardDeviation(vals []float32) float32 {
	sum := float64(0)
	for _, val := range vals {
		sum += float64(val)
	}
	meanVal := sum / float64(len(vals))

	sum = float64(0)
	for _, val := range vals {
		sum += math.Pow(float64(val)-meanVal, 2)
	}
	sum = sum / (float64(len(vals)) - 1)
	sd := math.Sqrt(sum)

	return float32(sd)
}

// compressByte returns a compressed byte slice.
func compressByte(src []byte) []byte {
	compressedData := new(bytes.Buffer)
	compress(src, compressedData, 9)
	return compressedData.Bytes()
}

// decompressByte returns a decompressed byte slice.
func decompressByte(src []byte) []byte {
	compressedData := bytes.NewBuffer(src)
	deCompressedData := new(bytes.Buffer)
	decompress(compressedData, deCompressedData)
	return deCompressedData.Bytes()
}

// compress uses flate to compress a byte slice to a corresponding level
func compress(src []byte, dest io.Writer, level int) {
	compressor, _ := flate.NewWriter(dest, level)
	compressor.Write(src)
	compressor.Close()
}

// compress uses flate to decompress an io.Reader
func decompress(src io.Reader, dest io.Writer) {
	decompressor := flate.NewReader(src)
	io.Copy(dest, decompressor)
	decompressor.Close()
}

// src is seeds the random generator for generating random strings
var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandStringBytesMaskImprSrc prints a random string
func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
