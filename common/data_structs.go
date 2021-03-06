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
package common

type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ErrorResponse struct {
	ResponseCode    string `json:"response_code,omitempty"`
	ResponseMessage string `json:"response_message,omitempty"`
	StatusCode      int    `json:"-"`
	Kind            string `json:"kind,omitempty"`
}

func (e *ErrorResponse) Error() string {
	return e.ResponseMessage
}

type ApiProduct struct {
	Id            string `db:"id"`
	TenantId      string `db:"tenant_id"`
	Name          string `db:"name"`
	DisplayName   string `db:"display_name"`
	Description   string `db:"description"`
	ApiResources  string `db:"api_resources"`
	ApprovalType  string `db:"approval_type"`
	Scopes        string `db:"scopes"`
	Proxies       string `db:"proxies"`
	Environments  string `db:"environments"`
	Quota         string `db:"quota"`
	QuotaTimeUnit string `db:"quota_time_unit"`
	QuotaInterval int64  `db:"quota_interval"`
	CreatedAt     string `db:"created_at"`
	CreatedBy     string `db:"created_by"`
	UpdatedAt     string `db:"updated_at"`
	UpdatedBy     string `db:"updated_by"`
}

type App struct {
	Id          string `db:"id"`
	TenantId    string `db:"tenant_id"`
	Name        string `db:"name"`
	DisplayName string `db:"display_name"`
	AccessType  string `db:"access_type"`
	CallbackUrl string `db:"callback_url"`
	Status      string `db:"status"`
	AppFamily   string `db:"app_family"`
	CompanyId   string `db:"company_id"`
	DeveloperId string `db:"developer_id"`
	ParentId    string `db:"parent_id"`
	Type        string `db:"type"`
	CreatedAt   string `db:"created_at"`
	CreatedBy   string `db:"created_by"`
	UpdatedAt   string `db:"updated_at"`
	UpdatedBy   string `db:"updated_by"`
}

type AppCredential struct {
	Id             string `db:"id"`
	TenantId       string `db:"tenant_id"`
	ConsumerSecret string `db:"consumer_secret"`
	AppId          string `db:"app_id"`
	MethodType     string `db:"method_type"`
	Status         string `db:"status"`
	IssuedAt       string `db:"issued_at"`
	ExpiresAt      string `db:"expires_at"`
	AppStatus      string `db:"app_status"`
	Scopes         string `db:"scopes"`
	CreatedAt      string `db:"created_at"`
	CreatedBy      string `db:"created_by"`
	UpdatedAt      string `db:"updated_at"`
	UpdatedBy      string `db:"updated_by"`
}

type Company struct {
	Id          string `db:"id"`
	TenantId    string `db:"tenant_id"`
	Name        string `db:"name"`
	DisplayName string `db:"display_name"`
	Status      string `db:"status"`
	CreatedAt   string `db:"created_at"`
	CreatedBy   string `db:"created_by"`
	UpdatedAt   string `db:"updated_at"`
	UpdatedBy   string `db:"updated_by"`
}

type Developer struct {
	Id                string `db:"id"`
	TenantId          string `db:"tenant_id"`
	UserName          string `db:"username"`
	FirstName         string `db:"first_name"`
	LastName          string `db:"last_name"`
	Password          string `db:"password"`
	Email             string `db:"email"`
	Status            string `db:"status"`
	EncryptedPassword string `db:"encrypted_password"`
	Salt              string `db:"salt"`
	CreatedAt         string `db:"created_at"`
	CreatedBy         string `db:"created_by"`
	UpdatedAt         string `db:"updated_at"`
	UpdatedBy         string `db:"updated_by"`
}

type CompanyDeveloper struct {
	TenantId    string `db:"tenant_id"`
	CompanyId   string `db:"company_id"`
	DeveloperId string `db:"developer_id"`
	Roles       string `db:"roles"`
	CreatedAt   string `db:"created_at"`
	CreatedBy   string `db:"created_by"`
	UpdatedAt   string `db:"updated_at"`
	UpdatedBy   string `db:"updated_by"`
}
