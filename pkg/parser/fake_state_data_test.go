package parser_test

import (
	"log"

	"github.com/kpetremann/salt-exporter/pkg/event"
	"github.com/vmihailenco/msgpack/v5"
)

var False = false
var True = true

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

var expectedNewStateSlsJob = event.SaltEvent{
	Tag:          "salt/job/20220630000000000000/new",
	Type:         "new",
	Module:       event.JobModule,
	TargetNumber: 1,
	Data: event.EventData{
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
	IsTest:        false,
	IsMock:        false,
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

var expectedStateSlsReturn = event.SaltEvent{
	Tag:          "salt/job/20220630000000000000/ret/node1",
	Type:         "ret",
	Module:       event.JobModule,
	TargetNumber: 0,
	Data: event.EventData{
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
	IsScheduleJob:      false,
	IsTest:             false,
	IsMock:             false,
	StateModuleSuccess: &True,
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

var expectedNewStateSingle = event.SaltEvent{
	Tag:          "salt/job/20220630000000000000/new",
	Type:         "new",
	Module:       event.JobModule,
	TargetNumber: 1,
	Data: event.EventData{
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
	IsTest:        false,
	IsMock:        false,
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

var expectedStateSingleReturn = event.SaltEvent{
	Tag:          "salt/job/20220630000000000000/ret/node1",
	Type:         "ret",
	Module:       event.JobModule,
	TargetNumber: 0,
	Data: event.EventData{
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
	IsScheduleJob:      false,
	IsTest:             false,
	IsMock:             false,
	StateModuleSuccess: &True,
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

/*
	Fake state.sls job test=True mock=True

	salt/job/20220630000000000000/new	{
		"_stamp": "2022-06-30T00:00:00.000000",
		"fun": "state.sls",
		"arg": [
			"somestate",
			{
				"__kwarg__": true,
				"test": true,
				"mock": true
			}
		],
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

var expectedNewTestMockStateSlsJob = event.SaltEvent{
	Tag:          "salt/job/20220630000000000000/new",
	Type:         "new",
	Module:       event.JobModule,
	TargetNumber: 1,
	Data: event.EventData{
		Timestamp: "2022-06-30T00:00:00.000000",
		Fun:       "state.sls",
		Arg: []interface{}{
			"somestate",
			map[string]interface{}{
				"test": true,
				"mock": true,
			},
		},
		Jid:     "20220630000000000000",
		Minions: []string{"node1"},
		Missing: []string{},
		Tgt:     "node1",
		TgtType: "glob",
		User:    "salt_user",
	},
	IsScheduleJob: false,
	IsTest:        true,
	IsMock:        true,
}

func fakeNewTestMockStateSlsJobEvent() []byte {
	// Marshal the data using MsgPack
	fake := FakeData{
		Timestamp: "2022-06-30T00:00:00.000000",
		Fun:       "state.sls",
		Arg: []interface{}{
			"somestate",
			map[string]interface{}{
				"test": true,
				"mock": true,
			},
		},
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
	Fake state.sls ret

	salt/job/20220630000000000000/ret/node1	{
		"_stamp": "2022-06-30T00:00:00.000000",
		"cmd": "_return",
		"fun": "state.sls",
		"fun_args": [
			"somestate",
			{
				"test": true,
				"mock": true
			}
		],
		"id": "node1",
		"jid": "20220630000000000000",
		"out": "highstate",
		"retcode": 1,
		"return": {
			"somestate_|-dummy somestate_|-Dummy somestate_|-nop": {
				"__id__": "dummy somestate",
				"__run_num__": 0,
				"__sls__": "somestate",
				"changes": {},
				"comment": "Success!",
				"duration": 0.481,
				"name": "Dummy somestate",
				"result": true,
				"start_time": "09:17:08.822722"
			},
			"somestate_|-failed_|-failed_|-fail_with_changes": {
				"__id__": "failed",
				"__run_num__": 2,
				"__sls__": "test",
				"changes": {
					"testing": {
						"new": "Something pretended to change",
						"old": "Unchanged"
					}
				},
				"comment": "Failure!",
				"duration": 0.579,
				"name": "failed",
				"result": false,
				"start_time": "09:17:02.812345"
			},
		},
		"success": true
	}

*/

var expectedTestMockStateSlsReturn = event.SaltEvent{
	Tag:          "salt/job/20220630000000000000/ret/node1",
	Type:         "ret",
	Module:       event.JobModule,
	TargetNumber: 0,
	Data: event.EventData{
		Timestamp: "2022-06-30T00:00:00.000000",
		Cmd:       "_return",
		Fun:       "state.sls",
		FunArgs: []interface{}{
			"somestate",
			map[string]interface{}{
				"test": true,
				"mock": true,
			},
		},
		Id:      "node1",
		Jid:     "20220630000000000000",
		Out:     "highstate",
		Retcode: 1,
		Return: map[string]interface{}{
			"somestate_|-dummy somestate_|-Dummy somestate_|-nop": map[string]interface{}{
				"__id__":      "dummy somestate",
				"__run_num__": int8(0),
				"__sls__":     "somestate",
				"changes":     map[string]interface{}{},
				"comment":     "Success!",
				"duration":    0.481,
				"name":        "Dummy somestate",
				"result":      true,
				"start_time":  "09:17:08.822722",
			},
			"somestate_|-failed_|-failed_|-fail_with_changes": map[string]interface{}{
				"__id__":      "dummy somestate",
				"__run_num__": int8(2),
				"__sls__":     "somestate",
				"changes":     map[string]interface{}{},
				"comment":     "Failure!",
				"duration":    0.579,
				"name":        "failed",
				"result":      false,
				"start_time":  "09:17:08.812345",
			},
		},
		Success: true,
	},
	IsScheduleJob:      false,
	IsTest:             true,
	IsMock:             true,
	StateModuleSuccess: &False,
}

func fakeTestMockStateSlsReturnEvent() []byte {
	// Marshal the data using MsgPack
	fake := FakeData{
		Timestamp: "2022-06-30T00:00:00.000000",
		Cmd:       "_return",
		Fun:       "state.sls",
		FunArgs: []interface{}{
			"somestate",
			map[string]interface{}{
				"test": true,
				"mock": true,
			},
		},
		Id:      "node1",
		Out:     "highstate",
		Jid:     "20220630000000000000",
		Retcode: 1,
		Return: map[string]interface{}{
			"somestate_|-dummy somestate_|-Dummy somestate_|-nop": map[string]interface{}{
				"__id__":      "dummy somestate",
				"__run_num__": 0,
				"__sls__":     "somestate",
				"changes":     map[string]interface{}{},
				"comment":     "Success!",
				"duration":    0.481,
				"name":        "Dummy somestate",
				"result":      true,
				"start_time":  "09:17:08.822722",
			},
			"somestate_|-failed_|-failed_|-fail_with_changes": map[string]interface{}{
				"__id__":      "dummy somestate",
				"__run_num__": int8(2),
				"__sls__":     "somestate",
				"changes":     map[string]interface{}{},
				"comment":     "Failure!",
				"duration":    0.579,
				"name":        "failed",
				"result":      false,
				"start_time":  "09:17:08.812345",
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
