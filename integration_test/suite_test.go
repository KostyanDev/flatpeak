package integration_test

import (
	"app/internal/mocks"
	"app/internal/service"
	httpClient "app/internal/transport/http"

	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"testing"
)

type TestSuite struct {
	suite.Suite
	server     *httptest.Server
	handler    *httpClient.Handler
	router     *mux.Router
	mockClient *mocks.CarbonClient
	service    *service.Service
}

func (s *TestSuite) SetupSuite() {
	//cfg, err := config.New[config.Config]()
	//s.Require().NoError(err, "Failed to load config")
	log := logrus.New()

	s.mockClient = &mocks.CarbonClient{}
	s.service = service.New(context.Background(), log, s.mockClient)
	s.handler = httpClient.New(context.Background(), log, s.service)

	s.router = mux.NewRouter()
	httpClient.RegisterRoutes(s.router, s.handler)
	s.server = httptest.NewServer(s.router)
}

func (s *TestSuite) TearDownSuite() {
	s.server.Close()
}

// 🔹 Run the test suite
func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
