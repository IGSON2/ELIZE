package elizeutils

import "log"

func Errchk(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
