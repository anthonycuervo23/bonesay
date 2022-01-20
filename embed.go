package bonesay

import (
	"embed"
	"sort"
	"strings"
)

//go:embed_bones/*
var bonesDir embed.FS

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(path string) ([]byte, error) {
	return bonesDir.ReadFile(path)
}

// AssetNames returns the list of filename of the assets.
func AssetNames() []string {
	entries, err := bonesDir.ReadDir("bones")
	if err != nil {
		panic(err)
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		name := strings.TrimSuffix(entry.Name(), ".bone")
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

var bonesInBinary = AssetNames()

// BonesInBinary returns the list of bonefiles which are in binary.
// the list is memoized.
func BonesInBinary() []string {
	return bonesInBinary
}
