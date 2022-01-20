package bonesay

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestBone_Clone(t *testing.T) {
	tests := []struct {
		name string
		opts []Option
		from *Bone
		want *Bone
	}{
		{
			name: "without options",
			opts: []Option{},
			from: func() *Bone {
				bone, _ := New()
				return bone
			}(),
			want: func() *Bone {
				bone, _ := New()
				return bone
			}(),
		},
		{
			name: "with some options",
			opts: []Option{},
			from: func() *Bone {
				bone, _ := New(
					Type("docker"),
					BallonWidth(60),
				)
				return bone
			}(),
			want: func() *Bone {
				bone, _ := New(
					Type("docker"),
					BallonWidth(60),
				)
				return bone
			}(),
		},
		{
			name: "clone and some options",
			opts: []Option{
				Thinking(),
				Thoughts('o'),
			},
			from: func() *Bone {
				bone, _ := New(
					Type("docker"),
					BallonWidth(60),
				)
				return bone
			}(),
			want: func() *Bone {
				bone, _ := New(
					Type("docker"),
					BallonWidth(60),
					Thinking(),
					Thoughts('o'),
				)
				return bone
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.want.Clone(tt.opts...)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.want, got,
				cmp.AllowUnexported(Bone{}),
				cmpopts.IgnoreFields(Bone{}, "buf")); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}

	t.Run("random", func(t *testing.T) {
		bone, _ := New(
			Type(""),
			Thinking(),
			Thoughts('o'),
			Eyes("xx"),
			Tongue("u"),
			Random(),
		)

		cloned, _ := bone.Clone()

		if diff := cmp.Diff(bone, cloned,
			cmp.AllowUnexported(Bone{}),
			cmpopts.IgnoreFields(Bone{}, "buf")); diff != "" {
			t.Errorf("(-want, +got)\n%s", diff)
		}
	})

	t.Run("error", func(t *testing.T) {
		bone, err := New()
		if err != nil {
			t.Fatal(err)
		}

		wantErr := errors.New("error")
		_, err = bone.Clone(func(*Bone) error {
			return wantErr
		})
		if wantErr != err {
			t.Fatalf("want %v, but got %v", wantErr, err)
		}
	})
}

func Test_adjustTo2Chars(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "empty",
			s:    "",
			want: "  ",
		},
		{
			name: "1 character",
			s:    "1",
			want: "1 ",
		},
		{
			name: "2 characters",
			s:    "12",
			want: "12",
		},
		{
			name: "3 characters",
			s:    "123",
			want: "12",
		},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			if got := adjustTo2Chars(tt.s); got != tt.want {
				t.Errorf("adjustTo2Chars() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotFound_Error(t *testing.T) {
	file := "test"
	n := &NotFound{
		Bonefile: file,
	}
	want := fmt.Sprintf("not found %q bonefile", file)
	if want != n.Error() {
		t.Fatalf("want %q but got %q", want, n.Error())
	}
}
