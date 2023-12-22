package parser_test

import (
	"log"

	"github.com/kpetremann/salt-exporter/pkg/event"
	"github.com/vmihailenco/msgpack/v5"
)

/*
Fake new beacon message of type /status

	salt/beacon/host1.example.com/status/2023-10-09T11:36:02.182345	{
		{
			"id": "host1.example.com",
			"data": {
				"loadavg": {
					"1-min": 0.35,
					"5-min": 0.48,
					"15-min": 0.26
				}
			},
			"_stamp": "2023-10-09T11:36:02.205686"
		}
	}
*/
var expectedBeacon = event.SaltEvent{
	Tag:          "salt/beacon/host1.example.com/status/2023-10-09T11:36:02.182345",
	Type:         "status",
	Module:       event.BeaconModule,
	TargetNumber: 0,
	Data: event.EventData{
		Timestamp: "2023-10-09T11:36:02.205686",
		ID:        "host1.example.com",
		Minions:   []string{},
	},
	IsScheduleJob: false,
}

func fakeBeaconEvent() []byte {
	// Marshal the data using MsgPack
	fake := FakeData{
		Timestamp: "2023-10-09T11:36:02.205686",
		Minions:   []string{},
		ID:        "host1.example.com",
	}

	fakeBody, err := msgpack.Marshal(fake)
	if err != nil {
		log.Fatalln(err)
	}

	fakeMessage := []byte("salt/beacon/host1.example.com/status/2023-10-09T11:36:02.182345\n\n")
	fakeMessage = append(fakeMessage, fakeBody...)

	return fakeMessage
}
