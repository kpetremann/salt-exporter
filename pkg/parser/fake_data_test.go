package parser_test

import (
	"log"

	"github.com/vmihailenco/msgpack/v5"
)

type FakeData struct {
	Arg       []any    `msgpack:"arg"`
	Cmd       string   `msgpack:"cmd"`
	Fun       string   `msgpack:"fun"`
	FunArgs   []any    `msgpack:"fun_args"`
	ID        string   `msgpack:"id"`
	Jid       string   `msgpack:"jid"`
	Minions   []string `msgpack:"minions"`
	Missing   []string `msgpack:"missing"`
	Retcode   int      `msgpack:"retcode"`
	Return    any      `msgpack:"return"`
	Schedule  string   `msgpack:"schedule"`
	Success   *bool    `msgpack:"success"`
	Tgt       any      `msgpack:"tgt"`
	TgtType   string   `msgpack:"tgt_type"`
	Timestamp string   `msgpack:"_stamp"`
	User      string   `msgpack:"user"`
	Out       string   `msgpack:"out"`
}

func fakeEventAsMap(event []byte) map[string]any {
	var m any

	if err := msgpack.Unmarshal(event, &m); err != nil {
		log.Fatalln(err)
	}

	return map[string]any{"body": event}
}
