package events

import (
	"log"

	"github.com/vmihailenco/msgpack/v5"
)

func getNewStateEvent() SaltEvent {
	return SaltEvent{
		Tag:          "salt/job/20220630000000000000/new",
		Type:         "new",
		TargetNumber: 1,
		Data: EventData{
			Timestamp: "2022-06-30T00:00:00.000000",
			Fun:       "state.sls",
			Arg:       []interface{}{"test"},
			Jid:       "20220630000000000000",
			Minions:   []string{"node1"},
			Missing:   []string{},
			Tgt:       "node1",
			TgtType:   "glob",
			User:      "salt_user",
		},
		IsScheduleJob: false,
	}
}

/*
	Fake state.sls job

	salt/job/20220630000000000000/new	{
		"_stamp": "2022-06-30T00:00:00.000000",
		"arg": [
			"test"
		],
		"fun": "state.sls",
		"jid": "20220630000000000000",
		"minions": [
			"node1"
		],
		"missing": [],
		"tgt": "node1",
		"tgt_type": "glob",
		"user": "salt_user"
	}
*/

var expectedNewStateSlsJob = SaltEvent{
	Tag:          "salt/job/20220630000000000000/new",
	Type:         "new",
	TargetNumber: 1,
	Data: EventData{
		Timestamp: "2022-06-30T00:00:00.000000",
		Fun:       "state.sls",
		Arg:       []interface{}{"test"},
		Jid:       "20220630000000000000",
		Minions:   []string{"node1"},
		Missing:   []string{},
		Tgt:       "node1",
		TgtType:   "glob",
		User:      "salt_user",
	},
	IsScheduleJob: false,
}

func fakeNewStateSlsJobEvent() []byte {
	// Marshal the data using MsgPack
	fake := FakeData{
		Timestamp: "2022-06-30T00:00:00.000000",
		Fun:       "state.sls",
		Arg:       []interface{}{"test"},
		Jid:       "20220630000000000000",
		Minions:   []string{"node1"},
		Missing:   []string{},
		Tgt:       "node1",
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
	Fake state.sls ret

	salt/job/20220630000000000000/ret/node1	{
		"_stamp": "2022-06-30T00:00:00.000000",
		"cmd": "_return",
		"fun": "state.sls",
		"fun_args": [
			"test"
		],
		"id": "node1",
		"jid": "20220630000000000000",
		"out": "highstate",
		"retcode": 0,
		"return": {
			"test_|-dummy test_|-Dummy test_|-nop": {
				"__id__": "dummy test",
				"__run_num__": 0,
				"__sls__": "test",
				"changes": {},
				"comment": "Success!",
				"duration": 0.481,
				"name": "Dummy test",
				"result": true,
				"start_time": "09:17:08.822722"
			}
		},
		"success": true
	}

*/

var expectedStateSlsReturn = SaltEvent{
	Tag:          "salt/job/20220630000000000000/ret/node1",
	Type:         "ret",
	TargetNumber: 0,
	Data: EventData{
		Timestamp: "2022-06-30T00:00:00.000000",
		Cmd:       "_return",
		Fun:       "state.sls",
		FunArgs:   []interface{}{"test"},
		Id:        "node1",
		Jid:       "20220630000000000000",
		Out:       "highstate",
		Retcode:   0,
		Return: map[string]interface{}{
			"test_|-dummy test_|-Dummy test_|-nop": map[string]interface{}{
				"__id__":      "dummy test",
				"__run_num__": int8(0),
				"__sls__":     "test",
				"changes":     map[string]interface{}{},
				"comment":     "Success!",
				"duration":    0.481,
				"name":        "Dummy test",
				"result":      true,
				"start_time":  "09:17:08.822722",
			},
		},
		Success: true,
	},
	IsScheduleJob: false,
}

func fakeStateSlsReturnEvent() []byte {
	// Marshal the data using MsgPack
	fake := FakeData{
		Timestamp: "2022-06-30T00:00:00.000000",
		Cmd:       "_return",
		Fun:       "state.sls",
		FunArgs:   []interface{}{"test"},
		Id:        "node1",
		Out:       "highstate",
		Jid:       "20220630000000000000",
		Retcode:   0,
		Return: map[string]interface{}{
			"test_|-dummy test_|-Dummy test_|-nop": map[string]interface{}{
				"__id__":      "dummy test",
				"__run_num__": 0,
				"__sls__":     "test",
				"changes":     map[string]interface{}{},
				"comment":     "Success!",
				"duration":    0.481,
				"name":        "Dummy test",
				"result":      true,
				"start_time":  "09:17:08.822722",
			},
		},
		Success: true,
	}

	fakeBody, err := msgpack.Marshal(fake)
	if err != nil {
		log.Fatalln(err)
	}

	fakeMessage := []byte("salt/job/20220630000000000000/ret/node1\n\n")
	fakeMessage = append(fakeMessage, fakeBody...)

	return fakeMessage
}

/*

	Fake state.single

	salt/job/20220630000000000000/new	{
		"_stamp": "2022-06-30T00:00:00.000000",
		"arg": [
			{
				"__kwarg__": true,
				"fun": "test.nop",
				"name": "toto"
			}
		],
		"fun": "state.single",
		"jid": "20220630000000000000",
		"minions": [
			"node1"
		],
		"missing": [],
		"tgt": "node1",
		"tgt_type": "glob",
		"user": "salt_user"
	}

*/

var expectedNewStateSingle = SaltEvent{
	Tag:          "salt/job/20220630000000000000/new",
	Type:         "new",
	TargetNumber: 1,
	Data: EventData{
		Timestamp: "2022-06-30T00:00:00.000000",
		Arg: []interface{}{
			map[string]interface{}{
				"__kwarg__": true,
				"fun":       "test.nop",
				"name":      "toto",
			},
		},
		Fun:     "state.single",
		Jid:     "20220630000000000000",
		Minions: []string{"node1"},
		Missing: []string{},
		Tgt:     "node1",
		TgtType: "glob",
		User:    "salt_user",
	},
	IsScheduleJob: false,
}

func fakeNewStateSingleEvent() []byte {
	// Marshal the data using MsgPack
	fake := FakeData{
		Timestamp: "2022-06-30T00:00:00.000000",
		Arg: []interface{}{
			map[string]interface{}{
				"__kwarg__": true,
				"fun":       "test.nop",
				"name":      "toto",
			},
		},
		Fun:     "state.single",
		Jid:     "20220630000000000000",
		Minions: []string{"node1"},
		Missing: []string{},
		Tgt:     "node1",
		TgtType: "glob",
		User:    "salt_user",
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

	Fake state.single return

	salt/job/20220630000000000000/ret/node1	{
		"_stamp": "2022-06-30T00:00:00.000000",
		"cmd": "_return",
		"fun": "state.single",
		"fun_args": [
			{
				"fun": "test.nop",
				"name": "toto"
			}
		],
		"id": "node1",
		"jid": "20220630000000000000",
		"out": "highstate",
		"retcode": 0,
		"return": {
			"test_|-toto_|-toto_|-nop": {
				"__id__": "toto",
				"__run_num__": 0,
				"__sls__": null,
				"changes": {},
				"comment": "Success!",
				"duration": 0.49,
				"name": "toto",
				"result": true,
				"start_time": "09:20:38.462572"
			}
		},
		"success": true
	}

*/

var expectedStateSingleReturn = SaltEvent{
	Tag:          "salt/job/20220630000000000000/ret/node1",
	Type:         "ret",
	TargetNumber: 0,
	Data: EventData{
		Timestamp: "2022-06-30T00:00:00.000000",
		Cmd:       "_return",
		Fun:       "state.single",
		FunArgs: []interface{}{
			map[string]interface{}{
				"fun":  "test.nop",
				"name": "toto",
			},
		},
		Id:      "node1",
		Jid:     "20220630000000000000",
		Out:     "highstate",
		Retcode: 0,
		Return: map[string]interface{}{
			"test_|-toto_|-toto_|-nop": map[string]interface{}{
				"__id__":      "toto",
				"__run_num__": int8(0),
				"__sls__":     nil,
				"changes":     map[string]interface{}{},
				"comment":     "Success!",
				"duration":    0.49,
				"name":        "toto",
				"result":      true,
				"start_time":  "09:20:38.462572",
			},
		},
		Success: true,
	},
	IsScheduleJob: false,
}

func fakeStateSingleReturnEvent() []byte {
	// Marshal the data using MsgPack
	fake := FakeData{
		Timestamp: "2022-06-30T00:00:00.000000",
		Cmd:       "_return",
		Fun:       "state.single",
		FunArgs: []interface{}{
			map[string]interface{}{
				"fun":  "test.nop",
				"name": "toto",
			},
		},
		Id:      "node1",
		Jid:     "20220630000000000000",
		Out:     "highstate",
		Retcode: 0,
		Return: map[string]interface{}{
			"test_|-toto_|-toto_|-nop": map[string]interface{}{
				"__id__":      "toto",
				"__run_num__": 0,
				"__sls__":     nil,
				"changes":     map[string]interface{}{},
				"comment":     "Success!",
				"duration":    0.49,
				"name":        "toto",
				"result":      true,
				"start_time":  "09:20:38.462572",
			},
		},
		Success: true,
	}

	fakeBody, err := msgpack.Marshal(fake)
	if err != nil {
		log.Fatalln(err)
	}

	fakeMessage := []byte("salt/job/20220630000000000000/ret/node1\n\n")
	fakeMessage = append(fakeMessage, fakeBody...)

	return fakeMessage
}
