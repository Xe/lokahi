package integration

import (
	"errors"
	"fmt"
	"testing"
)

// har har
func TestSuite(t *testing.T) {
	s := &Suite{}

	t.Run("basic setup/teardown", func(t *testing.T) {
		err := s.Setup()
		if err != nil {
			t.Fatal(err)
		}

		err = s.Teardown()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("err tools", func(t *testing.T) {
		for _, cs := range []error{nil, errors.New("hi there")} {
			t.Run(fmt.Sprint(cs), func(t *testing.T) {
				s.SetErr(cs)

				err := s.GetErr()
				if err != cs {
					t.Fatalf("expected s.err == %v, got: %v", cs, err)
				}

				err = s.WantAnError()
				if err == nil && cs == nil {
					t.Fatal("expected s.WantAnError to return an error as s.err is nil")
				}
			})
		}
	})
}
