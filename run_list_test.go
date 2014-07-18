package chef_test

import (
	"github.com/go-chef/chef"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var (
	rl = chef.RunList{"recipe[foo]", "recipe[baz]", "role[banana]"}
)

func TestNodeRunList(t *testing.T) {
	Convey("Node.RunList() should be a RunList", t, func() {
		So(rl, ShouldHaveSameTypeAs, chef.RunList{})
	})

	Convey("Node.RunList() should be populated", t, func() {
		So(rl, ShouldContain, "recipe[foo]")
		So(rl, ShouldContain, "recipe[baz]")
		So(rl, ShouldContain, "role[banana]")
	})

	rl = chef.RunList{}
	Convey("Empty RunList should be valid", t, func() {
		So(rl, ShouldBeEmpty)
	})

}
