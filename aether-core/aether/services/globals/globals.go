// Services > Globals
// This file collects all globally accessible entities.

package globals

import (
	"aether-core/aether/services/configstore"
	// "fmt"
	"github.com/asdine/storm"
	"github.com/jmoiron/sqlx"
	cdir "github.com/shibukawa/configdir"
	"os"
	"path/filepath"
)

var FrontendConfig *configstore.FrontendConfig
var FrontendTransientConfig *configstore.FrontendTransientConfig

var BackendConfig *configstore.BackendConfig
var BackendTransientConfig *configstore.BackendTransientConfig

var DbInstance *sqlx.DB
var KvInstance *storm.DB

// GetDbSize gets the size of the database. This is here and not in toolbox because we need to access GetSQLiteDBLocation().
func GetDbSize() int {
	dbLoc := filepath.Join(BackendConfig.GetSQLiteDBLocation(), "AetherDB.db")
	fi, _ := os.Stat(dbLoc)
	// get the size
	size := fi.Size() / 1000000
	return int(size)
}

func GetBackendConfigLocation() string {
	configDirs := cdir.New(BackendTransientConfig.OrgIdentifier, BackendTransientConfig.AppIdentifier)
	folders := configDirs.QueryFolders(cdir.Global)
	return filepath.Join(folders[0].Path, "backend", "backend_config.json")
}

func GetFrontendConfigLocation() string {
	configDirs := cdir.New(FrontendTransientConfig.OrgIdentifier, FrontendTransientConfig.AppIdentifier)
	folders := configDirs.QueryFolders(cdir.Global)
	return filepath.Join(folders[0].Path, "frontend", "frontend_config.json")
}
