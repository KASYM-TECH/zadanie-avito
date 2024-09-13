package server

import (
	"avito/controllers"
	"avito/log"
	"github.com/gorilla/mux"
	"net/http"
)

type Router struct {
	Router *mux.Router
	logger log.Logger
}

type Controllers struct {
	UserCnt   *controllers.UserController
	DummyCnt  *controllers.DummyController
	TenderCnt *controllers.TenderController
	OrgCnt    *controllers.OrganizationController
	BidCnt    *controllers.BidController
}

func NewRouter(logger log.Logger) *Router {
	return &Router{Router: mux.NewRouter(), logger: logger}
}

func (r *Router) AddRoutes(m *Middleware, cts Controllers) {
	register := func(path, method string, h http.HandlerFunc) {
		r.Router.HandleFunc(path, h).Methods(method)
	}

	register("/api/ping", "GET", m.Wrap(cts.DummyCnt.Ping))

	register("/api/auth/signup", "POST", m.Wrap(cts.UserCnt.Signup))

	register("/api/organizations/new", "POST", m.Wrap(cts.OrgCnt.Create))
	register("/api/organizations/bond", "POST", m.Wrap(cts.OrgCnt.MakeResponsible))

	register("/api/tenders/new", "POST", m.Wrap(cts.TenderCnt.Create))
	register("/api/tenders/{tenderId}/status", "GET", m.Wrap(cts.TenderCnt.GetStatus))
	register("/api/tenders/{tenderId}/status", "PUT", m.Wrap(cts.TenderCnt.SetStatus))
	register("/api/tenders/my", "GET", m.Wrap(cts.TenderCnt.GetByUsername))
	register("/api/tenders", "GET", m.Wrap(cts.TenderCnt.GetPublished))
	register("/api/tenders/{tenderId}/edit", "PATCH", m.Wrap(cts.TenderCnt.Edit))
	register("/api/tenders/{tenderId}/rollback/{version}", "PUT", m.Wrap(cts.TenderCnt.Rollback))

	register("/api/bids/new", "POST", m.Wrap(cts.BidCnt.Create))
	register("/api/bids/my", "GET", m.Wrap(cts.BidCnt.GetByUsername))
	register("/api/bids/{tenderId}/list", "GET", m.Wrap(cts.BidCnt.GetByTenderId))
	register("/api/bids/{bidId}/status", "GET", m.Wrap(cts.BidCnt.GetStatus))
	register("/api/bids/{bidId}/status", "PUT", m.Wrap(cts.BidCnt.SetStatus))
	register("/api/bids/{bidId}/edit", "PATCH", m.Wrap(cts.BidCnt.Edit))
	register("/api/bids/{bidId}/submit_decision", "PUT", m.Wrap(cts.BidCnt.SubmitDecision))
	register("/api/bids/{bidId}/feedback", "PUT", m.Wrap(cts.BidCnt.SubmitFeedback))
	register("/api/bids/{bidId}/rollback/{version}", "PUT", m.Wrap(cts.BidCnt.Rollback))
	register("/api/bids/{tenderId}/reviews", "GET", m.Wrap(cts.BidCnt.Reviews))
}
