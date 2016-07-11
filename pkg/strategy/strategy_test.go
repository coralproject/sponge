/* package source_test is doing unit tests for the source package */
package strategy_test

import (
	"os"

	uuidimported "github.com/pborman/uuid"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"

	. "github.com/coralproject/sponge/pkg/strategy"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Stubing the Strategy Configuration
func fakeStrategy() Strategy {

	var cdatabase CredentialDatabase

	cdatabase = CredentialDatabase{
		Database: "coral",
		Username: "user",
		Password: "password",
		Host:     "host",
		Port:     "5432",
		Adapter:  "mysql",
		Type:     "source",
	}

	cfields := make([]map[string]interface{}, 8)

	cfields[0] = map[string]interface{}{
		"foreign":  "commentid",
		"local":    "CommentID",
		"relation": "Identity",
		"type":     "int",
	}
	cfields[1] = map[string]interface{}{
		"foreign":  "commentbody",
		"local":    "Body",
		"relation": "Identity",
		"type":     "[]byte",
	}
	cfields[2] = map[string]interface{}{
		"foreign":  "parentid",
		"local":    "ParentID",
		"relation": "Identity",
		"type":     "int",
	}
	cfields[3] = map[string]interface{}{
		"foreign":  "assetid",
		"local":    "AssetID",
		"relation": "Identity",
		"type":     "int",
	}
	cfields[4] = map[string]interface{}{
		"foreign":  "statusid",
		"local":    "status",
		"relation": "Status",
		"type":     "string",
	}
	cfields[5] = map[string]interface{}{
		"foreign":        "createdate",
		"local":          "DateCreated",
		"relation":       "Parse",
		"type":           "timedate",
		"datetimeformat": "February 1st, 2006",
	}
	cfields[6] = map[string]interface{}{
		"foreign":  "updatedate",
		"local":    "DateUpdated",
		"relation": "Parse",
		"type":     "timedate",
	}
	cfields[7] = map[string]interface{}{
		"foreign":  "approvedate",
		"local":    "DateApproved",
		"relation": "Parse",
		"type":     "timedate",
	}

	afields := make([]map[string]interface{}, 3)
	afields[0] = map[string]interface{}{
		"foreign":  "assetid",
		"local":    "AssetID",
		"relation": "identity",
		"type":     "int",
	}
	afields[1] = map[string]interface{}{
		"foreign":  "sourceid",
		"local":    "SourceID",
		"relation": "identity",
		"type":     "int",
	}
	afields[2] = map[string]interface{}{
		"foreign":  "asseturl",
		"local":    "URL",
		"relation": "identity",
		"type":     "[]byte",
	}

	ufields := make([]map[string]interface{}, 6)
	ufields[0] = map[string]interface{}{
		"foreign":  "userid",
		"local":    "UserID",
		"relation": "identity",
		"type":     "int",
	}
	ufields[1] = map[string]interface{}{
		"foreign":  "userdisplayname",
		"local":    "UserName",
		"relation": "identity",
		"type":     "[]byte",
	}

	var status = map[string]string{
		"ModeratorApproved": "1",
		"Untouched":         "2",
	}
	var fakeConf = Strategy{
		Name: "New York Times",
		Map: Map{
			Foreign:        "mysql",
			DateTimeFormat: "2006-01-02 15:04:05",
			Entities: map[string]Entity{
				"comments": Entity{
					Foreign:        "crnr_comment",
					Local:          "comments",
					OrderBy:        "commentid",
					ID:             "commentid",
					Fields:         cfields,
					Status:         status,
					PillarEndpoint: "/api/import/comment",
				},
				"assets": Entity{
					Foreign:        "crnr_asset",
					Local:          "assets",
					OrderBy:        "assetid",
					ID:             "assetid",
					Fields:         afields,
					PillarEndpoint: "/api/import/asset",
				},
				"users": Entity{
					Foreign:        "crnr_comment",
					Local:          "users",
					OrderBy:        "userid",
					ID:             "commentid",
					Fields:         ufields,
					PillarEndpoint: "/api/import/user",
				},
			},
		},
		Credentials: Credentials{
			Database: cdatabase,
		},
	}

	return fakeConf
}

var _ = Describe("Getting configuration fields from stubbed strategy", func() {

	var (
		fakeConf Strategy
	)

	BeforeEach(func() {
		fakeConf = fakeStrategy()
	})

	Describe("from a stub strategy object", func() {
		Context("with a valid strategy file", func() {
			It("should get the right id", func() {
				Expect(fakeConf.GetIDField("comments")).To(Equal("commentid"))
			})
			It("should get the right credentials", func() {
				credential, err := fakeConf.GetCredential("mysql", "source")
				Expect(err).Should(BeNil())

				dcredential, ok := credential.(CredentialDatabase)
				Expect(ok).To(Equal(true))
				Expect(dcredential.Database).To(Equal("coral"))
				Expect(dcredential.Username).To(Equal("user"))
				Expect(dcredential.Password).To(Equal("password"))
				Expect(dcredential.Host).To(Equal("host"))
				Expect(dcredential.Port).To(Equal("5432"))
				Expect(dcredential.Adapter).To(Equal("mysql"))
				Expect(dcredential.Type).To(Equal("source"))
			})
			It("should get the foreign field", func() {
				strategy := fakeConf.GetMap()
				Expect(strategy.Foreign).To(Equal("mysql"))
			})
			It("should get the right table", func() {
				tables := fakeConf.GetEntities()
				Expect(tables["comments"].Foreign).To(Equal("crnr_comment"))
			})
			It("should have the format for strategy", func() {
				Expect(fakeConf.GetDateTimeFormat("comments", "DateUpdated")).To(Equal("2006-01-02 15:04:05"))
				Expect(fakeConf.GetDateTimeFormat("comments", "DateCreated")).To(Equal("February 1st, 2006"))
			})
			It("should have the status", func() {
				Expect(fakeConf.GetStatus("comments", "ModeratorApproved")).To(Equal("1"))
			})
			It("should have an entity", func() {
				Expect(fakeConf.GetEntityForeignName("comments")).To(Equal("crnr_comment"))
			})
		})
		// Context("with an invalid strategy file", func() {
		// 	It("should not be valid", func() {
		// 		Expect().To(Equal())
		// 	})
		// })

	})
})

var _ = Describe("Getting configuration from strategy file", func() {
	BeforeEach(func() {
		// Initialize logging
		logLevel := func() int {
			ll, err := cfg.Int("LOGGING_LEVEL")
			if err != nil {
				return log.NONE
			}
			return ll
		}

		log.Init(os.Stderr, logLevel)

		Init(uuidimported.New())
	})

	Describe("", func() {
		Context("with a valid strategy file", func() {
			It("should be valid", func() {
				validstrategyfile := os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/test/strategy_api_test.json"
				Expect(Validate(validstrategyfile)).To(Equal(true))
			})
		})
		It("should not be valid", func() {
			notvalidstrategyfile := os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/test/not_valid_Strategy_test.json"
			Expect(Validate(notvalidstrategyfile)).To(Equal(false))
		})
	})
})
