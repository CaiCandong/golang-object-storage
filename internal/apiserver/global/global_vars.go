package global

var (
	ListenAddr = ""
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
