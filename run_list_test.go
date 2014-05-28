package chef

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var (
	runList RunList = []string{"recipe[foo]", "recipe[baz]", "role[banana]"}
)

func TestNodeRunList(t *testing.T) {
	rl := runList

	Convey("Node.RunList() should be a RunList", t, func() {
		So(rl, ShouldHaveSameTypeAs, RunList{})
	})

	Convey("Node.RunList() should be populated", t, func() {
		So(rl, ShouldContain, "recipe[foo]")
		So(rl, ShouldContain, "recipe[baz]")
		So(rl, ShouldContain, "role[banana]")
	})

	rl = RunList{}
	Convey("Empty RunList should be valid", t, func() {
		So(rl, ShouldBeEmpty)
	})

}
