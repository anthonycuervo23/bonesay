package bonesay

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Say to return bonesay string.
func Say(phrase string, options ...Option) (string, error) {
	bone, err := New(options...)
	if err != nil {
		return "", err
	}
	return bone.Say(phrase)
}

// LocationType indicates the type of BONEPATH.
type LocationType int

const (
	// InBinary indicates the BONEPATH in binary.
	InBinary LocationType = iota

	// InDirectory indicates the BONEPATH in your directory.
	InDirectory
)

// BonePath is information of the BONEPATH.
type BonePath struct {
	// Name is name of the BONEPATH.
	// If you specified `BONEPATH=/foo/bar`, Name is `/foo/bar`.
	Name string
	// BoneFiles are name of the bonefile which are trimmed ".bone" suffix.
	BoneFiles []string
	// LocationType is the type of BONEPATH
	LocationType LocationType
}

// Lookup will look for the target bonefile in the specified path.
// If it exists, it returns the bonefile information and true value.
func (c *BonePath) Lookup(target string) (*BoneFile, bool) {
	for _, bonefile := range c.BoneFiles {
		if bonefile == target {
			return &BoneFile{
				Name:         bonefile,
				BasePath:     c.Name,
				LocationType: c.LocationType,
			}, true
		}
	}
	return nil, false
}

// BoneFile is information of the bonefile.
type BoneFile struct {
	// Name is name of the bonefile.
	Name string
	// BasePath is the path which the bonepath is in.
	BasePath string
	// LocationType is the type of BONEPATH
	LocationType LocationType
}

// ReadAll reads the bonefile content.
// If LocationType is InBinary, the file read from binary.
// otherwise reads from file system.
func (c *BoneFile) ReadAll() ([]byte, error) {
	joinedPath := filepath.Join(c.BasePath, c.Name+".bone")
	if c.LocationType == InBinary {
		return Asset(joinedPath)
	}
	return ioutil.ReadFile(joinedPath)
}

// Bones to get list of bones
func Bones() ([]*BonePath, error) {
	bonePaths, err := bonesFromBonePath()
	if err != nil {
		return nil, err
	}
	bonePaths = append(bonePaths, &BonePath{
		Name:         "bones",
		BoneFiles:    BonesInBinary(),
		LocationType: InBinary,
	})
	return bonePaths, nil
}

func bonesFromBonePath() ([]*BonePath, error) {
	bonePaths := make([]*BonePath, 0)
	bonePath := os.Getenv("BONEATH")
	if bonePath == "" {
		return bonePaths, nil
	}
	paths := splitPath(bonePath)
	for _, path := range paths {
		dirEntries, err := ioutil.ReadDir(path)
		if err != nil {
			return nil, err
		}
		path := &BonePath{
			Name:         path,
			BoneFiles:    []string{},
			LocationType: InDirectory,
		}
		for _, entry := range dirEntries {
			name := entry.Name()
			if strings.HasSuffix(name, ".bone") {
				name = strings.TrimSuffix(name, ".bone")
				path.BoneFiles = append(path.BoneFiles, name)
			}
		}
		sort.Strings(path.BoneFiles)
		bonePaths = append(bonePaths, path)
	}
	return bonePaths, nil
}

// GetBone to get bone's ascii art
func (bone *Bone) GetBone() (string, error) {
	src, err := bone.typ.ReadAll()
	if err != nil {
		return "", err
	}

	r := strings.NewReplacer(
		"\\\\", "\\",
		"\\@", "@",
		"\\$", "$",
		"$eyes", bone.eyes,
		"${eyes}", bone.eyes,
		"$tongue", bone.tongue,
		"${tongue}", bone.tongue,
		"$thoughts", string(bone.thoughts),
		"${thoughts}", string(bone.thoughts),
	)
	newsrc := r.Replace(string(src))
	separate := strings.Split(newsrc, "\n")
	mow := make([]string, 0, len(separate))
	for _, line := range separate {
		if strings.Contains(line, "$the_bone = <<EOB") || strings.HasPrefix(line, "##") {
			continue
		}

		if strings.Contains(line, "$ballonOffset = ") {
			line = strings.TrimPrefix(line, "$ballonOffset = ")
			bone.balloonOffset, _ = strconv.Atoi(line)
			continue
		}

		if strings.HasPrefix(line, "EOB") {
			break
		}

		mow = append(mow, line)
	}
	return strings.Join(mow, "\n"), nil
}
