package test

import (
	"testing"
	"os"
	"github.com/jmacdonald/liberator/filesystem/directory"
)

func TestSize(t *testing.T) {
	var expectedSize int64 = 512020

	// Call the Size function and make sure it returns the expected value.
	dir, _ := os.Getwd()
	if x := directory.Size(dir + "/sample"); x != expectedSize {
		t.Errorf("directory.Size(%v) should return %v, but returned %v instead.", dir, expectedSize, x)
	}
}
