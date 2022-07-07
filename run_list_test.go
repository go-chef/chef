package chef

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	rl = RunList{"recipe[foo]", "recipe[baz]", "role[banana]"}
)

func TestNodeRunList(t *testing.T) {
	assert.IsType(t, RunList{}, rl, "Runlist type")
	assert.Contains(t, rl, "recipe[foo]", "Runlist contents")
	assert.Contains(t, rl, "recipe[baz]", "Runlist contents")
	assert.Contains(t, rl, "role[banana]", "Runlist contents")

	rl = RunList{}
	assert.Empty(t, rl, "Empty runlist")
}
