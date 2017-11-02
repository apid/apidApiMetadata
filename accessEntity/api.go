// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package accessEntity

import (
	"encoding/json"
	"fmt"
	"github.com/apid/apidVerifyApiKey/common"
	"net/http"
	"strconv"
	"strings"
)

const (
	AccessEntityPath         = "/entities"
	EndpointApp              = "/apps"
	EndpointApiProduct       = "/apiproducts"
	EndpointCompany          = "/companies"
	EndpointCompanyDeveloper = "/companydevelopers"
	EndpointDeveloper        = "/developers"
	EndpointAppCredentials   = "/appcredentials"
)

const (
	IdentifierAppId          = "appid"
	IdentifierApiProductName = "apiproductname"
	IdentifierAppName        = "appname"
	IdentifierApiResource    = "apiresource"
	IdentifierDeveloperId    = "developerid"
	IdentifierDeveloperEmail = "developeremail"
	IdentifierConsumerKey    = "consumerkey"
	IdentifierCompanyName    = "companyname"
)

const (
	AppTypeDeveloper = "DEVELOPER"
	AppTypeCompany   = "COMPANY"
)

var (
	Identifiers = map[string]bool{
		"appid":          true,
		"apiproductname": true,
		"appname":        true,
		"apiresource":    true,
		"developerid":    true,
		"companyname":    true,
		"developeremail": true,
		"consumerkey":    true,
	}

	ErrInvalidPar = &common.ErrorResponse{
		ResponseCode:    strconv.Itoa(INVALID_PARAMETERS),
		ResponseMessage: "invalid identifiers",
		StatusCode:      http.StatusBadRequest,
	}

	IdentifierTree = map[string]map[string][]string{
		EndpointApiProduct: {
			IdentifierApiProductName: {},
			IdentifierAppId:          {IdentifierApiResource},
			IdentifierAppName:        {IdentifierApiResource, IdentifierDeveloperEmail, IdentifierDeveloperId, IdentifierCompanyName},
			IdentifierConsumerKey:    {IdentifierApiResource},
		},
		EndpointApp: {
			IdentifierAppId:       {},
			IdentifierAppName:     {IdentifierDeveloperEmail, IdentifierDeveloperId, IdentifierCompanyName},
			IdentifierConsumerKey: {},
		},
		EndpointCompany: {
			IdentifierAppId:       {},
			IdentifierCompanyName: {},
			IdentifierConsumerKey: {},
		},
		EndpointCompanyDeveloper: {
			IdentifierCompanyName: {},
		},
		EndpointAppCredentials: {
			IdentifierConsumerKey: {},
		},
		EndpointDeveloper: {
			IdentifierDeveloperEmail: {},
			IdentifierAppId:          {},
			IdentifierDeveloperId:    {},
			IdentifierConsumerKey:    {},
		},
	}
)

const (
	INVALID_PARAMETERS = iota
	DB_ERROR
	DATA_ERROR
)

type ApiManager struct {
	DbMan            DbManagerInterface
	AccessEntityPath string
	apiInitialized   bool
}

func (a *ApiManager) InitAPI() {
	if a.apiInitialized {
		return
	}
	services.API().HandleFunc(a.AccessEntityPath+EndpointApp, a.HandleApps).Methods("GET")
	services.API().HandleFunc(a.AccessEntityPath+EndpointApiProduct, a.HandleApiProducts).Methods("GET")
	services.API().HandleFunc(a.AccessEntityPath+EndpointCompany, a.HandleCompanies).Methods("GET")
	services.API().HandleFunc(a.AccessEntityPath+EndpointCompanyDeveloper, a.HandleCompanyDevelopers).Methods("GET")
	services.API().HandleFunc(a.AccessEntityPath+EndpointDeveloper, a.HandleDevelopers).Methods("GET")
	services.API().HandleFunc(a.AccessEntityPath+EndpointAppCredentials, a.HandleAppCredentials).Methods("GET")
	a.apiInitialized = true
	log.Debug("API endpoints initialized")
}

func (a *ApiManager) HandleApps(w http.ResponseWriter, r *http.Request) {
	ids, org, err := extractIdentifiers(r.URL.Query())
	if err != nil {
		common.WriteError(w, err.Error(), INVALID_PARAMETERS, http.StatusBadRequest)
	}
	res, errRes := a.getApp(org, ids)
	if errRes != nil {
		writeJson(errRes, w, r)
		return
	}
	writeJson(res, w, r)
}

func (a *ApiManager) HandleApiProducts(w http.ResponseWriter, r *http.Request) {
	ids, org, err := extractIdentifiers(r.URL.Query())
	if err != nil {
		common.WriteError(w, err.Error(), INVALID_PARAMETERS, http.StatusBadRequest)
	}
	details, errRes := a.getApiProduct(org, ids)
	if errRes != nil {
		w.WriteHeader(errRes.StatusCode)
		writeJson(errRes, w, r)
		return
	}
	writeJson(details, w, r)
}

func (a *ApiManager) HandleCompanies(w http.ResponseWriter, r *http.Request) {
	ids, org, err := extractIdentifiers(r.URL.Query())
	if err != nil {
		common.WriteError(w, err.Error(), INVALID_PARAMETERS, http.StatusBadRequest)
	}
	details, errRes := a.getCompany(org, ids)
	if errRes != nil {
		w.WriteHeader(errRes.StatusCode)
		writeJson(errRes, w, r)
		return
	}
	writeJson(details, w, r)
}

func (a *ApiManager) HandleCompanyDevelopers(w http.ResponseWriter, r *http.Request) {
	ids, org, err := extractIdentifiers(r.URL.Query())
	if err != nil {
		common.WriteError(w, err.Error(), INVALID_PARAMETERS, http.StatusBadRequest)
	}
	details, errRes := a.getCompanyDeveloper(org, ids)
	if errRes != nil {
		w.WriteHeader(errRes.StatusCode)
		writeJson(errRes, w, r)
		return
	}
	writeJson(details, w, r)
}
func (a *ApiManager) HandleDevelopers(w http.ResponseWriter, r *http.Request) {
	ids, org, err := extractIdentifiers(r.URL.Query())
	if err != nil {
		common.WriteError(w, err.Error(), INVALID_PARAMETERS, http.StatusBadRequest)
	}
	details, errRes := a.getDeveloper(org, ids)
	if errRes != nil {
		w.WriteHeader(errRes.StatusCode)
		writeJson(errRes, w, r)
		return
	}
	writeJson(details, w, r)
}

func (a *ApiManager) HandleAppCredentials(w http.ResponseWriter, r *http.Request) {
	ids, org, err := extractIdentifiers(r.URL.Query())
	if err != nil {
		common.WriteError(w, err.Error(), INVALID_PARAMETERS, http.StatusBadRequest)
	}
	details, errRes := a.getAppCredential(org, ids)
	if errRes != nil {
		w.WriteHeader(errRes.StatusCode)
		writeJson(errRes, w, r)
		return
	}
	writeJson(details, w, r)
}

func extractIdentifiers(pars map[string][]string) (map[string]string, string, error) {
	m := make(map[string]string)
	org := ""
	if orgs := pars["organization"]; len(orgs) > 0 {
		org = orgs[0]
	}
	for k, v := range pars {
		k = strings.ToLower(k)
		if Identifiers[k] {
			if len(v) == 1 {
				m[k] = v[0]
			} else {
				return nil, org, fmt.Errorf("each identifier must have only 1 value")
			}
		}
	}
	return m, org, nil
}

func (a *ApiManager) getCompanyDeveloper(org string, ids map[string]string) (*CompanyDevelopersSuccessResponse, *common.ErrorResponse) {
	valid, keyVals := parseIdentifiers(EndpointDeveloper, ids)
	if !valid {
		return nil, ErrInvalidPar
	}
	priKey, priVal, secKey, secVal := keyVals[0], keyVals[1], keyVals[2], keyVals[3]

	devs, err := a.DbMan.GetCompanyDevelopers(priKey, priVal, secKey, secVal)
	if err != nil {
		log.Errorf("getCompanyDeveloper: %v", err)
		return nil, newDbError(err)
	}

	if len(devs) == 0 {
		return &CompanyDevelopersSuccessResponse{
			CompanyDevelopers: nil,
			Organization:      org,
		}, nil
	}

	var details []*CompanyDeveloperDetails
	for _, dev := range devs {
		comName, err := a.DbMan.GetComNameByComId(dev.CompanyId)
		if err != nil {
			log.Errorf("getCompanyDeveloper: %v", err)
			return nil, newDbError(err)
		}
		email, err := a.DbMan.GetDevEmailByDevId(dev.DeveloperId)
		if err != nil {
			log.Errorf("getCompanyDeveloper: %v", err)
			return nil, newDbError(err)
		}
		detail := makeComDevDetails(&dev, comName, email, priKey, priVal)
		details = append(details, detail)
	}
	return &CompanyDevelopersSuccessResponse{
		CompanyDevelopers: details,
		Organization:      org,
	}, nil
}

func (a *ApiManager) getDeveloper(org string, ids map[string]string) (*DeveloperSuccessResponse, *common.ErrorResponse) {
	valid, keyVals := parseIdentifiers(EndpointDeveloper, ids)
	if !valid {
		return nil, ErrInvalidPar
	}
	priKey, priVal, secKey, secVal := keyVals[0], keyVals[1], keyVals[2], keyVals[3]

	devs, err := a.DbMan.GetDevelopers(priKey, priVal, secKey, secVal)
	if err != nil {
		log.Errorf("getDeveloper: %v", err)
		return nil, newDbError(err)
	}

	if len(devs) == 0 {
		return &DeveloperSuccessResponse{
			Developer:    nil,
			Organization: org,
		}, nil
	}
	dev := &devs[0]

	attrs := a.DbMan.GetKmsAttributes(dev.TenantId, dev.Id)[dev.Id]
	comNames, err := a.DbMan.GetComNamesByDevId(dev.Id)
	if err != nil {
		log.Errorf("getDeveloper: %v", err)
		return nil, newDbError(err)
	}
	appNames, err := a.DbMan.GetAppNamesByDevId(dev.Id)
	if err != nil {
		log.Errorf("getDeveloper: %v", err)
		return nil, newDbError(err)
	}
	details := makeDevDetails(dev, appNames, comNames, attrs, priKey, priVal)
	return &DeveloperSuccessResponse{
		Developer:    details,
		Organization: org,
	}, nil
}

func (a *ApiManager) getCompany(org string, ids map[string]string) (*CompanySuccessResponse, *common.ErrorResponse) {
	valid, keyVals := parseIdentifiers(EndpointCompany, ids)
	if !valid {
		return nil, ErrInvalidPar
	}
	priKey, priVal, secKey, secVal := keyVals[0], keyVals[1], keyVals[2], keyVals[3]

	coms, err := a.DbMan.GetCompanies(priKey, priVal, secKey, secVal)
	if err != nil {
		log.Errorf("getCompany: %v", err)
		return nil, newDbError(err)
	}

	if len(coms) == 0 {
		return &CompanySuccessResponse{
			Company:      nil,
			Organization: org,
		}, nil
	}
	com := &coms[0]

	attrs := a.DbMan.GetKmsAttributes(com.TenantId, com.Id)[com.Id]
	appNames, err := a.DbMan.GetAppNamesByComId(com.Id)
	if err != nil {
		log.Errorf("getCompany: %v", err)
		return nil, newDbError(err)
	}
	details := makeCompanyDetails(com, appNames, attrs, priKey, priVal)
	return &CompanySuccessResponse{
		Company:      details,
		Organization: org,
	}, nil
}

func (a *ApiManager) getApiProduct(org string, ids map[string]string) (*ApiProductSuccessResponse, *common.ErrorResponse) {
	valid, keyVals := parseIdentifiers(EndpointApiProduct, ids)
	if !valid {
		return nil, ErrInvalidPar
	}
	priKey, priVal, secKey, secVal := keyVals[0], keyVals[1], keyVals[2], keyVals[3]

	prods, err := a.DbMan.GetApiProducts(priKey, priVal, secKey, secVal)
	if err != nil {
		log.Errorf("getApiProduct: %v", err)
		return nil, newDbError(err)
	}

	var attrs []common.Attribute
	if len(prods) == 0 {
		return &ApiProductSuccessResponse{
			ApiProduct:   nil,
			Organization: org,
		}, nil
	}
	prod := &prods[0]
	attrs = a.DbMan.GetKmsAttributes(prod.TenantId, prod.Id)[prod.Id]
	details, errRes := makeApiProductDetails(prod, attrs, priKey, priVal, secKey, secVal)
	if errRes != nil {
		return nil, errRes
	}

	return &ApiProductSuccessResponse{
		ApiProduct:   details,
		Organization: org,
	}, nil
}

func (a *ApiManager) getAppCredential(org string, ids map[string]string) (*AppCredentialSuccessResponse, *common.ErrorResponse) {
	valid, keyVals := parseIdentifiers(EndpointApiProduct, ids)
	if !valid {
		return nil, ErrInvalidPar
	}
	priKey, priVal, secKey, secVal := keyVals[0], keyVals[1], keyVals[2], keyVals[3]

	appCreds, err := a.DbMan.GetAppCredentials(priKey, priVal, secKey, secVal)
	if err != nil {
		log.Errorf("getAppCredential: %v", err)
		return nil, newDbError(err)
	}

	if len(appCreds) == 0 {
		return &AppCredentialSuccessResponse{
			AppCredential: nil,
			Organization:  org,
		}, nil
	}
	appCred := &appCreds[0]
	attrs := a.DbMan.GetKmsAttributes(appCred.TenantId, appCred.Id)[appCred.Id]
	apps, err := a.DbMan.GetApps(IdentifierAppId, appCred.AppId, "", "")
	if err != nil {
		log.Errorf("getAppCredential: %v", err)
		return nil, newDbError(err)
	}

	if len(apps) == 0 {
		log.Errorf("getAppCredential: No App with id=%v", appCred.AppId)
		return &AppCredentialSuccessResponse{
			AppCredential: nil,
			Organization:  org,
		}, nil
	}
	app := &apps[0]
	cd := a.getCredDetails(appCred, app.Status)
	devStatus := ""
	if app.DeveloperId != "" {
		devStatus, err = a.DbMan.GetStatus(app.DeveloperId, AppTypeDeveloper)
		if err != nil {
			log.Errorf("getAppCredential error get status: %v", err)
			return nil, newDbError(err)
		}
	}
	//TODO: isValidKey
	cks := makeConsumerKeyStatusDetails(app, cd, devStatus, "")
	//TODO: redirectUris
	details := makeAppCredentialDetails(appCred, cks, []string{}, attrs, priKey, priVal)
	return &AppCredentialSuccessResponse{
		AppCredential: details,
		Organization:  org,
	}, nil
}

func (a *ApiManager) getApp(org string, ids map[string]string) (*AppSuccessResponse, *common.ErrorResponse) {
	valid, keyVals := parseIdentifiers(EndpointApp, ids)
	if !valid {
		return nil, ErrInvalidPar
	}
	priKey, priVal, secKey, secVal := keyVals[0], keyVals[1], keyVals[2], keyVals[3]

	apps, err := a.DbMan.GetApps(priKey, priVal, secKey, secVal)
	if err != nil {
		log.Errorf("getApp: %v", err)
		return nil, newDbError(err)
	}

	var app *common.App
	var attrs []common.Attribute

	if len(apps) == 0 {
		return &AppSuccessResponse{
			App:          nil,
			Organization: org,
		}, nil
	}

	app = &apps[0]
	attrs = a.DbMan.GetKmsAttributes(app.TenantId, app.Id)[app.Id]
	prods, err := a.DbMan.GetApiProductNamesByAppId(app.Id)
	if err != nil {
		log.Errorf("getApp error getting productNames: %v", err)
		return nil, newDbError(err)
	}
	parStatus, err := a.DbMan.GetStatus(app.ParentId, app.Type)
	if err != nil {
		log.Errorf("getApp error getting parent status: %v", err)
		return nil, newDbError(err)
	}
	creds, err := a.DbMan.GetAppCredentials(IdentifierAppId, app.Id, "", "")
	if err != nil {
		log.Errorf("getApp error getting parent status: %v", err)
		return nil, newDbError(err)
	}
	var credDetails []*CredentialDetails
	for _, cred := range creds {
		credDetails = append(credDetails, a.getCredDetails(&cred, app.Status))
	}

	details, errRes := makeAppDetails(app, parStatus, prods, credDetails, attrs, priKey, priVal, secKey, secVal)
	if errRes != nil {
		return nil, errRes
	}
	return &AppSuccessResponse{
		App:          details,
		Organization: org,
	}, nil
}

func makeConsumerKeyStatusDetails(app *common.App, c *CredentialDetails, devStatus, isValidKey string) *ConsumerKeyStatusDetails {
	return &ConsumerKeyStatusDetails{
		AppCredential:   c,
		AppID:           c.AppID,
		AppName:         app.Name,
		AppStatus:       app.Status,
		AppType:         app.Type,
		DeveloperID:     app.DeveloperId,
		DeveloperStatus: devStatus,
		IsValidKey:      isValidKey,
	}
}

func makeAppCredentialDetails(ac *common.AppCredential, cks *ConsumerKeyStatusDetails, redirectUrl []string, attrs []common.Attribute, priKey, priVal string) *AppCredentialDetails {
	return &AppCredentialDetails{
		AppID:                  ac.AppId,
		AppName:                cks.AppName,
		Attributes:             attrs,
		ConsumerKey:            ac.Id,
		ConsumerKeyStatus:      cks,
		ConsumerSecret:         ac.ConsumerSecret,
		DeveloperID:            cks.DeveloperID,
		PrimaryIdentifierType:  priKey,
		PrimaryIdentifierValue: priVal,
		RedirectUris:           redirectUrl, //TODO
		Scopes:                 common.JsonToStringArray(ac.Scopes),
		Status:                 ac.Status,
	}
}

func makeApiProductDetails(prod *common.ApiProduct, attrs []common.Attribute, priKey, priVal, secKey, secVal string) (*ApiProductDetails, *common.ErrorResponse) {
	var a *ApiProductDetails
	if prod != nil {
		var quotaLimit int
		var err error
		if prod.Quota != "" {
			quotaLimit, err = strconv.Atoi(prod.Quota)
			if err != nil {
				return nil, newDataError(err)
			}
		}

		a = &ApiProductDetails{
			ApiProxies:               common.JsonToStringArray(prod.Proxies),
			ApiResources:             common.JsonToStringArray(prod.ApiResources),
			ApprovalType:             prod.ApprovalType,
			Attributes:               attrs,
			CreatedAt:                prod.CreatedAt,
			CreatedBy:                prod.CreatedBy,
			Description:              prod.Description,
			DisplayName:              prod.DisplayName,
			Environments:             common.JsonToStringArray(prod.Environments),
			ID:                       prod.Id,
			LastModifiedAt:           prod.UpdatedAt,
			LastModifiedBy:           prod.UpdatedBy,
			Name:                     prod.Name,
			QuotaInterval:            prod.QuotaInterval,
			QuotaLimit:               int64(quotaLimit),
			QuotaTimeUnit:            prod.QuotaTimeUnit,
			Scopes:                   common.JsonToStringArray(prod.Scopes),
			PrimaryIdentifierType:    priKey,
			PrimaryIdentifierValue:   priVal,
			SecondaryIdentifierType:  secKey,
			SecondaryIdentifierValue: secVal,
		}
	} else {
		a = new(ApiProductDetails)
	}
	return a, nil
}

func makeAppDetails(app *common.App, parentStatus string, prods []string, creds []*CredentialDetails, attrs []common.Attribute, priKey, priVal, secKey, secVal string) (*AppDetails, *common.ErrorResponse) {
	var a *AppDetails
	if app != nil {
		a = &AppDetails{
			AccessType:               app.AccessType,
			ApiProducts:              prods,
			AppCredentials:           creds,
			AppFamily:                app.AppFamily,
			AppParentID:              app.ParentId,
			AppParentStatus:          parentStatus,
			AppType:                  app.Type,
			Attributes:               attrs,
			CallbackUrl:              app.CallbackUrl,
			CreatedAt:                app.CreatedAt,
			CreatedBy:                app.CreatedBy,
			DisplayName:              app.DisplayName,
			Id:                       app.Id,
			KeyExpiresIn:             "", //TODO
			LastModifiedAt:           app.UpdatedAt,
			LastModifiedBy:           app.UpdatedBy,
			Name:                     app.Name,
			Scopes:                   []string{}, //TODO
			Status:                   app.Status,
			PrimaryIdentifierType:    priKey,
			PrimaryIdentifierValue:   priVal,
			SecondaryIdentifierType:  secKey,
			SecondaryIdentifierValue: secVal,
		}
	} else {
		a = new(AppDetails)
	}
	return a, nil
}

func makeCompanyDetails(com *common.Company, appNames []string, attrs []common.Attribute, priKey, priVal string) *CompanyDetails {
	return &CompanyDetails{
		Apps:           appNames,
		Attributes:     attrs,
		CreatedAt:      com.CreatedAt,
		CreatedBy:      com.CreatedBy,
		DisplayName:    com.DisplayName,
		ID:             com.Id,
		LastModifiedAt: com.UpdatedAt,
		LastModifiedBy: com.UpdatedBy,
		Name:           com.Name,
		PrimaryIdentifierType:  priKey,
		PrimaryIdentifierValue: priVal,
		Status:                 com.Status,
	}
}

func makeDevDetails(dev *common.Developer, appNames []string, comNames []string, attrs []common.Attribute, priKey, priVal string) *DeveloperDetails {
	return &DeveloperDetails{
		Apps:                   appNames,
		Attributes:             attrs,
		Companies:              comNames,
		CreatedAt:              dev.CreatedAt,
		CreatedBy:              dev.CreatedBy,
		Email:                  dev.Email,
		FirstName:              dev.FirstName,
		ID:                     dev.Id,
		LastModifiedAt:         dev.UpdatedAt,
		LastModifiedBy:         dev.UpdatedBy,
		LastName:               dev.LastName,
		Password:               dev.Password,
		PrimaryIdentifierType:  priKey,
		PrimaryIdentifierValue: priVal,
		Status:                 dev.Status,
		UserName:               dev.UserName,
	}
}

func makeComDevDetails(comDev *common.CompanyDeveloper, comName, devEmail, priKey, priVal string) *CompanyDeveloperDetails {
	return &CompanyDeveloperDetails{
		CompanyName:            comName,
		CreatedAt:              comDev.CreatedAt,
		CreatedBy:              comDev.CreatedBy,
		DeveloperEmail:         devEmail,
		LastModifiedAt:         comDev.UpdatedAt,
		LastModifiedBy:         comDev.UpdatedBy,
		PrimaryIdentifierType:  priKey,
		PrimaryIdentifierValue: priVal,
		Roles: common.JsonToStringArray(comDev.Roles),
	}
}

func (a *ApiManager) getCredDetails(cred *common.AppCredential, appStatus string) *CredentialDetails {

	return &CredentialDetails{
		ApiProductReferences: []string{}, //TODO
		AppID:                cred.AppId,
		AppStatus:            appStatus,
		Attributes:           a.DbMan.GetKmsAttributes(cred.TenantId, cred.Id)[cred.Id],
		ConsumerKey:          cred.Id,
		ConsumerSecret:       cred.ConsumerSecret,
		ExpiresAt:            cred.ExpiresAt,
		IssuedAt:             cred.IssuedAt,
		MethodType:           cred.MethodType,
		Scopes:               common.JsonToStringArray(cred.Scopes),
		Status:               cred.Status,
	}
}

func parseIdentifiers(endpoint string, ids map[string]string) (valid bool, keyVals []string) {
	if len(ids) > 2 {
		return false, nil
	}
	if m := IdentifierTree[endpoint]; m != nil {
		for key, val := range ids {
			if m[key] != nil {
				keyVals = append(keyVals, key, val)
				for _, id := range m[key] {
					if ids[id] != "" {
						keyVals = append(keyVals, id, ids[id])
						return true, keyVals
					}
				}
				if len(ids) == 2 {
					return false, nil
				}
				keyVals = append(keyVals, "", "")
				return true, keyVals
			}
		}
	}
	return false, nil
}

func newDbError(err error) *common.ErrorResponse {
	return &common.ErrorResponse{
		ResponseCode:    strconv.Itoa(DB_ERROR),
		ResponseMessage: err.Error(),
		StatusCode:      http.StatusInternalServerError,
	}
}

func newDataError(err error) *common.ErrorResponse {
	return &common.ErrorResponse{
		ResponseCode:    strconv.Itoa(DATA_ERROR),
		ResponseMessage: err.Error(),
		StatusCode:      http.StatusInternalServerError,
	}
}

func writeJson(obj interface{}, w http.ResponseWriter, r *http.Request) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		log.Error("unable to marshal errorResponse: " + err.Error())
		w.Write([]byte("unable to marshal errorResponse: " + err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
	}
}
