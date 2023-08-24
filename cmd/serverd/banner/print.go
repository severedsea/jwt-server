package banner

import (
	"log"
)

// Print prints the ASCII banner to the console
func Print() {
	// Note: Generated ASCII banner online: http://patorjk.com/software/taag/#p=display&f=ANSI%20Shadow
	log.Println(`---------------------------`)
	log.Println(`     ██╗██╗    ██╗████████╗`)
	log.Println(`     ██║██║    ██║╚══██╔══╝`)
	log.Println(`     ██║██║ █╗ ██║   ██║   `)
	log.Println(`██   ██║██║███╗██║   ██║   `)
	log.Println(`╚█████╔╝╚███╔███╔╝   ██║   `)
	log.Println(` ╚════╝  ╚══╝╚══╝    ╚═╝   `)
	log.Println(`---------------------------`)
}
