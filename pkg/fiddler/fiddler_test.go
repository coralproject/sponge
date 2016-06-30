package fiddler_test

import (
	"fmt"
	"os"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	uuidimported "github.com/pborman/uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	. "github.com/coralproject/sponge/pkg/fiddler"
)

var _ = Describe("Transform row of data", func() {

	var (
		oStrategy  string
		oPillarURL string
	)

	BeforeEach(func() {
		// Save original enviroment variables
		oStrategy = os.Getenv("STRATEGY_CONF")
		oPillarURL = os.Getenv("PILLAR_URL")

		// Initialize log
		logLevel := func() int {
			ll, err := cfg.Int("LOGGING_LEVEL")
			if err != nil {
				return log.NONE
			}
			return ll
		}

		log.Init(os.Stderr, logLevel)

		// Mock strategy configuration
		strategyConf := os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/test/strategy_fiddler_test.json"
		e := os.Setenv("STRATEGY_CONF", strategyConf)
		if e != nil {
			fmt.Println("It could not setup the mock strategy conf variable")
		}

		u := uuidimported.New()

		// Initialize fiddler
		Init(u)
	})

	AfterEach(func() {
		// recover the environment variables
		os.Setenv("STRATEGY_CONF", oStrategy)
		os.Setenv("PILLAR_URL", oPillarURL)
	})

	DescribeTable("no error",
		func(row map[string]interface{}, model string) {
			_, _, err := TransformRow(row, model)
			Expect(err).Should(BeNil())
		},
		Entry("asssets with no document array", map[string]interface{}{"assetid": "3416344", "asseturl": "http://www.nytimes.com/interactive/2014/11/24/us/north-dakota-oil-boom-politics.html", "updatedate": "2014-12-04 00:01:11", "createdate": "2014-12-04 00:01:11"}, "assets"),
		Entry("action likes with document array", map[string]interface{}{
			"provider.icon": "http://cdn.echoenabled.com/images/echo.png", "object.accumulators.likesCount": 1,
			"object.status": "Untouched", "actor.avatar": "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.permalink": "",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.avatar:https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png": "",
			"actor._id": "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
			"actor.id":  "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
			"targets": []map[string]string{
				map[string]string{"conversationID": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f",
					"id": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f"},
			},
			"object.tags.0":        "replyto_Tropicat",
			"object.content":       "Probably nothing since otherwise these folks would just have been sitting around the hotel waiting for meetings to start.",
			"object.context.0.uri": "http://washingtonpost.com/news/to-your-health/wp/2015/05/31/no-stranger-to-brutal-sports-injuries-kerry-faces-a-long-road-to-recovery/",
			"source.name":          "washpost.com", "provider.name": "echo", "object.likes.0.actor.title": "Yersinia_pestis",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.published:2015-06-03T16:50:15Z": "",
			"updated": "2015-06-03 09:50:15.668 -0700 PDT", "ip": "10.128.133.132",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.title:Yersinia_pestis":                                                                        "",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.id:http://washingtonpost.com/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/": "",
			"object.id":           "http://washingtonpost.com/ECHO/item/fc3be552-cb73-45e5-9d50-73b1b754663b",
			"actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"id":           "http://js-kit.com/activities/post/fc3be552-cb73-45e5-9d50-73b1b754663b",
			"verbs":        []string{"http://activitystrea.ms/schema/1.0/post"},
			"provider.uri": "http://aboutecho.com/", "object.content_type": "html",
			"object.objectTypes.0": "http://activitystrea.ms/schema/1.0/comment", "actor.status": "ModeratorApproved",
			"postedTime": "2015-05-31 17:00:12.683 -0700 PDT", "_id": "ObjectIdHex(\"556ba08cd710290035cf6c74\")",
			"object.published": "2015-06-01T00:00:12Z",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.objectTypes.0:http://activitystrea.ms/schema/1.0/person": "",
			"object.context.0.title": "", "actor.title": "Zeus Mom",
			"object.likes.0.actor.id":            "http://washingtonpost.com/user0/",
			"object.likes.0.actor.avatar":        "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.likes.0.actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"object.likes.0.published":           "2015-06-03T16:50:15Z",
		}, "actionslikes"),
		Entry("action likes with array type", map[string]interface{}{
			"provider.icon": "http://cdn.echoenabled.com/images/echo.png", "object.accumulators.likesCount": 1,
			"object.status": "Untouched", "actor.avatar": "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.permalink": "",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.avatar:https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png": "",
			"actor._id": "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
			"actor.id":  "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
			"targets": []map[string]string{
				map[string]string{"conversationID": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f",
					"id": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f"},
			},
			"object.tags.0":        "replyto_Tropicat",
			"object.content":       "Probably nothing since otherwise these folks would just have been sitting around the hotel waiting for meetings to start.",
			"object.context.0.uri": "http://washingtonpost.com/news/to-your-health/wp/2015/05/31/no-stranger-to-brutal-sports-injuries-kerry-faces-a-long-road-to-recovery/",
			"source.name":          "washpost.com", "provider.name": "echo", "object.likes.0.actor.title": "Yersinia_pestis",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.published:2015-06-03T16:50:15Z": "",
			"updated": "2015-06-03 09:50:15.668 -0700 PDT", "ip": "10.128.133.132",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.title:Yersinia_pestis":                                                                        "",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.id:http://washingtonpost.com/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/": "",
			"object.id":           "http://washingtonpost.com/ECHO/item/fc3be552-cb73-45e5-9d50-73b1b754663b",
			"actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"id":           "http://js-kit.com/activities/post/fc3be552-cb73-45e5-9d50-73b1b754663b",
			"verbs":        []string{"http://activitystrea.ms/schema/1.0/post"},
			"provider.uri": "http://aboutecho.com/", "object.content_type": "html",
			"object.objectTypes.0": "http://activitystrea.ms/schema/1.0/comment", "actor.status": "ModeratorApproved",
			"postedTime": "2015-05-31 17:00:12.683 -0700 PDT", "_id": "ObjectIdHex(\"556ba08cd710290035cf6c74\")",
			"object.published": "2015-06-01T00:00:12Z",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.objectTypes.0:http://activitystrea.ms/schema/1.0/person": "",
			"object.context.0.title": "", "actor.title": "Zeus Mom",
			"object.likes.0.actor.id":            "http://washingtonpost.com/user0/",
			"object.likes.0.actor.avatar":        "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.likes.0.actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"object.likes.0.published":           "2015-06-03T16:50:15Z",
			"object.likes.1.actor.id":            "http://washingtonpost.com/user1/",
			"object.likes.1.actor.avatar":        "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.likes.1.actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"object.likes.1.published":           "2015-06-03T16:50:15Z",
		}, "actionslikes"),
		Entry("array type with null values", map[string]interface{}{
			"provider.icon": "http://cdn.echoenabled.com/images/echo.png", "object.accumulators.likesCount": 1,
			"object.status": "Untouched", "actor.avatar": "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.permalink": "",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.avatar:https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png": "",
			"actor._id": "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
			"actor.id":  "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
			"targets": []map[string]string{
				map[string]string{"conversationID": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f",
					"id": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f"},
			},
			"object.tags.0":        "replyto_Tropicat",
			"object.content":       "Probably nothing since otherwise these folks would just have been sitting around the hotel waiting for meetings to start.",
			"object.context.0.uri": "http://washingtonpost.com/news/to-your-health/wp/2015/05/31/no-stranger-to-brutal-sports-injuries-kerry-faces-a-long-road-to-recovery/",
			"source.name":          "washpost.com", "provider.name": "echo", "object.likes.0.actor.title": "Yersinia_pestis",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.published:2015-06-03T16:50:15Z": "",
			"updated": "2015-06-03 09:50:15.668 -0700 PDT", "ip": "10.128.133.132",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.title:Yersinia_pestis":                                                                        "",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.id:http://washingtonpost.com/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/": "",
			"object.id":           "http://washingtonpost.com/ECHO/item/fc3be552-cb73-45e5-9d50-73b1b754663b",
			"actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"id":           "http://js-kit.com/activities/post/fc3be552-cb73-45e5-9d50-73b1b754663b",
			"verbs":        []string{"http://activitystrea.ms/schema/1.0/post"},
			"provider.uri": "http://aboutecho.com/", "object.content_type": "html",
			"object.objectTypes.0": "http://activitystrea.ms/schema/1.0/comment", "actor.status": "ModeratorApproved",
			"postedTime": "2015-05-31 17:00:12.683 -0700 PDT", "_id": "ObjectIdHex(\"556ba08cd710290035cf6c74\")",
			"object.published": "2015-06-01T00:00:12Z",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.objectTypes.0:http://activitystrea.ms/schema/1.0/person": "",
			"object.context.0.title": "", "actor.title": "Zeus Mom",
			"object.likes.0.actor.id":            "http://washingtonpost.com/user0/",
			"object.likes.0.actor.avatar":        "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.likes.0.actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"object.likes.0.published":           "2015-06-03T16:50:15Z",
			"object.likes.1.actor.id":            "http://washingtonpost.com/user1/",
			"object.likes.1.actor.avatar":        "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.likes.1.actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"object.likes.1.published":           "2015-06-03T16:50:15Z",
		}, "actionslikes"),
	)

	DescribeTable("with error",
		func(row map[string]interface{}, model string) {
			_, result, err := TransformRow(row, model)
			Expect(err).ShouldNot(BeNil())
			Expect(result).Should(BeNil())
		},
		Entry("incorrect model name", map[string]interface{}{}, "papafrita"),
	)

	DescribeTable("size",
		func(row map[string]interface{}, model string, size int) {
			_, result, _ := TransformRow(row, model)
			Expect(len(result)).To(Equal(size))
		},
		Entry("asssets with no document array, result size is 1", map[string]interface{}{"assetid": "3416344", "asseturl": "http://www.nytimes.com/interactive/2014/11/24/us/north-dakota-oil-boom-politics.html", "updatedate": "2014-12-04 00:01:11", "createdate": "2014-12-04 00:01:11"}, "assets", 1),
		Entry("action likes with document array, result size is 1", map[string]interface{}{
			"provider.icon": "http://cdn.echoenabled.com/images/echo.png", "object.accumulators.likesCount": 1,
			"object.status": "Untouched", "actor.avatar": "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.permalink": "",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.avatar:https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png": "",
			"actor._id": "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
			"actor.id":  "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
			"targets": []map[string]string{
				map[string]string{"conversationID": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f",
					"id": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f"},
			},
			"object.tags.0":        "replyto_Tropicat",
			"object.content":       "Probably nothing since otherwise these folks would just have been sitting around the hotel waiting for meetings to start.",
			"object.context.0.uri": "http://washingtonpost.com/news/to-your-health/wp/2015/05/31/no-stranger-to-brutal-sports-injuries-kerry-faces-a-long-road-to-recovery/",
			"source.name":          "washpost.com", "provider.name": "echo", "object.likes.0.actor.title": "Yersinia_pestis",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.published:2015-06-03T16:50:15Z": "",
			"updated": "2015-06-03 09:50:15.668 -0700 PDT", "ip": "10.128.133.132",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.title:Yersinia_pestis":                                                                        "",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.id:http://washingtonpost.com/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/": "",
			"object.id":           "http://washingtonpost.com/ECHO/item/fc3be552-cb73-45e5-9d50-73b1b754663b",
			"actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"id":           "http://js-kit.com/activities/post/fc3be552-cb73-45e5-9d50-73b1b754663b",
			"verbs":        []string{"http://activitystrea.ms/schema/1.0/post"},
			"provider.uri": "http://aboutecho.com/", "object.content_type": "html",
			"object.objectTypes.0": "http://activitystrea.ms/schema/1.0/comment", "actor.status": "ModeratorApproved",
			"postedTime": "2015-05-31 17:00:12.683 -0700 PDT", "_id": "ObjectIdHex(\"556ba08cd710290035cf6c74\")",
			"object.published": "2015-06-01T00:00:12Z",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.objectTypes.0:http://activitystrea.ms/schema/1.0/person": "",
			"object.context.0.title": "", "actor.title": "Zeus Mom",
			"object.likes.0.actor.id":            "http://washingtonpost.com/user0/",
			"object.likes.0.actor.avatar":        "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.likes.0.actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"object.likes.0.published":           "2015-06-03T16:50:15Z",
		}, "actionslikes", 1),
		Entry("action likes with array type", map[string]interface{}{
			"provider.icon": "http://cdn.echoenabled.com/images/echo.png", "object.accumulators.likesCount": 1,
			"object.status": "Untouched", "actor.avatar": "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.permalink": "",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.avatar:https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png": "",
			"actor._id": "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
			"actor.id":  "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
			"targets": []map[string]string{
				map[string]string{"conversationID": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f",
					"id": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f"},
			},
			"object.tags.0":        "replyto_Tropicat",
			"object.content":       "Probably nothing since otherwise these folks would just have been sitting around the hotel waiting for meetings to start.",
			"object.context.0.uri": "http://washingtonpost.com/news/to-your-health/wp/2015/05/31/no-stranger-to-brutal-sports-injuries-kerry-faces-a-long-road-to-recovery/",
			"source.name":          "washpost.com", "provider.name": "echo", "object.likes.0.actor.title": "Yersinia_pestis",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.published:2015-06-03T16:50:15Z": "",
			"updated": "2015-06-03 09:50:15.668 -0700 PDT", "ip": "10.128.133.132",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.title:Yersinia_pestis":                                                                        "",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.id:http://washingtonpost.com/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/": "",
			"object.id":           "http://washingtonpost.com/ECHO/item/fc3be552-cb73-45e5-9d50-73b1b754663b",
			"actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"id":           "http://js-kit.com/activities/post/fc3be552-cb73-45e5-9d50-73b1b754663b",
			"verbs":        []string{"http://activitystrea.ms/schema/1.0/post"},
			"provider.uri": "http://aboutecho.com/", "object.content_type": "html",
			"object.objectTypes.0": "http://activitystrea.ms/schema/1.0/comment", "actor.status": "ModeratorApproved",
			"postedTime": "2015-05-31 17:00:12.683 -0700 PDT", "_id": "ObjectIdHex(\"556ba08cd710290035cf6c74\")",
			"object.published": "2015-06-01T00:00:12Z",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.objectTypes.0:http://activitystrea.ms/schema/1.0/person": "",
			"object.context.0.title": "", "actor.title": "Zeus Mom",
			"object.likes.0.actor.id":            "http://washingtonpost.com/user0/",
			"object.likes.0.actor.avatar":        "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.likes.0.actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"object.likes.0.published":           "2015-06-03T16:50:15Z",
			"object.likes.1.actor.id":            "http://washingtonpost.com/user1/",
			"object.likes.1.actor.avatar":        "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.likes.1.actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"object.likes.1.published":           "2015-06-03T16:50:15Z",
		}, "actionslikes", 2),
		Entry("array type with null values", map[string]interface{}{
			"provider.icon": "http://cdn.echoenabled.com/images/echo.png", "object.accumulators.likesCount": 1,
			"object.status": "Untouched", "actor.avatar": "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.permalink": "",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.avatar:https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png": "",
			"actor._id": "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
			"actor.id":  "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
			"targets": []map[string]string{
				map[string]string{"conversationID": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f",
					"id": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f"},
			},
			"object.tags.0":        "replyto_Tropicat",
			"object.content":       "Probably nothing since otherwise these folks would just have been sitting around the hotel waiting for meetings to start.",
			"object.context.0.uri": "http://washingtonpost.com/news/to-your-health/wp/2015/05/31/no-stranger-to-brutal-sports-injuries-kerry-faces-a-long-road-to-recovery/",
			"source.name":          "washpost.com", "provider.name": "echo", "object.likes.0.actor.title": "Yersinia_pestis",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.published:2015-06-03T16:50:15Z": "",
			"updated": "2015-06-03 09:50:15.668 -0700 PDT", "ip": "10.128.133.132",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.title:Yersinia_pestis":                                                                        "",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.id:http://washingtonpost.com/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/": "",
			"object.id":           "http://washingtonpost.com/ECHO/item/fc3be552-cb73-45e5-9d50-73b1b754663b",
			"actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"id":           "http://js-kit.com/activities/post/fc3be552-cb73-45e5-9d50-73b1b754663b",
			"verbs":        []string{"http://activitystrea.ms/schema/1.0/post"},
			"provider.uri": "http://aboutecho.com/", "object.content_type": "html",
			"object.objectTypes.0": "http://activitystrea.ms/schema/1.0/comment", "actor.status": "ModeratorApproved",
			"postedTime": "2015-05-31 17:00:12.683 -0700 PDT", "_id": "ObjectIdHex(\"556ba08cd710290035cf6c74\")",
			"object.published": "2015-06-01T00:00:12Z",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.objectTypes.0:http://activitystrea.ms/schema/1.0/person": "",
			"object.context.0.title": "", "actor.title": "Zeus Mom",
			"object.likes.0.actor.id":            "http://washingtonpost.com/user0/",
			"object.likes.0.actor.avatar":        "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.likes.0.actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"object.likes.0.published":           "2015-06-03T16:50:15Z",
			"object.likes.1.actor.id":            "http://washingtonpost.com/user1/",
			"object.likes.1.actor.avatar":        "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.likes.1.actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"object.likes.1.published":           "2015-06-03T16:50:15Z",
		}, "actionslikes", 2),
	)

	DescribeTable("date",
		func(row map[string]interface{}, model string, date string) {
			_, result, _ := TransformRow(row, model)
			Expect(result[0]["date_updated"]).To(Equal(date))
		},
		Entry("asssets with no document array, expected date", map[string]interface{}{"assetid": "3416344", "asseturl": "http://www.nytimes.com/interactive/2014/11/24/us/north-dakota-oil-boom-politics.html", "updatedate": "2014-12-04 00:01:11", "createdate": "2014-12-04 00:01:11"}, "assets", "2014-12-04T00:01:11Z"),
		// Entry("action likes with document array, result size is 1", map[string]interface{}{
		// 	"provider.icon": "http://cdn.echoenabled.com/images/echo.png", "object.accumulators.likesCount": 1,
		// 	"object.status": "Untouched", "actor.avatar": "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
		// 	"object.permalink": "",
		// 	"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.avatar:https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png": "",
		// 	"actor._id": "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
		// 	"actor.id":  "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
		// 	"targets": []map[string]string{
		// 		map[string]string{"conversationID": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f",
		// 			"id": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f"},
		// 	},
		// 	"object.tags.0":        "replyto_Tropicat",
		// 	"object.content":       "Probably nothing since otherwise these folks would just have been sitting around the hotel waiting for meetings to start.",
		// 	"object.context.0.uri": "http://washingtonpost.com/news/to-your-health/wp/2015/05/31/no-stranger-to-brutal-sports-injuries-kerry-faces-a-long-road-to-recovery/",
		// 	"source.name":          "washpost.com", "provider.name": "echo", "object.likes.0.actor.title": "Yersinia_pestis",
		// 	"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.published:2015-06-03T16:50:15Z": "",
		// 	"updated": "2015-06-03 09:50:15.668 -0700 PDT", "ip": "10.128.133.132",
		// 	"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.title:Yersinia_pestis":                                                                        "",
		// 	"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.id:http://washingtonpost.com/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/": "",
		// 	"object.id":           "http://washingtonpost.com/ECHO/item/fc3be552-cb73-45e5-9d50-73b1b754663b",
		// 	"actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
		// 	"id":           "http://js-kit.com/activities/post/fc3be552-cb73-45e5-9d50-73b1b754663b",
		// 	"verbs":        []string{"http://activitystrea.ms/schema/1.0/post"},
		// 	"provider.uri": "http://aboutecho.com/", "object.content_type": "html",
		// 	"object.objectTypes.0": "http://activitystrea.ms/schema/1.0/comment", "actor.status": "ModeratorApproved",
		// 	"postedTime": "2015-05-31 17:00:12.683 -0700 PDT", "_id": "ObjectIdHex(\"556ba08cd710290035cf6c74\")",
		// 	"object.published": "2015-06-01T00:00:12Z",
		// 	"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.objectTypes.0:http://activitystrea.ms/schema/1.0/person": "",
		// 	"object.context.0.title": "", "actor.title": "Zeus Mom",
		// 	"object.likes.0.actor.id":            "http://washingtonpost.com/user0/",
		// 	"object.likes.0.actor.avatar":        "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
		// 	"object.likes.0.actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
		// 	"object.likes.0.published":           "2015-06-03T16:50:15Z",
		// }, "actionslikes", 1),
	)

	DescribeTable("url",
		func(row map[string]interface{}, model string, url string) {
			_, result, _ := TransformRow(row, model)
			Expect(result[0]["url"]).To(Equal(url))
		},
		Entry("asssets with no document array, expected url", map[string]interface{}{"assetid": "3416344", "asseturl": "http://www.nytimes.com/interactive/2014/11/24/us/north-dakota-oil-boom-politics.html", "updatedate": "2014-12-04 00:01:11", "createdate": "2014-12-04 00:01:11"}, "assets", "http://www.nytimes.com/interactive/2014/11/24/us/north-dakota-oil-boom-politics.html"),
	)

	DescribeTable("ID",
		func(row map[string]interface{}, model string, expectedid string) {
			id, _, _ := TransformRow(row, model)
			Expect(id).To(Equal(expectedid))
		},
		Entry("asssets with no document array, expected id", map[string]interface{}{"assetid": "3416344", "asseturl": "http://www.nytimes.com/interactive/2014/11/24/us/north-dakota-oil-boom-politics.html", "updatedate": "2014-12-04 00:01:11", "createdate": "2014-12-04 00:01:11"}, "assets", "3416344"),
		Entry("action likes with array type", map[string]interface{}{
			"provider.icon": "http://cdn.echoenabled.com/images/echo.png", "object.accumulators.likesCount": 1,
			"object.status": "Untouched", "actor.avatar": "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.permalink": "",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.avatar:https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png": "",
			"actor._id": "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
			"actor.id":  "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
			"targets": []map[string]string{
				map[string]string{"conversationID": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f",
					"id": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f"},
			},
			"object.tags.0":        "replyto_Tropicat",
			"object.content":       "Probably nothing since otherwise these folks would just have been sitting around the hotel waiting for meetings to start.",
			"object.context.0.uri": "http://washingtonpost.com/news/to-your-health/wp/2015/05/31/no-stranger-to-brutal-sports-injuries-kerry-faces-a-long-road-to-recovery/",
			"source.name":          "washpost.com", "provider.name": "echo", "object.likes.0.actor.title": "Yersinia_pestis",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.published:2015-06-03T16:50:15Z": "",
			"updated": "2015-06-03 09:50:15.668 -0700 PDT", "ip": "10.128.133.132",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.title:Yersinia_pestis":                                                                        "",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.id:http://washingtonpost.com/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/": "",
			"object.id":           "http://washingtonpost.com/ECHO/item/fc3be552-cb73-45e5-9d50-73b1b754663b",
			"actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"id":           "http://js-kit.com/activities/post/fc3be552-cb73-45e5-9d50-73b1b754663b",
			"verbs":        []string{"http://activitystrea.ms/schema/1.0/post"},
			"provider.uri": "http://aboutecho.com/", "object.content_type": "html",
			"object.objectTypes.0": "http://activitystrea.ms/schema/1.0/comment", "actor.status": "ModeratorApproved",
			"postedTime": "2015-05-31 17:00:12.683 -0700 PDT", "_id": "ObjectIdHex(\"556ba08cd710290035cf6c74\")",
			"object.published": "2015-06-01T00:00:12Z",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.objectTypes.0:http://activitystrea.ms/schema/1.0/person": "",
			"object.context.0.title": "", "actor.title": "Zeus Mom",
			"object.likes.0.actor.id":            "http://washingtonpost.com/user0/",
			"object.likes.0.actor.avatar":        "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.likes.0.actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"object.likes.0.published":           "2015-06-03T16:50:15Z",
			"object.likes.1.actor.id":            "http://washingtonpost.com/user1/",
			"object.likes.1.actor.avatar":        "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.likes.1.actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"object.likes.1.published":           "2015-06-03T16:50:15Z",
		}, "actionslikes", "ObjectIdHex(\"556ba08cd710290035cf6c74\")"),
		Entry("array type with null values", map[string]interface{}{
			"provider.icon": "http://cdn.echoenabled.com/images/echo.png", "object.accumulators.likesCount": 1,
			"object.status": "Untouched", "actor.avatar": "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.permalink": "",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.avatar:https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png": "",
			"actor._id": "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
			"actor.id":  "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
			"targets": []map[string]string{
				map[string]string{"conversationID": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f",
					"id": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f"},
			},
			"object.tags.0":        "replyto_Tropicat",
			"object.content":       "Probably nothing since otherwise these folks would just have been sitting around the hotel waiting for meetings to start.",
			"object.context.0.uri": "http://washingtonpost.com/news/to-your-health/wp/2015/05/31/no-stranger-to-brutal-sports-injuries-kerry-faces-a-long-road-to-recovery/",
			"source.name":          "washpost.com", "provider.name": "echo", "object.likes.0.actor.title": "Yersinia_pestis",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.published:2015-06-03T16:50:15Z": "",
			"updated": "2015-06-03 09:50:15.668 -0700 PDT", "ip": "10.128.133.132",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.title:Yersinia_pestis":                                                                        "",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.id:http://washingtonpost.com/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/": "",
			"object.id":           "http://washingtonpost.com/ECHO/item/fc3be552-cb73-45e5-9d50-73b1b754663b",
			"actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"id":           "http://js-kit.com/activities/post/fc3be552-cb73-45e5-9d50-73b1b754663b",
			"verbs":        []string{"http://activitystrea.ms/schema/1.0/post"},
			"provider.uri": "http://aboutecho.com/", "object.content_type": "html",
			"object.objectTypes.0": "http://activitystrea.ms/schema/1.0/comment", "actor.status": "ModeratorApproved",
			"postedTime": "2015-05-31 17:00:12.683 -0700 PDT", "_id": "ObjectIdHex(\"556ba08cd710290035cf6c74\")",
			"object.published": "2015-06-01T00:00:12Z",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.objectTypes.0:http://activitystrea.ms/schema/1.0/person": "",
			"object.context.0.title": "", "actor.title": "Zeus Mom",
			"object.likes.0.actor.id":            "http://washingtonpost.com/user0/",
			"object.likes.0.actor.avatar":        "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.likes.0.actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"object.likes.0.published":           "2015-06-03T16:50:15Z",
			"object.likes.1.actor.id":            "http://washingtonpost.com/user1/",
			"object.likes.1.actor.avatar":        "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.likes.1.actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"object.likes.1.published":           "2015-06-03T16:50:15Z",
		}, "actionslikes", "ObjectIdHex(\"556ba08cd710290035cf6c74\")"),
	)

	DescribeTable("Different IDs",
		func(row map[string]interface{}, model string) {
			_, result, _ := TransformRow(row, model)
			result0 := result[0]["source"].(map[string]interface{})["user"].(map[string]interface{})["source"].(map[string]interface{})["id"]
			result1 := result[1]["source"].(map[string]interface{})["user"].(map[string]interface{})["source"].(map[string]interface{})["id"]
			Expect(result0).ToNot(Equal(result1))
		},
		Entry("action likes with array type", map[string]interface{}{
			"provider.icon": "http://cdn.echoenabled.com/images/echo.png", "object.accumulators.likesCount": 1,
			"object.status": "Untouched", "actor.avatar": "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.permalink": "",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.avatar:https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png": "",
			"actor._id": "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
			"actor.id":  "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
			"targets": []map[string]string{
				map[string]string{"conversationID": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f",
					"id": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f"},
			},
			"object.tags.0":        "replyto_Tropicat",
			"object.content":       "Probably nothing since otherwise these folks would just have been sitting around the hotel waiting for meetings to start.",
			"object.context.0.uri": "http://washingtonpost.com/news/to-your-health/wp/2015/05/31/no-stranger-to-brutal-sports-injuries-kerry-faces-a-long-road-to-recovery/",
			"source.name":          "washpost.com", "provider.name": "echo", "object.likes.0.actor.title": "Yersinia_pestis",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.published:2015-06-03T16:50:15Z": "",
			"updated": "2015-06-03 09:50:15.668 -0700 PDT", "ip": "10.128.133.132",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.title:Yersinia_pestis":                                                                        "",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.id:http://washingtonpost.com/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/": "",
			"object.id":           "http://washingtonpost.com/ECHO/item/fc3be552-cb73-45e5-9d50-73b1b754663b",
			"actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"id":           "http://js-kit.com/activities/post/fc3be552-cb73-45e5-9d50-73b1b754663b",
			"verbs":        []string{"http://activitystrea.ms/schema/1.0/post"},
			"provider.uri": "http://aboutecho.com/", "object.content_type": "html",
			"object.objectTypes.0": "http://activitystrea.ms/schema/1.0/comment", "actor.status": "ModeratorApproved",
			"postedTime": "2015-05-31 17:00:12.683 -0700 PDT", "_id": "ObjectIdHex(\"556ba08cd710290035cf6c74\")",
			"object.published": "2015-06-01T00:00:12Z",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.objectTypes.0:http://activitystrea.ms/schema/1.0/person": "",
			"object.context.0.title": "", "actor.title": "Zeus Mom",
			"object.likes.0.actor.id":            "http://washingtonpost.com/user0/",
			"object.likes.0.actor.avatar":        "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.likes.0.actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"object.likes.0.published":           "2015-06-03T16:50:15Z",
			"object.likes.1.actor.id":            "http://washingtonpost.com/user1/",
			"object.likes.1.actor.avatar":        "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.likes.1.actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"object.likes.1.published":           "2015-06-03T16:50:15Z",
		}, "actionslikes"),
		Entry("array type with null values", map[string]interface{}{
			"provider.icon": "http://cdn.echoenabled.com/images/echo.png", "object.accumulators.likesCount": 1,
			"object.status": "Untouched", "actor.avatar": "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.permalink": "",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.avatar:https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png": "",
			"actor._id": "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
			"actor.id":  "http://washingtonpost.com/yH5FvK5Hcr6lmQeD6Xcx8fJkV59ZvvsMzHeNJ1Se1fpoGeVxxhJF5A%3D%3D/",
			"targets": []map[string]string{
				map[string]string{"conversationID": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f",
					"id": "http://washingtonpost.com/ECHO/item/2d1d3956-08a4-4aaa-9ffd-22182fbb5b8f"},
			},
			"object.tags.0":        "replyto_Tropicat",
			"object.content":       "Probably nothing since otherwise these folks would just have been sitting around the hotel waiting for meetings to start.",
			"object.context.0.uri": "http://washingtonpost.com/news/to-your-health/wp/2015/05/31/no-stranger-to-brutal-sports-injuries-kerry-faces-a-long-road-to-recovery/",
			"source.name":          "washpost.com", "provider.name": "echo", "object.likes.0.actor.title": "Yersinia_pestis",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.published:2015-06-03T16:50:15Z": "",
			"updated": "2015-06-03 09:50:15.668 -0700 PDT", "ip": "10.128.133.132",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.title:Yersinia_pestis":                                                                        "",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.id:http://washingtonpost.com/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/": "",
			"object.id":           "http://washingtonpost.com/ECHO/item/fc3be552-cb73-45e5-9d50-73b1b754663b",
			"actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"id":           "http://js-kit.com/activities/post/fc3be552-cb73-45e5-9d50-73b1b754663b",
			"verbs":        []string{"http://activitystrea.ms/schema/1.0/post"},
			"provider.uri": "http://aboutecho.com/", "object.content_type": "html",
			"object.objectTypes.0": "http://activitystrea.ms/schema/1.0/comment", "actor.status": "ModeratorApproved",
			"postedTime": "2015-05-31 17:00:12.683 -0700 PDT", "_id": "ObjectIdHex(\"556ba08cd710290035cf6c74\")",
			"object.published": "2015-06-01T00:00:12Z",
			"object.likes_obj.http://washingtonpostcom/E2vgkjF8Dr3osyRlhbNLC%2BwKrkvbT4tmsKR0XWQVNYpoGeVxxhJF5A%3D%3D/.actor.objectTypes.0:http://activitystrea.ms/schema/1.0/person": "",
			"object.context.0.title": "", "actor.title": "Zeus Mom",
			"object.likes.0.actor.id":            "http://washingtonpost.com/user0/",
			"object.likes.0.actor.avatar":        "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.likes.0.actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"object.likes.0.published":           "2015-06-03T16:50:15Z",
			"object.likes.1.actor.id":            "http://washingtonpost.com/user1/",
			"object.likes.1.actor.avatar":        "https://wpidentity.s3.amazonaws.com/assets/images/avatar-default.png",
			"object.likes.1.actor.objectTypes.0": "http://activitystrea.ms/schema/1.0/person",
			"object.likes.1.published":           "2015-06-03T16:50:15Z",
		}, "actionslikes"),
	)
})

var _ = Describe("Getters", func() {
	var (
		oStrategy  string
		oPillarURL string
	)

	BeforeEach(func() {
		// Save original enviroment variables
		oStrategy = os.Getenv("STRATEGY_CONF")
		oPillarURL = os.Getenv("PILLAR_URL")

		// Initialize log
		logLevel := func() int {
			ll, err := cfg.Int("LOGGING_LEVEL")
			if err != nil {
				return log.NONE
			}
			return ll
		}

		log.Init(os.Stderr, logLevel)

		// Mock strategy configuration
		strategyConf := os.Getenv("GOPATH") + "/src/github.com/coralproject/sponge/test/strategy_fiddler_test.json"
		e := os.Setenv("STRATEGY_CONF", strategyConf)
		if e != nil {
			fmt.Println("It could not setup the mock strategy conf variable")
		}

		u := uuidimported.New()

		// Initialize fiddler
		Init(u)
	})

	AfterEach(func() {
		// recover the environment variables
		os.Setenv("STRATEGY_CONF", oStrategy)
		os.Setenv("PILLAR_URL", oPillarURL)
	})

	Describe("working with the valid strategy file", func() {
		Context("and assets", func() {
			It("should get the right id", func() {
				Expect(GetID("assets")).To(Equal("assetid"))
			})
		})
	})
})

// // Signature:  GetCollections() []string {
// func TestGetCollections(t *testing.T) {
//
// 	collections := GetCollections()
//
// 	if !equal(collections, expectedCollections) {
// 		t.Fatalf("got %s , expected %s", collections, expectedCollections)
// 	}
// }
//
// func equal(a []string, b []string) bool {
// 	if len(a) != len(b) {
// 		return false
// 	}
// 	for _, v := range a {
// 		if vnotinb(v, b) {
// 			return false
// 		}
// 	}
//
// 	return true
// }
//
// func vnotinb(v string, b []string) bool {
//
// 	for _, k := range b {
// 		if v == k {
// 			return false
// 		}
// 	}
// 	return true
// }

// // Signature: appendField(source []map[string]interface{}, item interface{}) []map[string]interface{}
// func TestappendField(t *testing.T) {
//
// 	var source []map[string]interface{}
//
// 	source[0] = make(map[string]interface{})
// 	source[0]["asset_id"] = 1
// 	source[1] = make(map[string]interface{})
// 	source[1]["comment_id"] = 2
//
// 	var item map[string]int
//
// 	item = make(map[string]int)
// 	item["user_id"] = 3
//
// 	result := fiddler.appendField(source, item)
//
// 	if result[3]["user_id"] == 3 {
// 		t.Fatalf("got  %v, expected 3", result[3]["user_id"])
// 	}
// }
