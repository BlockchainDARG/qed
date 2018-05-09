// Copyright © 2018 Banco Bilbao Vizcaya Argentaria S.A.  All rights reserved.
// Use of this source code is governed by an Apache 2 License
// that can be found in the LICENSE file

package main

import (
	"os"

	"qed/cli"
)

func main() {
	ctx := cli.NewContext()
	if err := cli.NewQedCommand(ctx).Execute(); err != nil {
		os.Exit(-1)
	}
}
