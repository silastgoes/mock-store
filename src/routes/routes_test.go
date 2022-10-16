package routes

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/silastgoes/mock-store/src/controllers/mocks"
)

func TestLoadRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)

	srv := mocks.NewMockProductControlService(ctrl)
	rs := NewRouterService(srv)

	srv.EXPECT().Index(gomock.Any(), gomock.Any()).Return().AnyTimes()
	srv.EXPECT().New(gomock.Any(), gomock.Any()).Return().AnyTimes()
	srv.EXPECT().Insert(gomock.Any(), gomock.Any()).Return().AnyTimes()
	srv.EXPECT().Delete(gomock.Any(), gomock.Any()).Return().AnyTimes()
	srv.EXPECT().Edit(gomock.Any(), gomock.Any()).Return().AnyTimes()
	srv.EXPECT().Update(gomock.Any(), gomock.Any()).Return().AnyTimes()

	rs.LoadRoutes()
}
