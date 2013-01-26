// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	flag "github.com/ogier/pflag"
)

var (
	inFile  *string = flag.StringP("input", "i", "", "name of the input file")
	outFile *string = flag.StringP("output", "o", "", "name of the output file; defaults to stdout")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: catbyte [options]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NFlag() == 0 {
		usage()
	}

	if *inFile == "" {
		log.Fatal("Enter an input file name")
	}

	rbytes, err := ioutil.ReadFile(*inFile)
	if err != nil {
		log.Fatal(err)
	}

	strbytes := fmt.Sprintf("%X", rbytes)

	if *outFile == "" {
		fmt.Fprintf(os.Stdout, strbytes+"\n")
		return
	}

	err = ioutil.WriteFile(*outFile, []byte(strbytes), 0600)
	if err != nil {
		log.Fatal(err)
	}
}
