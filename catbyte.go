// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var (
	inFile  *string = flag.String("input", "", "name of the input file")
	outFile *string = flag.String("output", "", "name of the output file; defaults to stdout")

	bin   *bool = flag.Bool("binary", false, "binary output")
	octal *bool = flag.Bool("octal", false, "octal output")
	hex   *bool = flag.Bool("hex", false, "hexadecimal output")
	b64   *bool = flag.Bool("base64", false, "base64 output")
)

func main() {
	flag.Parse()

	if flag.NFlag() == 0 {
		flag.Usage()
	}

	if *inFile == "" {
		log.Fatal("Enter an input file name")
	}

	readBytes, err := ioutil.ReadFile(*inFile)
	if err != nil {
		log.Fatal(err)
	}

	var outputBytes []byte
	switch {
	case *bin:
		var buf bytes.Buffer
		for b := range readBytes {
			buf.Write([]byte(fmt.Sprintf("%b", b)))
		}
		outputBytes = buf.Bytes()

	case *octal:
		var buf bytes.Buffer
		for b := range readBytes {
			buf.Write([]byte(fmt.Sprintf("%o", b)))
		}
		outputBytes = buf.Bytes()

	case *hex:
		outputBytes = []byte(fmt.Sprintf("%X", readBytes))

	case *b64:
		outputBytes = []byte(base64.StdEncoding.EncodeToString(readBytes))

	default:
		log.Fatal("Select an output format")
	}

	if *outFile == "" {
		os.Stdout.Write(outputBytes)
		os.Stdout.Write([]byte("\n"))
	} else {
		err = ioutil.WriteFile(*outFile, []byte(outputBytes), 0600)
		if err != nil {
			log.Fatal(err)
		}
	}
}
