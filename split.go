//go:build !windows
// +build !windows

package bonesay

import "strings"

func splitPath(s string) []string {
	return strings.Split(s, ":")
}
