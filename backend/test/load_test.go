package test

import (
	"avito/db/model"
	"avito/domain"
	"avito/test/basic"
	"fmt"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"net/http"
	"testing"
	"time"
)

func TestLoadStress(t *testing.T) {
	test := basic.InitTest(t)

	martinOrg := basic.CreateOrgEmployee(test, "Martin")

	tenderReq := domain.CreateTenderReq{
		Name:            "n1",
		Description:     "d1",
		ServiceType:     model.TenderServiceTypeConstruction,
		Status:          model.TenderStatusCreated,
		OrganizationId:  martinOrg.OrgId,
		CreatorUsername: martinOrg.Username,
	}
	tender, tenderResp := basic.CreateTender(test, tenderReq)
	test.Assertions.Equal(http.StatusOK, tenderResp.StatusCode())

	tenders, myResp := basic.GetTenderByUsername(test, martinOrg.Username, 0, 3)
	test.Assertions.Equal(http.StatusOK, myResp.StatusCode())
	test.Assertions.Len(tenders, 1)

	rate := vegeta.Rate{Freq: 1000, Per: time.Second}
	duration := 4 * time.Second
	pinpointer := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    test.URL + "/api/tenders/" + tender.Id + "/status",
	})
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(pinpointer, rate, duration, "LOAD TESTING") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
	test.Assertions.Less(metrics.Latencies.P99, time.Millisecond*50)
}
