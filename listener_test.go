package apidVerifyApiKey

import (
	"encoding/json"
	"github.com/30x/apid"
	"github.com/apigee-labs/transicator/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
)

var _ = Describe("listener", func() {

	Context("KMS create/updates verification via changes", func() {
		It("Create KMS tables via changes, and Verify via verifyApiKey", func(done Done) {
			var event = common.ChangeList{}
			closed := 0
			/* API Product */
			srvItems := common.Row{
				"id": {
					Value: "ch_api_product_2",
				},
				"apid_resources": {
					Value: "{}",
				},
				"environments": {
					Value: "{Env_0, Env_1}",
				},
				"tenant_id": {
					Value: "tenant_id_0",
				},
				"_change_selector": {
					Value: "test_org0",
				},
			}

			/* DEVELOPER */
			devItems := common.Row{
				"id": {
					Value: "ch_developer_id_2",
				},
				"status": {
					Value: "Active",
				},
				"tenant_id": {
					Value: "tenant_id_0",
				},
				"_change_selector": {
					Value: "test_org0",
				},
			}

			/* APP */
			appItems := common.Row{
				"id": {
					Value: "ch_application_id_2",
				},
				"developer_id": {
					Value: "ch_developer_id_2",
				},
				"status": {
					Value: "Approved",
				},
				"tenant_id": {
					Value: "tenant_id_0",
				},
				"_change_selector": {
					Value: "test_org0",
				},
			}

			/* CRED */
			credItems := common.Row{
				"id": {
					Value: "ch_app_credential_2",
				},
				"app_id": {
					Value: "ch_application_id_2",
				},
				"tenant_id": {
					Value: "tenant_id_0",
				},
				"status": {
					Value: "Approved",
				},
				"_change_selector": {
					Value: "test_org0",
				},
			}

			/* APP_CRED_APIPRD_MAPPER */
			mpItems := common.Row{
				"apiprdt_id": {
					Value: "ch_api_product_2",
				},
				"app_id": {
					Value: "ch_application_id_2",
				},
				"appcred_id": {
					Value: "ch_app_credential_2",
				},
				"status": {
					Value: "Approved",
				},
				"_change_selector": {
					Value: "test_org0",
				},
				"tenant_id": {
					Value: "tenant_id_0",
				},
			}

			event.Changes = []common.Change{
				{
					Table:     "kms.api_product",
					NewRow:    srvItems,
					Operation: 1,
				},
				{
					Table:     "kms.developer",
					NewRow:    devItems,
					Operation: 1,
				},
				{
					Table:     "kms.app",
					NewRow:    appItems,
					Operation: 1,
				},
				{
					Table:     "kms.app_credential",
					NewRow:    credItems,
					Operation: 1,
				},
				{
					Table:     "kms.app_credential_apiproduct_mapper",
					NewRow:    mpItems,
					Operation: 1,
				},
			}

			h := &test_handler{
				"checkDatabase post Insertion",
				func(e apid.Event) {
					defer GinkgoRecover()

					// ignore the first event, let standard listener process it
					changeSet := e.(*common.ChangeList)
					if len(changeSet.Changes) > 0 || closed == 1 {
						return
					}
					v := url.Values{
						"key":       []string{"ch_app_credential_2"},
						"uriPath":   []string{"/test"},
						"scopeuuid": []string{"XYZ"},
						"action":    []string{"verify"},
					}
					rsp, err := verifyAPIKey(v)
					Expect(err).ShouldNot(HaveOccurred())
					var respj kmsResponseSuccess
					json.Unmarshal(rsp, &respj)
					Expect(respj.Type).Should(Equal("APIKeyContext"))
					Expect(respj.RspInfo.Key).Should(Equal("ch_app_credential_2"))
					closed = 1
					close(done)
				},
			}

			apid.Events().Listen("ApigeeSync", h)
			apid.Events().Emit("ApigeeSync", &event)
			apid.Events().Emit("ApigeeSync", &common.ChangeList{})
		})
	})

	It("Modify tables in KMS tables, and verify via verifyApiKey", func(done Done) {
		closed := 0
		var event = common.ChangeList{}
		var event2 = common.ChangeList{}

		/* Orig data */
		/* API Product */
		srvItemsOld := common.Row{
			"id": {
				Value: "ch_api_product_0",
			},
			"apid_resources": {
				Value: "{}",
			},
			"environments": {
				Value: "{Env_0, Env_1}",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"_change_selector": {
				Value: "test_org0",
			},
		}

		/* DEVELOPER */
		devItemsOld := common.Row{
			"id": {
				Value: "ch_developer_id_0",
			},
			"status": {
				Value: "Active",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"_change_selector": {
				Value: "test_org0",
			},
		}

		/* APP */
		appItemsOld := common.Row{
			"id": {
				Value: "ch_application_id_0",
			},
			"developer_id": {
				Value: "ch_developer_id_0",
			},
			"status": {
				Value: "Approved",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"_change_selector": {
				Value: "test_org0",
			},
		}

		/* CRED */
		credItemsOld := common.Row{
			"id": {
				Value: "ch_app_credential_0",
			},
			"app_id": {
				Value: "ch_application_id_0",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"status": {
				Value: "Approved",
			},
			"_change_selector": {
				Value: "test_org0",
			},
		}

		/* APP_CRED_APIPRD_MAPPER */
		mpItemsOld := common.Row{
			"apiprdt_id": {
				Value: "ch_api_product_0",
			},
			"app_id": {
				Value: "ch_application_id_0",
			},
			"appcred_id": {
				Value: "ch_app_credential_0",
			},
			"status": {
				Value: "Approved",
			},
			"_change_selector": {
				Value: "test_org0",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
		}

		/* New to be replaced data */
		/* API PRODUCT */
		srvItemsNew := common.Row{
			"id": {
				Value: "ch_api_product_1",
			},
			"apid_resources": {
				Value: "{}",
			},
			"environments": {
				Value: "{Env_0, Env_1}",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"_change_selector": {
				Value: "test_org0",
			},
		}

		/* DEVELOPER */
		devItemsNew := common.Row{
			"id": {
				Value: "ch_developer_id_1",
			},
			"status": {
				Value: "Active",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"_change_selector": {
				Value: "test_org0",
			},
		}

		/* APP */
		appItemsNew := common.Row{
			"id": {
				Value: "ch_application_id_1",
			},
			"developer_id": {
				Value: "ch_developer_id_1",
			},
			"status": {
				Value: "Approved",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"_change_selector": {
				Value: "test_org0",
			},
		}

		/* CRED */
		credItemsNew := common.Row{
			"id": {
				Value: "ch_app_credential_1",
			},
			"app_id": {
				Value: "ch_application_id_1",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
			"status": {
				Value: "Approved",
			},
			"_change_selector": {
				Value: "test_org0",
			},
		}

		/* APP_CRED_APIPRD_MAPPER */
		mpItemsNew := common.Row{
			"apiprdt_id": {
				Value: "ch_api_product_1",
			},
			"app_id": {
				Value: "ch_application_id_1",
			},
			"appcred_id": {
				Value: "ch_app_credential_1",
			},
			"status": {
				Value: "Approved",
			},
			"_change_selector": {
				Value: "test_org0",
			},
			"tenant_id": {
				Value: "tenant_id_0",
			},
		}

		event.Changes = []common.Change{
			{
				Table:     "kms.api_product",
				NewRow:    srvItemsOld,
				Operation: 1,
			},
			{
				Table:     "kms.developer",
				NewRow:    devItemsOld,
				Operation: 1,
			},
			{
				Table:     "kms.app",
				NewRow:    appItemsOld,
				Operation: 1,
			},
			{
				Table:     "kms.app_credential",
				NewRow:    credItemsOld,
				Operation: 1,
			},
			{
				Table:     "kms.app_credential_apiproduct_mapper",
				NewRow:    mpItemsOld,
				Operation: 1,
			},
		}

		event2.Changes = []common.Change{
			{
				Table:     "kms.api_product",
				OldRow:    srvItemsOld,
				NewRow:    srvItemsNew,
				Operation: 2,
			},
			{
				Table:     "kms.developer",
				OldRow:    devItemsOld,
				NewRow:    devItemsNew,
				Operation: 2,
			},
			{
				Table:     "kms.app",
				OldRow:    appItemsOld,
				NewRow:    appItemsNew,
				Operation: 2,
			},
			{
				Table:     "kms.app_credential",
				OldRow:    credItemsOld,
				NewRow:    credItemsNew,
				Operation: 2,
			},
			{
				Table:     "kms.app_credential_apiproduct_mapper",
				OldRow:    mpItemsOld,
				NewRow:    mpItemsNew,
				Operation: 2,
			},
		}

		h := &test_handler{
			"checkDatabase post Insertion",
			func(e apid.Event) {
				defer GinkgoRecover()

				// ignore the first event, let standard listener process it
				changeSet := e.(*common.ChangeList)
				if len(changeSet.Changes) > 0 || closed == 1 {
					return
				}
				v := url.Values{
					"key":       []string{"ch_app_credential_1"},
					"uriPath":   []string{"/test"},
					"scopeuuid": []string{"XYZ"},
					"action":    []string{"verify"},
				}
				rsp, err := verifyAPIKey(v)
				Expect(err).ShouldNot(HaveOccurred())
				var respj kmsResponseSuccess
				json.Unmarshal(rsp, &respj)
				Expect(respj.Type).Should(Equal("APIKeyContext"))
				Expect(respj.RspInfo.Key).Should(Equal("ch_app_credential_1"))
				closed = 1
				close(done)
			},
		}

		apid.Events().Listen("ApigeeSync", h)
		apid.Events().Emit("ApigeeSync", &event)
		apid.Events().Emit("ApigeeSync", &event2)
		apid.Events().Emit("ApigeeSync", &common.ChangeList{})
	})

})

type test_handler struct {
	description string
	f           func(event apid.Event)
}

func (t *test_handler) String() string {
	return t.description
}

func (t *test_handler) Handle(event apid.Event) {
	t.f(event)
}

func addScopes(db apid.DB) {
	txn, _ := db.Begin()
	txn.Exec("INSERT INTO DATA_SCOPE (id, _change_selector, apid_cluster_id, scope, org, env) "+
		"VALUES"+
		"($1,$2,$3,$4,$5,$6)",
		"ABCDE",
		"some_cluster_id",
		"some_cluster_id",
		"tenant_id_xxxx",
		"test_org0",
		"Env_0",
	)
	txn.Exec("INSERT INTO DATA_SCOPE (id, _change_selector, apid_cluster_id, scope, org, env) "+
		"VALUES"+
		"($1,$2,$3,$4,$5,$6)",
		"XYZ",
		"test_org0",
		"somecluster_id",
		"tenant_id_0",
		"test_org0",
		"Env_0",
	)
	log.Info("Inserted DATA_SCOPE for test")
	txn.Commit()
}
