package banner

import (
	"fmt"
)

// prints the version message
const version = "v0.0.1"

func PrintVersion() {
	fmt.Printf("Current pathcrawler version %s\n", version)
}

// Prints the Colorful banner
func PrintBanner() {
	banner := `
                   __   __                                  __           
    ____   ____ _ / /_ / /_   _____ _____ ____ _ _      __ / /___   _____
   / __ \ / __  // __// __ \ / ___// ___// __  /| | /| / // // _ \ / ___/
  / /_/ // /_/ // /_ / / / // /__ / /   / /_/ / | |/ |/ // //  __// /    
 / .___/ \__,_/ \__//_/ /_/ \___//_/    \__,_/  |__/|__//_/ \___//_/     
/_/`
	fmt.Printf("%s\n%55s\n\n", banner, "Current pathcrawler version "+version)
}
