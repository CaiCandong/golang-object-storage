package global

import (
	"go.uber.org/zap"
	"golang-object-storage/internal/pkg/rs"
)

var (
	ListenAddr = ""
	RsConfig   = rs.DefaultConfig()
	Logger     *zap.SugaredLogger
)

// func CheckSharedVars()  {
//     existError := false
//     if ListenAddr == "" {
//         log.Printf("dataServer address '%s' is invalid\n", ListenAddr)
//         existError = true
//     }

//     if existError {
//         err_utils.PanicNonNilError(fmt.Errorf("Error: please checkout [listenAddr='%s']\n", ListenAddr))
//     }
// }
