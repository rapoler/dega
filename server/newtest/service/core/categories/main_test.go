package categories

import (
	"log"
	"os"
	"testing"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/newtest"
	"github.com/factly/dega-server/service/core/model"
	"gopkg.in/h2non/gock.v1"
)

func TestMain(m *testing.M) {
	defer gock.Disable()
	newtest.MockServer()

	newtest.SetupSqlite("./categories.db")
	config.DB.AutoMigrate(&model.Category{}, &model.Medium{})
	defer gock.DisableNetworking()

	exitValue := m.Run()
	if err := os.Remove("./categories.db"); err != nil {
		log.Fatal(err)
	}

	os.Exit(exitValue)
}
