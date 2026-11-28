package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/max-zit/backupzit/internal/imaging"
	"github.com/max-zit/backupzit/internal/repo"
)

// image-layout shows why a partition that can grow cannot shrink, below
// the table; the JSON output has it in "reason" already.
func TestImageLayoutShowsWhyAPartitionCannotShrink(t *testing.T) {
	img := &repo.DiskImage{Number: 0, Size: 40 << 30, Style: "gpt", Partitions: []repo.PartitionImage{
		{Number: 1, Offset: 1 << 20, Length: 100 << 20, GPTType: imaging.GPTTypeEFI, Included: true},
		{Number: 3, Offset: 101 << 20, Length: 39 << 30, GPTType: imaging.GPTTypeBasicData, Included: true, FileSystem: "NTFS"},
	}}
	const why = "cannot be shrunk: the file system in the backup was not shut down cleanly: start that Windows, shut it down fully"
	limits := map[int]imaging.PartitionLimits{
		1: {Reason: "EFI System partitions keep their size", Why: imaging.ReasonFixedSize},
		3: {Resizable: true, MinSize: 39 << 30, UsedBytes: 13 << 30, Reason: why, Why: imaging.ReasonNotClean},
	}
	var b bytes.Buffer
	printImageLayout(&b, img, limits, nil, 40<<30, 40<<30)
	out := b.String()
	if !strings.Contains(out, "3: "+why) {
		t.Errorf("reason missing:\n%s", out)
	}
	if strings.Count(out, "keep their size") != 1 {
		t.Errorf("the fixed-size reason belongs in the table, once:\n%s", out)
	}
}
