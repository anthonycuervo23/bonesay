package bonesay

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestBones(t *testing.T) {
	t.Run("no set BONEPATH env", func(t *testing.T) {
		bonePaths, err := Bones()
		if err != nil {
			t.Fatal(err)
		}
		if len(bonePaths) != 1 {
			t.Fatalf("want 1, but got %d", len(bonePaths))
		}
		bonePath := bonePaths[0]
		if len(bonePath.BoneFiles) == 0 {
			t.Fatalf("no bonefiles")
		}

		wantBonePath := &BonePath{
			Name:         "bones",
			LocationType: InBinary,
		}
		if diff := cmp.Diff(wantBonePath, bonePath,
			cmpopts.IgnoreFields(BonePath{}, "BoneFiles"),
		); diff != "" {
			t.Errorf("(-want, +got)\n%s", diff)
		}
	})

	t.Run("set BONEPATH env", func(t *testing.T) {
		bonepath := filepath.Join("testdata", "testdir")

		os.Setenv("BONEPATH", bonepath)
		defer os.Unsetenv("BONEPATH")

		bonePaths, err := Bones()
		if err != nil {
			t.Fatal(err)
		}
		if len(bonePaths) != 2 {
			t.Fatalf("want 2, but got %d", len(bonePaths))
		}

		wants := []*BonePath{
			{
				Name:         "testdata/testdir",
				LocationType: InDirectory,
			},
			{
				Name:         "bones",
				LocationType: InBinary,
			},
		}
		if diff := cmp.Diff(wants, bonePaths,
			cmpopts.IgnoreFields(BonePath{}, "BoneFiles"),
		); diff != "" {
			t.Errorf("(-want, +got)\n%s", diff)
		}

		if len(bonePaths[0].BoneFiles) != 1 {
			t.Fatalf("unexpected bonefiles len = %d, %+v",
				len(bonePaths[0].BoneFiles), bonePaths[0].BoneFiles,
			)
		}

		if bonePaths[0].BoneFiles[0] != "test" {
			t.Fatalf("want %q but got %q", "test", bonePaths[0].BoneFiles[0])
		}
	})

	t.Run("set BONEPATH env", func(t *testing.T) {
		os.Setenv("BONEPATH", "notfound")
		defer os.Unsetenv("BONEPATH")

		_, err := Bones()
		if err == nil {
			t.Fatal("want error")
		}
	})

}

func TestBonePath_Lookup(t *testing.T) {
	t.Run("looked for bonefile", func(t *testing.T) {
		c := &BonePath{
			Name:         "basepath",
			BoneFiles:    []string{"test"},
			LocationType: InBinary,
		}
		got, ok := c.Lookup("test")
		if !ok {
			t.Errorf("want %v", ok)
		}
		want := &BoneFile{
			Name:         "test",
			BasePath:     "basepath",
			LocationType: InBinary,
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("(-want, +got)\n%s", diff)
		}
	})

	t.Run("no bonefile", func(t *testing.T) {
		c := &BonePath{
			Name:         "basepath",
			BoneFiles:    []string{"test"},
			LocationType: InBinary,
		}
		got, ok := c.Lookup("no bonefile")
		if ok {
			t.Errorf("want %v", !ok)
		}
		if got != nil {
			t.Error("want nil")
		}
	})
}

func TestBoneFile_ReadAll(t *testing.T) {
	fromTestData := &BoneFile{
		Name:         "test",
		BasePath:     filepath.Join("testdata", "testdir"),
		LocationType: InDirectory,
	}
	fromTestdataContent, err := fromTestData.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	fromBinary := &BoneFile{
		Name:         "default",
		BasePath:     "bones",
		LocationType: InBinary,
	}
	fromBinaryContent, err := fromBinary.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(fromTestdataContent, fromBinaryContent) {
		t.Fatalf("testdata\n%s\n\nbinary%s\n", string(fromTestdataContent), string(fromBinaryContent))
	}

}

const defaultSay = ` ________ 
< bonesay >
 -------- 
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||`

func TestSay(t *testing.T) {
	type args struct {
		phrase  string
		options []Option
	}
	tests := []struct {
		name     string
		args     args
		wantFile string
		wantErr  bool
	}{
		{
			name: "default",
			args: args{
				phrase: "hello!",
			},
			wantFile: "default.bone",
			wantErr:  false,
		},
		{
			name: "nest",
			args: args{
				phrase: defaultSay,
				options: []Option{
					DisableWordWrap(),
				},
			},
			wantFile: "nest.bone",
			wantErr:  false,
		},
		{
			name: "error",
			args: args{
				phrase: "error",
				options: []Option{
					func(*Bone) error {
						return errors.New("error")
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Say(tt.args.phrase, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Say() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			filename := filepath.Join("testdata", tt.wantFile)
			content, err := ioutil.ReadFile(filename)
			if err != nil {
				t.Fatal(err)
			}
			want := string(content)
			if want != got {
				t.Fatalf("want\n%s\n\ngot\n%s", want, got)
			}
		})
	}
}
