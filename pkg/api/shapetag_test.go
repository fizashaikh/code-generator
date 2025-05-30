package api_test

import (
	"testing"

	"github.com/aws-controllers-k8s/code-generator/pkg/api"
)

func TestShapeTagJoin(t *testing.T) {
	s := api.ShapeTags{
		{Key: "location", Val: "query"},
		{Key: "locationName", Val: "abc"},
		{Key: "type", Val: "string"},
	}

	expected := `location:"query" locationName:"abc" type:"string"`

	o := s.Join(" ")
	o2 := s.String()
	if expected != o {
		t.Errorf("Expected %s, but received %s", expected, o)
	}
	if expected != o2 {
		t.Errorf("Expected %s, but received %s", expected, o2)
	}
}
