package connection

import (
	"context"

	guard "github.com/go-go-code/goost-guard"
)

type Connection interface {
	Close()
}

var _ctx *context.Context

var cfg map[string]any
var connections []Connection

func init() {

	cfg = map[string]any{}

	guard.Deploy()
}

func add(c Connection) {

	connections = append(connections, c)
}

func SetContext(ctx *context.Context) {

	_ctx = ctx
}

func SetConfig(configs map[string]any) {

	cfg = configs
}

func Close() {

	for _, s := range connections {
		s.Close()
	}
}
