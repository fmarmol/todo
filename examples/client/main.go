/*basic implementation of client*/
package main

import "flag"

const (
	LIST   string = "list"
	CREATE        = "create"
	UPDATE        = "update"
	DELETE        = "delete"
)

func main() {
	action := flag.String("action", "list", "choose action from list|create|update|delete")
	uid := flag.Int("uid", 1, "uuid task")
	flag.Parse()
}
