// Copyright © 2024 Mark Summerfield. All rights reserved.
// License: GPL-3

package rbset

import (
    "fmt"
    _ "embed"
    )

//go:embed Version.dat
var Version string

func Hello() string {
    return fmt.Sprintf("Hello rbset v%s", Version)
}
