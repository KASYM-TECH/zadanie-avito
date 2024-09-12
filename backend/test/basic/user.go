package basic

import (
	"avito/db/model"
	"avito/domain"
	"context"
	"github.com/txix-open/isp-kit/http/httpcli"
	"net/http"
)

type EmployeeOrg struct {
	Token      string
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

	tokens, loginResp := LoginUser(test, username)
	test.Assertions.Equal(http.StatusOK, loginResp.StatusCode())

	orgId, orgResp := CreateOrganization(test, tokens.AccessToken)
	test.Assertions.Equal(http.StatusOK, orgResp.StatusCode())

	_, bondResp := Bond(test, userId, orgId, tokens.AccessToken)
	test.Assertions.Equal(http.StatusOK, bondResp.StatusCode())

	return EmployeeOrg{
		EmployeeId: userId,
		Username:   username,
		OrgId:      orgId,
		Token:      tokens.AccessToken,
	}
}

func Bond(test *Test, userId, orgId, token string) (string, *httpcli.Response) {
	assert := test.Assertions

	req := domain.BondReq{
		UserId:         userId,
		OrganizationId: orgId,
	}

	var bondId string
	resp, err := test.Cli.Post(test.URL+"/api/organizations/bond").
		JsonRequestBody(&req).
		JsonResponseBody(&bondId).
		Header("Authorization", token).
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

func LoginUser(test *Test, username string) (domain.LoginResponse, *httpcli.Response) {
	assert := test.Assertions

	req := domain.LoginRequest{
		Username: username,
	}

	var tokens domain.LoginResponse
	resp, err := test.Cli.Post(test.URL + "/api/auth/login").
		JsonRequestBody(&req).
		JsonResponseBody(&tokens).
		Do(context.Background())

	assert.NoError(err)

	return tokens, resp
}

func CreateOrganization(test *Test, token string) (string, *httpcli.Response) {
	assert := test.Assertions

	req := domain.CreateOrganizationReq{
		Name:        "test_org_name",
		Description: "org",
		Type:        model.OrganizationTypeJSC,
	}

	var userId string
	resp, err := test.Cli.Post(test.URL+"/api/organizations/new").
		JsonRequestBody(&req).
		JsonResponseBody(&userId).
		Header("Authorization", token).
		Do(context.Background())

	assert.NoError(err)

	return userId, resp
}
