package system

import (
	"context"
	"testing"

	"github.com/smartystreets/assertions"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/SailGame/Core/data/memory"
	"github.com/SailGame/Core/pb/core"
	"github.com/SailGame/Core/server"
)

func TestAccountLogin(t *testing.T) {
	Convey("Login Successfully", t, func() {
		f := newFixture()
		f.init(&server.CoreServerConfig{
			MStorage: memory.NewStorage(),
		})
		uc := f.newUserClient()
		userName := "test"
		ret, err := uc.mCoreClient.Login(context.TODO(), &core.LoginArgs{
			UserName: userName,
			Password: "",
		})
		So(err, assertions.ShouldBeNil)
		So(ret.Err, assertions.ShouldEqual, core.ErrorNumber_OK)
		So(ret.Token, assertions.ShouldNotBeBlank)
		So(ret.GetAccount().GetUserName(), assertions.ShouldEqual, userName)
	})
}
