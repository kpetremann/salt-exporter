package parser_test

import (
	"log"

	"github.com/vmihailenco/msgpack/v5"
)

type FakeData struct {
	Arg       []interface{} `msgpack:"arg"`
	Cmd       string        `msgpack:"cmd"`
	Fun       string        `msgpack:"fun"`
	FunArgs   []interface{} `msgpack:"fun_args"`
	ID        string        `msgpack:"id"`
	Jid       string        `msgpack:"jid"`
	Minions   []string      `msgpack:"minions"`
	Missing   []string      `msgpack:"missing"`
	Retcode   int           `msgpack:"retcode"`
	Return    interface{}   `msgpack:"return"`
	Schedule  string        `msgpack:"schedule"`
	Success   *bool         `msgpack:"success"`
	Tgt       interface{}   `msgpack:"tgt"`
	TgtType   string        `msgpack:"tgt_type"`
	Timestamp string        `msgpack:"_stamp"`
	User      string        `msgpack:"user"`
	Out       string        `msgpack:"out"`
}

func fakeEventAsMap(event []byte) map[string]interface{} {
	var m interface{}

	if err := msgpack.Unmarshal(event, &m); err != nil {
		log.Fatalln(err)
	}

	return map[string]interface{}{"body": event}
}
