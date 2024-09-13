package basic

import (
	"avito/db/model"
	"avito/domain"
	"context"
	"github.com/txix-open/isp-kit/http/httpcli"
	"net/http"
)

type EmployeeOrg struct {
	OrgId      string
	EmployeeId string
	Username   string
}

func CreateOrgEmployee(test *Test, username string) EmployeeOrg {
	userId, userResp := CreateUser(test, domain.SignupRequest{
		Username:  username,
		FirstName: "first",
		LastName:  "last",
	})
	test.Assertions.Equal(http.StatusOK, userResp.StatusCode())

	orgId, orgResp := CreateOrganization(test)
	test.Assertions.Equal(http.StatusOK, orgResp.StatusCode())

	_, bondResp := Bond(test, userId, orgId)
	test.Assertions.Equal(http.StatusOK, bondResp.StatusCode())

	return EmployeeOrg{
		EmployeeId: userId,
		Username:   username,
		OrgId:      orgId,
	}
}

func Bond(test *Test, userId, orgId string) (string, *httpcli.Response) {
	assert := test.Assertions

	req := domain.BondReq{
		UserId:         userId,
		OrganizationId: orgId,
	}

	var bondId string
	resp, err := test.Cli.Post(test.URL + "/api/organizations/bond").
		JsonRequestBody(&req).
		JsonResponseBody(&bondId).
		Do(context.Background())

	assert.NoError(err)

	return bondId, resp
}

func CreateUser(test *Test, req domain.SignupRequest) (string, *httpcli.Response) {
	assert := test.Assertions

	var userId string
	resp, err := test.Cli.Post(test.URL + "/api/auth/signup").
		JsonRequestBody(&req).
		JsonResponseBody(&userId).
		Do(context.Background())

	assert.NoError(err)

	return userId, resp
}

func CreateOrganization(test *Test) (string, *httpcli.Response) {
	assert := test.Assertions

	req := domain.CreateOrganizationReq{
		Name:        "test_org_name",
		Description: "org",
		Type:        model.OrganizationTypeJSC,
	}

	var userId string
	resp, err := test.Cli.Post(test.URL + "/api/organizations/new").
		JsonRequestBody(&req).
		JsonResponseBody(&userId).
		Do(context.Background())

	assert.NoError(err)

	return userId, resp
}
