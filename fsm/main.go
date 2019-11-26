package main

import (
	"fmt"
	"io/ioutil"

	"github.com/looplab/fsm"
	"gopkg.in/yaml.v2"
)

type FlexRule struct {
	BeginWith string      `json:"begin_with" yaml:"begin_with"`
	EndWith   string      `json:"end_with" yaml:"end_with"`
	Events    []TaskEvent `json:"events" yaml:"events"`
}

type TaskEvent struct {
	Name     string
	Messages []string
	Next     []string
}

type FlexTask struct {
	BeginWith string
	EndWith   string
	Current   string
	RawRule   FlexRule
	Params    map[string]string
	FSM       *fsm.FSM
}

func LoadRule(path string) FlexRule {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic("error file")
	}
	var rule FlexRule
	err = yaml.Unmarshal(data, &rule)
	if err != nil {
		panic(err.Error())
	}
	return rule
}

func (f *FlexTask) InitFsmFromRawRule() {
	call := fsm.Callbacks{
		"enter_state": func(e *fsm.Event) { f.enterState(e) },
		"leave_state": func(e *fsm.Event) {
			e.Async()
		},
	}
	events := fsm.Events{}
	f.FSM = fsm.NewFSM(
		f.BeginWith,
		events,
		call,
	)
}

func NewFlexTask() *FlexTask {
	d := &FlexTask{
		BeginWith: "open",
		Current:   "open",
		EndWith:   "end",
		RawRule:   LoadRule("rule.yml"),
	}
	d.FSM = fsm.NewFSM(
		"open",
		fsm.Events{
			{Name: "open-level2", Src: []string{"open"}, Dst: "level2"},
			{Name: "open-level3", Src: []string{"open"}, Dst: "level3"},
			{Name: "level2-level4", Src: []string{"level2"}, Dst: "level4"},
			{Name: "level3-level4", Src: []string{"level3"}, Dst: "level4"},
			{Name: "level4-close", Src: []string{"level4"}, Dst: "close"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) { d.enterState(e) },
			"leave_state": func(e *fsm.Event) {
				e.Async()
			},
		},
	)
	fmt.Println(d)
	return d
}

func (f *FlexTask) enterState(e *fsm.Event) {
	fmt.Printf("enter state,    %v is %s, statue :  %v\n", f.Params, e.Dst, e)
}

func main() {
	var chose string
	door := NewFlexTask()
	fmt.Println(door.FSM.Current(), door.FSM.AvailableTransitions())
	fmt.Scanf("%s", &chose)
	err := door.FSM.Event(chose)
	if err != nil {
		fmt.Println(err)
	}
	err = door.FSM.Transition()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(door.FSM.Current(), door.FSM.AvailableTransitions())
	fmt.Scanf("%s", &chose)
	err = door.FSM.Event(chose)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("----")

	// fsmy := fsm.NewFSM(
	// 	"closed",
	// 	fsm.Events{
	// 		{Name: "openx", Src: []string{"closed"}, Dst: "open"},
	// 		{Name: "close", Src: []string{"open"}, Dst: "closed"},
	// 	},
	// 	fsm.Callbacks{
	// 		"leave_state": func(e *fsm.Event) {
	// 			e.Async()
	// 		},
	// 	},
	// )
	// fmt.Println(fsmy.Current())
	// err = fsmy.Event("openx")
	// if e, ok := err.(fsm.AsyncError); !ok && e.Err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(fsmy.Current())
	// // err = fsmy.Transition()
	// // if err != nil {
	// // 	fmt.Println(err)
	// // }
	// // fmt.Println(fsmy.Current())
	// err = fsmy.Event("close")
	// if e, ok := err.(fsm.AsyncError); !ok && e.Err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(fsmy.Current())
	// err = fsmy.Transition()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(fsmy.Current())
}
