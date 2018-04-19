package snap

import (
	"fmt"
	"strings"

	"github.com/Confbase/cfg/lib/dotcfg"
)

func Ls(lineMode bool) {
	snaps := dotcfg.MustLoadSnaps()
	if lineMode {
		for _, s := range snaps.Snapshots {
			fmt.Println(s.Name)
		}
	} else {
		snapNames := make([]string, len(snaps.Snapshots))
		for i, s := range snaps.Snapshots {
			snapNames[i] = s.Name
		}
		fmt.Println(strings.Join(snapNames, " "))
	}
}
