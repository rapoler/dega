package test

import (
	"database/sql/driver"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gavv/httpexpect/v2"
	"github.com/jinzhu/gorm/dialects/postgres"
)

// DATA
var postData = map[string]interface{}{
	"title":    "Post",
	"subtitle": "post subtitle",
	"slug":     "post",
	"status":   "publish",
	"page":     true,
	"excerpt":  "post excerpt",
	"description": postgres.Jsonb{
		RawMessage: []byte(`{"time":1617039625490,"blocks":[{"type":"paragraph","data":{"text":"Test Description"}}],"version":"2.19.0"}`),
	},
	"html_description":   "<p>Test Description</p>",
	"is_featured":        false,
	"is_sticky":          true,
	"is_highlighted":     true,
	"featured_medium_id": uint(1),
	"published_date":     time.Now(),
	"format_id":          uint(1),
	"category_ids":       []uint{1},
	"tag_ids":            []uint{1},
	"claim_ids":          []uint{1},
	"author_ids":         []uint{1},
}

var postColumns = []string{"id", "created_at", "updated_at", "deleted_at", "created_by_id", "updated_by_id", "title", "subtitle", "slug", "status", "page", "excerpt", "description", "html_description", "is_featured", "is_sticky", "is_highlighted", "featured_medium_id", "format_id", "published_date", "space_id"}

func TestPosts(t *testing.T) {
	// Setup Mock DB
	mock := SetupMockDB()

	// Start test server
	testServer := httptest.NewServer(TestRouter())
	defer testServer.Close()

	// Setup httpexpect
	e := httpexpect.New(t, testServer.URL)

	// posts testcases
	t.Run("get list of post ids", func(t *testing.T) {
		PostCountMock(mock, 1)
		PostSelectMock(mock)

		resp := e.POST(path).
			WithJSON(Query{
				Query: `{
					posts {
						nodes {
							id
						}
					}
				}`,
			}).Expect().
			JSON().
			Object()

		CheckJSON(resp, map[string]interface{}{
			"nodes": []map[string]interface{}{
				{"id": "1"},
			},
		}, "posts")
		ExpectationsMet(t, mock)
	})

	t.Run("get list of posts in ascending order of slug", func(t *testing.T) {
		PostCountMock(mock, 1)

		mock.ExpectQuery(`SELECT \* FROM "posts" (.+) ORDER BY slug asc`).
			WillReturnRows(sqlmock.NewRows(postColumns).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, postData["title"], postData["subtitle"], postData["slug"], postData["status"], postData["page"], postData["excerpt"], postData["description"], postData["html_description"], postData["is_featured"], postData["is_sticky"], postData["is_highlighted"], postData["featured_medium_id"], postData["format_id"], postData["published_date"], 1))

		resp := e.POST(path).
			WithJSON(Query{
				Query: `{
				posts(sortBy: "slug", sortOrder: "asc") {
						nodes {
							id
							title
							html_description
						}
					}
				}`,
			}).Expect().
			JSON().
			Object()

		CheckJSON(resp, map[string]interface{}{
			"nodes": []map[string]interface{}{
				{"id": "1", "title": postData["title"], "html_description": postData["html_description"]},
			},
		}, "posts")
		ExpectationsMet(t, mock)
	})

	t.Run("get list of draft posts from some spaces", func(t *testing.T) {
		PostCountMock(mock, 1)
		PostSelectMock(mock, 1, 2, "draft")

		resp := e.POST(path).
			WithJSON(Query{
				Query: `{
				posts(status:"draft", spaces:[1,2]) {
						nodes {
							id
							title
							html_description
						}
					}
				}`,
			}).Expect().
			JSON().
			Object()

		CheckJSON(resp, map[string]interface{}{
			"nodes": []map[string]interface{}{
				{"id": "1", "title": postData["title"], "html_description": postData["html_description"]},
			},
		}, "posts")
		ExpectationsMet(t, mock)
	})

	t.Run("filter posts based on categories and tags and formats", func(t *testing.T) {
		PostCountMock(mock, 1)

		mock.ExpectQuery(`SELECT (.+) FROM "posts" INNER JOIN post_categories (.+) INNER JOIN post_tags (.+)category_id IN \(2,3\) (.+)tag_id IN \(1,2\) (.+)format_id IN \(1\)\)`).
			WithArgs("publish").
			WillReturnRows(sqlmock.NewRows(postColumns).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, postData["title"], postData["subtitle"], postData["slug"], postData["status"], postData["page"], postData["excerpt"], postData["description"], postData["html_description"], postData["is_featured"], postData["is_sticky"], postData["is_highlighted"], postData["featured_medium_id"], postData["format_id"], postData["published_date"], 1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "post_categories"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"post_id", "category_id"}).AddRow(1, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "post_tags"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"post_id", "tag_id"}).AddRow(1, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		resp := e.POST(path).
			WithJSON(Query{
				Query: `{
				posts(categories:[2,3], tags:[1,2], formats:[1]) {
						nodes {
							id
							title
							html_description
							categories {
								id
							}
							tags {
								id
							}
							format {
								id
							}
						}
					}
				}`,
			}).Expect().
			JSON().
			Object()

		CheckJSON(resp, map[string]interface{}{
			"nodes": []map[string]interface{}{
				{"id": "1", "title": postData["title"], "html_description": postData["html_description"], "tags": []map[string]interface{}{{"id": "1"}}, "categories": []map[string]interface{}{{"id": "1"}}, "format": map[string]interface{}{"id": "1"}},
			},
		}, "posts")
		ExpectationsMet(t, mock)
	})

	t.Run("filter posts based on authors", func(t *testing.T) {
		PostCountMock(mock, 1)

		mock.ExpectQuery(`SELECT (.+) FROM "posts" INNER JOIN post_authors (.+)author_id IN \(5,6\)`).
			WithArgs("publish").
			WillReturnRows(sqlmock.NewRows(postColumns).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, postData["title"], postData["subtitle"], postData["slug"], postData["status"], postData["page"], postData["excerpt"], postData["description"], postData["html_description"], postData["is_featured"], postData["is_sticky"], postData["is_highlighted"], postData["featured_medium_id"], postData["format_id"], postData["published_date"], 1))

		resp := e.POST(path).
			WithJSON(Query{
				Query: `{
				posts(users:[5,6]) {
						nodes {
							id
							title
							html_description
						}
					}
				}`,
			}).Expect().
			JSON().
			Object()

		CheckJSON(resp, map[string]interface{}{
			"nodes": []map[string]interface{}{
				{"id": "1", "title": postData["title"], "html_description": postData["html_description"]},
			},
		}, "posts")
		ExpectationsMet(t, mock)
	})

}

func PostSelectMock(mock sqlmock.Sqlmock, args ...driver.Value) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "posts"`)).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows(postColumns).
			AddRow(1, time.Now(), time.Now(), nil, 1, 1, postData["title"], postData["subtitle"], postData["slug"], postData["status"], postData["page"], postData["excerpt"], postData["description"], postData["html_description"], postData["is_featured"], postData["is_sticky"], postData["is_highlighted"], postData["featured_medium_id"], postData["format_id"], postData["published_date"], 1))
}

func PostCountMock(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "posts"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}
