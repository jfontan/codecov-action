package main

import (
	"testing"
	"time"
)

func TestDummy(t *testing.T) {
	_ = newGithubClient("", time.Second)
}
