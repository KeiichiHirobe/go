// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var gooses, goarches []string

func main() {
	data, err := os.ReadFile("../../../go/build/syslist.go")
	if err != nil {
		log.Fatal(err)
	}
	const (
		goosPrefix   = `const goosList = `
		goarchPrefix = `const goarchList = `
	)
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, goosPrefix) {
			text, err := strconv.Unquote(strings.TrimPrefix(line, goosPrefix))
			if err != nil {
				log.Fatalf("parsing goosList: %v", err)
			}
			gooses = strings.Fields(text)
		}
		if strings.HasPrefix(line, goarchPrefix) {
			text, err := strconv.Unquote(strings.TrimPrefix(line, goarchPrefix))
			if err != nil {
				log.Fatalf("parsing goarchList: %v", err)
			}
			goarches = strings.Fields(text)
		}
	}

	for _, target := range gooses {
		if target == "nacl" {
			continue
		}
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "// Code generated by gengoos.go using 'go generate'. DO NOT EDIT.\n\n")
		if target == "linux" {
			fmt.Fprintf(&buf, "// +build !android\n") // must explicitly exclude android for linux
		}
		if target == "solaris" {
			fmt.Fprintf(&buf, "// +build !illumos\n") // must explicitly exclude illumos for solaris
		}
		if target == "darwin" {
			fmt.Fprintf(&buf, "// +build !ios\n") // must explicitly exclude ios for darwin
		}
		fmt.Fprintf(&buf, "// +build %s\n\n", target) // must explicitly include target for bootstrapping purposes
		fmt.Fprintf(&buf, "package sys\n\n")
		fmt.Fprintf(&buf, "const GOOS = `%s`\n\n", target)
		for _, goos := range gooses {
			value := 0
			if goos == target {
				value = 1
			}
			fmt.Fprintf(&buf, "const Goos%s = %d\n", strings.Title(goos), value)
		}
		err := os.WriteFile("zgoos_"+target+".go", buf.Bytes(), 0666)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, target := range goarches {
		if target == "amd64p32" {
			continue
		}
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "// Code generated by gengoos.go using 'go generate'. DO NOT EDIT.\n\n")
		fmt.Fprintf(&buf, "// +build %s\n\n", target) // must explicitly include target for bootstrapping purposes
		fmt.Fprintf(&buf, "package sys\n\n")
		fmt.Fprintf(&buf, "const GOARCH = `%s`\n\n", target)
		for _, goarch := range goarches {
			value := 0
			if goarch == target {
				value = 1
			}
			fmt.Fprintf(&buf, "const Goarch%s = %d\n", strings.Title(goarch), value)
		}
		err := os.WriteFile("zgoarch_"+target+".go", buf.Bytes(), 0666)
		if err != nil {
			log.Fatal(err)
		}
	}
}