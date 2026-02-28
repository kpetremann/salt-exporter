package parser_test

import (
	"log"
	"time"

	"github.com/kpetremann/salt-exporter/pkg/event"
	"github.com/vmihailenco/msgpack/v5"
)

/*
	Fake state.highstate ret with saltenv/pillarenv args

	salt/job/20260218151902735416/ret/test-node-00	{
		"_stamp": "2026-02-18T15:19:02.735941",
		"arg": [
			"saltenv=branch_name",
			"pillarenv=branch_name"
		],
		"cmd": "_return",
		"fun": "state.highstate",
		"fun_args": [
			"saltenv=branch_name",
			"pillarenv=branch_name"
		],
		"id": "test-node-00",
		"jid": "20260218151902735416",
		"out": "highstate",
		"retcode": 0,
		"return": {
			"file_|-hostname_file_|-/etc/hostname_|-managed": {
				"__id__": "hostname_file",
				"__run_num__": 0,
				"__sls__": "defaults",
				"changes": {},
				"comment": "File /etc/hostname is in the correct state",
				"duration": 14.258,
				"name": "/etc/hostname",
				"result": true,
				"start_time": "15:19:02.712934"
			}
		},
		"tgt": "test-node-00",
		"tgt_type": "glob"
	}
*/

var expectedStateHighstateWithEnvReturn = event.SaltEvent{
	Tag:          "salt/job/20260218151902735416/ret/test-node-00",
	Type:         "ret",
	Module:       event.JobModule,
	TargetNumber: 0,
	Data: event.EventData{
		Timestamp: "2026-02-18T15:19:02.735941",
		Arg:       []any{"saltenv=branch_name", "pillarenv=branch_name"},
		Cmd:       "_return",
		Fun:       "state.highstate",
		FunArgs:   []any{"saltenv=branch_name", "pillarenv=branch_name"},
		ID:        "test-node-00",
		Jid:       "20260218151902735416",
		Out:       "highstate",
		Retcode:   0,
		Return: map[string]any{
			"file_|-hostname_file_|-/etc/hostname_|-managed": map[string]any{
				"__id__":      "hostname_file",
				"__run_num__": int8(0),
				"__sls__":     "defaults",
				"changes":     map[string]any{},
				"comment":     "File /etc/hostname is in the correct state",
				"duration":    14.258,
				"name":        "/etc/hostname",
				"result":      true,
				"start_time":  "15:19:02.712934",
			},
		},
		Tgt:     "test-node-00",
		TgtType: "glob",
	},
	IsScheduleJob:      false,
	IsTest:             false,
	IsMock:             false,
	StateModuleSuccess: new(true),
	StateDuration:      new(time.Duration(14.258 * float64(time.Second))),
}

func fakeStateHighstateWithEnvReturnEvent() []byte {
	fake := FakeData{
		Timestamp: "2026-02-18T15:19:02.735941",
		Arg:       []any{"saltenv=branch_name", "pillarenv=branch_name"},
		Cmd:       "_return",
		Fun:       "state.highstate",
		FunArgs:   []any{"saltenv=branch_name", "pillarenv=branch_name"},
		ID:        "test-node-00",
		Jid:       "20260218151902735416",
		Out:       "highstate",
		Retcode:   0,
		Return: map[string]any{
			"file_|-hostname_file_|-/etc/hostname_|-managed": map[string]any{
				"__id__":      "hostname_file",
				"__run_num__": 0,
				"__sls__":     "defaults",
				"changes":     map[string]any{},
				"comment":     "File /etc/hostname is in the correct state",
				"duration":    14.258,
				"name":        "/etc/hostname",
				"result":      true,
				"start_time":  "15:19:02.712934",
			},
		},
		Tgt:     "test-node-00",
		TgtType: "glob",
	}

	fakeBody, err := msgpack.Marshal(fake)
	if err != nil {
		log.Fatalln(err)
	}

	fakeMessage := []byte("salt/job/20260218151902735416/ret/test-node-00\n\n")
	fakeMessage = append(fakeMessage, fakeBody...)

	return fakeMessage
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

var expectedNewStateSlsJob = event.SaltEvent{
	Tag:          "salt/job/20220630000000000000/new",
	Type:         "new",
	Module:       event.JobModule,
	TargetNumber: 1,
	Data: event.EventData{
		Timestamp: "2022-06-30T00:00:00.000000",
		Fun:       "state.sls",
		Arg:       []any{"test"},
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
		Arg:       []any{"test"},
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
		FunArgs:   []any{"test"},
		ID:        "node1",
		Jid:       "20220630000000000000",
		Out:       "highstate",
		Retcode:   0,
		Return: map[string]any{
			"test_|-dummy test_|-Dummy test_|-nop": map[string]any{
				"__id__":      "dummy test",
				"__run_num__": int8(0),
				"__sls__":     "test",
				"changes":     map[string]any{},
				"comment":     "Success!",
				"duration":    0.481,
				"name":        "Dummy test",
				"result":      true,
				"start_time":  "09:17:08.822722",
			},
		},
		Success: new(true),
	},
	IsScheduleJob:      false,
	IsTest:             false,
	IsMock:             false,
	StateModuleSuccess: new(true),
	StateDuration:      new(time.Duration(0.481 * float64(time.Second))),
}

func fakeStateSlsReturnEvent() []byte {
	// Marshal the data using MsgPack
	fake := FakeData{
		Timestamp: "2022-06-30T00:00:00.000000",
		Cmd:       "_return",
		Fun:       "state.sls",
		FunArgs:   []any{"test"},
		ID:        "node1",
		Out:       "highstate",
		Jid:       "20220630000000000000",
		Retcode:   0,
		Return: map[string]any{
			"test_|-dummy test_|-Dummy test_|-nop": map[string]any{
				"__id__":      "dummy test",
				"__run_num__": 0,
				"__sls__":     "test",
				"changes":     map[string]any{},
				"comment":     "Success!",
				"duration":    0.481,
				"name":        "Dummy test",
				"result":      true,
				"start_time":  "09:17:08.822722",
			},
		},
		Success: new(true),
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
		Arg: []any{
			map[string]any{
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
		Arg: []any{
			map[string]any{
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
		FunArgs: []any{
			map[string]any{
				"fun":  "test.nop",
				"name": "toto",
			},
		},
		ID:      "node1",
		Jid:     "20220630000000000000",
		Out:     "highstate",
		Retcode: 0,
		Return: map[string]any{
			"test_|-toto_|-toto_|-nop": map[string]any{
				"__id__":      "toto",
				"__run_num__": int8(0),
				"__sls__":     nil,
				"changes":     map[string]any{},
				"comment":     "Success!",
				"duration":    0.49,
				"name":        "toto",
				"result":      true,
				"start_time":  "09:20:38.462572",
			},
		},
		Success: new(true),
	},
	IsScheduleJob:      false,
	IsTest:             false,
	IsMock:             false,
	StateModuleSuccess: new(true),
	StateDuration:      new(time.Duration(0.49 * float64(time.Second))),
}

func fakeStateSingleReturnEvent() []byte {
	// Marshal the data using MsgPack
	fake := FakeData{
		Timestamp: "2022-06-30T00:00:00.000000",
		Cmd:       "_return",
		Fun:       "state.single",
		FunArgs: []any{
			map[string]any{
				"fun":  "test.nop",
				"name": "toto",
			},
		},
		ID:      "node1",
		Jid:     "20220630000000000000",
		Out:     "highstate",
		Retcode: 0,
		Return: map[string]any{
			"test_|-toto_|-toto_|-nop": map[string]any{
				"__id__":      "toto",
				"__run_num__": 0,
				"__sls__":     nil,
				"changes":     map[string]any{},
				"comment":     "Success!",
				"duration":    0.49,
				"name":        "toto",
				"result":      true,
				"start_time":  "09:20:38.462572",
			},
		},
		Success: new(true),
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
		Arg: []any{
			"somestate",
			map[string]any{
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
		Arg: []any{
			"somestate",
			map[string]any{
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
		FunArgs: []any{
			"somestate",
			map[string]any{
				"test": true,
				"mock": true,
			},
		},
		ID:      "node1",
		Jid:     "20220630000000000000",
		Out:     "highstate",
		Retcode: 1,
		Return: map[string]any{
			"somestate_|-dummy somestate_|-Dummy somestate_|-nop": map[string]any{
				"__id__":      "dummy somestate",
				"__run_num__": int8(0),
				"__sls__":     "somestate",
				"changes":     map[string]any{},
				"comment":     "Success!",
				"duration":    0.481,
				"name":        "Dummy somestate",
				"result":      true,
				"start_time":  "09:17:08.822722",
			},
			"somestate_|-failed_|-failed_|-fail_with_changes": map[string]any{
				"__id__":      "dummy somestate",
				"__run_num__": int8(2),
				"__sls__":     "somestate",
				"changes":     map[string]any{},
				"comment":     "Failure!",
				"duration":    0.579,
				"name":        "failed",
				"result":      false,
				"start_time":  "09:17:08.812345",
			},
		},
		Success: new(true),
	},
	IsScheduleJob:      false,
	IsTest:             true,
	IsMock:             true,
	StateModuleSuccess: new(false),
	StateDuration:      new(time.Duration((0.481 + 0.579) * float64(time.Second))),
}

func fakeTestMockStateSlsReturnEvent() []byte {
	// Marshal the data using MsgPack
	fake := FakeData{
		Timestamp: "2022-06-30T00:00:00.000000",
		Cmd:       "_return",
		Fun:       "state.sls",
		FunArgs: []any{
			"somestate",
			map[string]any{
				"test": true,
				"mock": true,
			},
		},
		ID:      "node1",
		Out:     "highstate",
		Jid:     "20220630000000000000",
		Retcode: 1,
		Return: map[string]any{
			"somestate_|-dummy somestate_|-Dummy somestate_|-nop": map[string]any{
				"__id__":      "dummy somestate",
				"__run_num__": 0,
				"__sls__":     "somestate",
				"changes":     map[string]any{},
				"comment":     "Success!",
				"duration":    0.481,
				"name":        "Dummy somestate",
				"result":      true,
				"start_time":  "09:17:08.822722",
			},
			"somestate_|-failed_|-failed_|-fail_with_changes": map[string]any{
				"__id__":      "dummy somestate",
				"__run_num__": int8(2),
				"__sls__":     "somestate",
				"changes":     map[string]any{},
				"comment":     "Failure!",
				"duration":    0.579,
				"name":        "failed",
				"result":      false,
				"start_time":  "09:17:08.812345",
			},
		},
		Success: new(true),
	}

	fakeBody, err := msgpack.Marshal(fake)
	if err != nil {
		log.Fatalln(err)
	}

	fakeMessage := []byte("salt/job/20220630000000000000/ret/node1\n\n")
	fakeMessage = append(fakeMessage, fakeBody...)

	return fakeMessage
}
