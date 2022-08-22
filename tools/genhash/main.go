// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"golang-object-storage/internal/pkg/utils"
	"os"
)

var (
	help = pflag.BoolP("help", "h", false, "Print this help message")
)

func main() {
	pflag.Usage = func() {
		fmt.Println(`Usage: genFileMd5 [OPTIONS] FILEPATH`)
		pflag.PrintDefaults()
	}
	pflag.Parse()

	if *help {
		pflag.Usage()
		return
	}
	if pflag.NArg() != 1 {
		pflag.Usage()
		os.Exit(1)
	}

	token := utils.GetFileHash(os.Args[1])
	if token == "" {
		fmt.Printf("Error: %s 文件不存在", os.Args[1])
		return
	}

	fmt.Println(token)
}
