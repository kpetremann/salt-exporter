package events

import (
	"log"

	"github.com/vmihailenco/msgpack/v5"
)

/*
	Fake new job message of type /new

	salt/job/20220630000000000000/new/localhost	{
		"_stamp": "2022-06-30T00:00:00.000000",
		"arg": [],
		"fun": "test.ping",
		"jid": "20220630000000000000",
		"minions": [
			"localhost"
		],
		"missing": [],
		"tgt": "localhost",
		"tgt_type": "glob",
		"user": "sudo_user"
	}
*/

var expectedNewJob = SaltEvent{
	Tag:          "salt/job/20220630000000000000/new",
	Type:         "new",
	TargetNumber: 1,
	Data: EventData{
		Timestamp: "2022-06-30T00:00:00.000000",
		Fun:       "test.ping",
		Jid:       "20220630000000000000",
		Minions:   []string{"localhost"},
		Tgt:       "localhost",
		TgtType:   "glob",
		User:      "salt_user",
	},
	IsScheduleJob: false,
}

func fakeNewJobEvent() []byte {
	// Marshal the data using MsgPack
	fake := FakeData{
		Timestamp: "2022-06-30T00:00:00.000000",
		Fun:       "test.ping",
		Jid:       "20220630000000000000",
		Minions:   []string{"localhost"},
		Tgt:       "localhost",
		TgtType:   "glob",
		User:      "salt_user",
	}

	fakeBody, err := msgpack.Marshal(fake)
	if err != nil {
		log.Fatalln(err)
	}

	fakeMessage := []byte("salt/job/20220630000000000000/new\n\n")
	fakeMessage = append(fakeMessage, fakeBody...)

	return fakeMessage
}

/*
	Fake new job message of type /ret

	salt/job/20220630000000000000/ret/localhost	{
		"_stamp": "2022-06-30T00:00:00.000000",
		"cmd": "_return",
		"fun": "test.ping",
		"fun_args": [],
		"id": "localhost",
		"jid": "20220630000000000000",
		"retcode": 0,
		"return": true,
		"success": true
	}

*/

var expectedReturnJob = SaltEvent{
	Tag:          "salt/job/20220630000000000000/ret/localhost",
	Type:         "ret",
	TargetNumber: 0,
	Data: EventData{
		Timestamp: "2022-06-30T00:00:00.000000",
		Cmd:       "_return",
		Fun:       "test.ping",
		Id:        "localhost",
		Jid:       "20220630000000000000",
		Retcode:   0,
		Return:    true,
		Success:   true,
	},
	IsScheduleJob: false,
}

func fakeRetJobEvent() []byte {
	// Marshal the data using MsgPack
	fake := FakeData{
		Timestamp: "2022-06-30T00:00:00.000000",
		Cmd:       "_return",
		Fun:       "test.ping",
		Id:        "localhost",
		Jid:       "20220630000000000000",
		Retcode:   0,
		Return:    true,
		Success:   true,
	}

	fakeBody, err := msgpack.Marshal(fake)
	if err != nil {
		log.Fatalln(err)
	}

	fakeMessage := []byte("salt/job/20220630000000000000/ret/localhost\n\n")
	fakeMessage = append(fakeMessage, fakeBody...)

	return fakeMessage
}

/*
	Fake manual scheduled job trigger
	salt/job/20220630000000000000/new	{
		"_stamp": "2022-06-30T00:00:00.000000",
		"arg": [
			"sync_all"
		],
		"fun": "schedule.run_job",
		"jid": "20220630000000000000",
		"minions": [
			"localhost"
		],
		"missing": [],
		"tgt": "localhost",
		"tgt_type": "glob",
		"user": "salt_user"
	}
*/

var expectedNewScheduleJob = SaltEvent{
	Tag:          "salt/job/20220630000000000000/new",
	Type:         "new",
	TargetNumber: 1,
	Data: EventData{
		Timestamp: "2022-06-30T00:00:00.000000",
		FunArgs:   []interface{}{"sync_all"},
		Fun:       "schedule.run_job",
		Jid:       "20220630000000000000",
		Minions:   []string{"localhost"},
		Tgt:       "localhost",
		TgtType:   "glob",
		User:      "salt_user",
	},
	IsScheduleJob: false,
}

func fakeNewScheduleJobEvent() []byte {
	// Marshal the data using MsgPack
	fake := FakeData{
		Timestamp: "2022-06-30T00:00:00.000000",
		FunArgs:   []interface{}{"sync_all"},
		Fun:       "schedule.run_job",
		Jid:       "20220630000000000000",
		Minions:   []string{"localhost"},
		Tgt:       "localhost",
		TgtType:   "glob",
		User:      "salt_user",
	}

	fakeBody, err := msgpack.Marshal(fake)
	if err != nil {
		log.Fatalln(err)
	}

	fakeMessage := []byte("salt/job/20220630000000000000/new\n\n")
	fakeMessage = append(fakeMessage, fakeBody...)

	return fakeMessage
}

/*
	Fake ack of manual triggered schedule job

	salt/job/20220630000000000000/ret/localhost	{
		"_stamp": "2022-06-30T00:00:00.000000",
		"cmd": "_return",
		"fun": "schedule.run_job",
		"fun_args": [
			"sync_all"
		],
		"id": "localhost",
		"jid": "20220630000000000000",
		"retcode": 0,
		"return": {
			"comment": "Scheduling Job sync_all on minion.",
			"result": true
		},
		"success": true
	}
*/

var expectedAckScheduleJob = SaltEvent{
	Tag:          "salt/job/20220630000000000000/ret/localhost",
	Type:         "ret",
	TargetNumber: 0,
	Data: EventData{
		Timestamp: "2022-06-30T00:00:00.000000",
		Cmd:       "_return",
		Fun:       "schedule.run_job",
		FunArgs:   []interface{}{"sync_all"},
		Id:        "localhost",
		Jid:       "20220630000000000000",
		Retcode:   0,
		Return: map[string]interface{}{
			"comment": "Scheduling Job sync_all on minion.",
			"result":  true,
		},
		Success: true,
	},
	IsScheduleJob: false,
}

func fakeAckScheduleJobEvent() []byte {
	// Marshal the data using MsgPack
	fake := FakeData{
		Timestamp: "2022-06-30T00:00:00.000000",
		Cmd:       "_return",
		Fun:       "schedule.run_job",
		FunArgs:   []interface{}{"sync_all"},
		Id:        "localhost",
		Jid:       "20220630000000000000",
		Retcode:   0,
		Return: map[string]interface{}{
			"comment": "Scheduling Job sync_all on minion.",
			"result":  true,
		},
		Success: true,
	}

	fakeBody, err := msgpack.Marshal(fake)
	if err != nil {
		log.Fatalln(err)
	}

	fakeMessage := []byte("salt/job/20220630000000000000/ret/localhost\n\n")
	fakeMessage = append(fakeMessage, fakeBody...)

	return fakeMessage
}

/*
	Fake schedule job return

	salt/job/20220630000000000000/ret/localhost	{
		"_stamp": "2022-06-30T00:00:00.000000",
		"arg": [],
		"cmd": "_return",
		"fun": "saltutil.sync_all",
		"fun_args": [],
		"id": "localhost",
		"jid": "20220630000000000000",
		"pid": 3969911,
		"retcode": 0,
		"return": {
			"beacons": [],
			"clouds": [],
			"engines": [],
			"executors": [],
			"grains": [],
			"log_handlers": [],
			"matchers": [],
			"modules": [],
			"output": [],
			"proxymodules": [],
			"renderers": [],
			"returners": [],
			"sdb": [],
			"serializers": [],
			"states": [],
			"thorium": [],
			"utils": []
		},
		"schedule": "sync_all",
		"success": true,
		"tgt": "localhost",
		"tgt_type": "glob"
	}
*/

var expectedScheduleJobReturn = SaltEvent{
	Tag:          "salt/job/20220630000000000000/ret/localhost",
	Type:         "ret",
	TargetNumber: 0,
	Data: EventData{
		Timestamp: "2022-06-30T00:00:00.000000",
		Cmd:       "_return",
		Fun:       "saltutil.sync_all",
		Id:        "localhost",
		Jid:       "20220630000000000000",
		Retcode:   0,
		Return: map[string]interface{}{
			"beacons":      []interface{}{},
			"clouds":       []interface{}{},
			"engines":      []interface{}{},
			"executors":    []interface{}{},
			"grains":       []interface{}{},
			"log_handlers": []interface{}{},
			"matchers":     []interface{}{},
			"modules":      []interface{}{},
			"output":       []interface{}{},
			"proxymodules": []interface{}{},
			"renderers":    []interface{}{},
			"returners":    []interface{}{},
			"sdb":          []interface{}{},
			"serializers":  []interface{}{},
			"states":       []interface{}{},
			"thorium":      []interface{}{},
			"utils":        []interface{}{},
		},
		Schedule: "sync_all",
		Success:  true,
		Tgt:      "localhost",
		TgtType:  "glob",
	},
	IsScheduleJob: true,
}

func fakeScheduleJobReturnEvent() []byte {
	// Marshal the data using MsgPack
	fake := FakeData{
		Timestamp: "2022-06-30T00:00:00.000000",
		Cmd:       "_return",
		Fun:       "saltutil.sync_all",
		Id:        "localhost",
		Jid:       "20220630000000000000",
		Retcode:   0,
		Return: map[string]interface{}{
			"beacons":      []interface{}{},
			"clouds":       []interface{}{},
			"engines":      []interface{}{},
			"executors":    []interface{}{},
			"grains":       []interface{}{},
			"log_handlers": []interface{}{},
			"matchers":     []interface{}{},
			"modules":      []interface{}{},
			"output":       []interface{}{},
			"proxymodules": []interface{}{},
			"renderers":    []interface{}{},
			"returners":    []interface{}{},
			"sdb":          []interface{}{},
			"serializers":  []interface{}{},
			"states":       []interface{}{},
			"thorium":      []interface{}{},
			"utils":        []interface{}{},
		},
		Schedule: "sync_all",
		Success:  true,
		Tgt:      "localhost",
		TgtType:  "glob",
	}

	fakeBody, err := msgpack.Marshal(fake)
	if err != nil {
		log.Fatalln(err)
	}

	fakeMessage := []byte("salt/job/20220630000000000000/ret/localhost\n\n")
	fakeMessage = append(fakeMessage, fakeBody...)

	return fakeMessage
}
