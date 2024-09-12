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
	register("/api/auth/login", "POST", m.Wrap(cts.UserCnt.Login))

	register("/api/organizations/new", "POST", m.Wrap(cts.OrgCnt.Create))
	register("/api/organizations/bond", "POST", m.Wrap(cts.OrgCnt.MakeResponsible))

	register("/api/tenders/new", "POST", m.WrapAuthed(cts.TenderCnt.Create))
	register("/api/tenders/{tenderId}/status", "GET", m.WrapAuthed(cts.TenderCnt.GetStatus))
	register("/api/tenders/{tenderId}/status", "PUT", m.WrapAuthed(cts.TenderCnt.SetStatus))
	register("/api/tenders/my", "GET", m.WrapAuthed(cts.TenderCnt.GetByUsername))
	register("/api/tenders", "GET", m.WrapAuthed(cts.TenderCnt.GetPublished))
	register("/api/tenders/{tenderId}/edit", "PATCH", m.WrapAuthed(cts.TenderCnt.Edit))
	register("/api/tenders/{tenderId}/rollback/{version}", "PUT", m.WrapAuthed(cts.TenderCnt.Rollback))

	register("/api/bids/new", "POST", m.WrapAuthed(cts.BidCnt.Create))
	register("/api/bids/my", "GET", m.WrapAuthed(cts.BidCnt.GetByUsername))
	register("/api/bids/{tenderId}/list", "GET", m.WrapAuthed(cts.BidCnt.GetByTenderId))
	register("/api/bids/{bidId}/status", "GET", m.WrapAuthed(cts.BidCnt.GetStatus))
	register("/api/bids/{bidId}/status", "PUT", m.WrapAuthed(cts.BidCnt.SetStatus))
	register("/api/bids/{bidId}/edit", "PATCH", m.WrapAuthed(cts.BidCnt.Edit))
	register("/api/bids/{bidId}/submit_decision", "PUT", m.WrapAuthed(cts.BidCnt.SubmitDecision))
	register("/api/bids/{bidId}/feedback", "PUT", m.WrapAuthed(cts.BidCnt.SubmitFeedback))
	register("/api/bids/{bidId}/rollback/{version}", "PUT", m.WrapAuthed(cts.BidCnt.Rollback))
	register("/api/bids/{tenderId}/reviews", "GET", m.WrapAuthed(cts.BidCnt.Reviews))
}
