package common

import (
	"fmt"

	"github.com/hibare/GoGeoIP/internal/constants"
)

func Banner() {
	fmt.Println( //nolint:forbidigo // printing banner
		`
 _       __                        _       __
| |     / /___ ___  ______  ____  (_)___  / /_
| | /| / / __ ` + ` / / / / __ \/ __ \/ / __ \/ __/
| |/ |/ / /_/ / /_/ / /_/ / /_/ / / / / / /_
|__/|__/\__,_/\__, / .___/\____/_/_/ /_/\__/
             /____/_/
			 `)
	fmt.Printf("\nVersion: %s\n", constants.Version)    //nolint:forbidigo // printing banner
	fmt.Printf("Build: %s\n", constants.BuildTimestamp) //nolint:forbidigo // printing banner
	fmt.Printf("Commit: %s\n\n", constants.CommitHash)  //nolint:forbidigo // printing banner
}
