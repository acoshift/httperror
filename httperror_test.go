package httperror

import (
	"errors"
	"testing"
)

func TestMerge(t *testing.T) {
	cases := []struct {
		err   error
		other error
		out   error
	}{
		{nil, nil, nil},
		{errors.New(""), nil, errors.New("")},
		{errors.New("a"), nil, errors.New("a")},
		{nil, errors.New(""), errors.New("")},
		{nil, errors.New("a"), errors.New("a")},
		{errors.New(""), errors.New(""), errors.New("; ")},
		{errors.New("a"), errors.New("b"), errors.New("a; b")},
		{BadRequest, nil, BadRequest},
		{nil, Conflict, Conflict},
		{emptyBadRequest, Conflict, BadRequestWith(Conflict)},
		{emptyConflict, BadRequest, ConflictWith(BadRequest)},
		{errors.New("a"), InternalServerError, errors.New(InternalServerError.Error() + "; a")},
		{emptyBadRequest, errors.New("invalid payload"), errors.New("bad_request: [400] invalid payload")},
	}

	for _, c := range cases {
		r := Merge(c.err, c.other)
		t.Log(r)
		if c.out != nil {
			if r.Error() != c.out.Error() {
				t.Errorf("expected merge result to be %v; got %v", c.out.Error(), r.Error())
			}
		} else {
			if r != nil {
				t.Errorf("expected merge result to be nil; got %v", r.Error())
			}
		}
	}
}
