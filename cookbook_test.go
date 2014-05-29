package chef

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewCookbook(t *testing.T) {
	Convey("Create a cookbook", t, func() {
		cook, err := NewCookbook(&Reader{"name": "blah"})
		So(cook.Name, ShouldEqual, "blah")
		So(err, ShouldBeNil)
	})
	Convey("Cookbook Decoder fail", t, func() {
		failCook, err := NewCookbook(&Reader{"name": struct{}{}})
		So(err, ShouldNotBeNil)
		So(failCook, ShouldBeNil)
	})

}
