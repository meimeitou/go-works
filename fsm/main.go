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
	Events    map[string]EventSet
	FSM       *fsm.FSM
}

type EventSet struct {
	Name     string
	Messages []string
	Next     []string
	Params   string
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
	f.BeginWith = f.RawRule.BeginWith
	f.EndWith = f.RawRule.EndWith
	f.Events = make(map[string]EventSet, 0)
	for _, item := range f.RawRule.Events {
		f.Events[item.Name] = EventSet{
			Name:     item.Name,
			Messages: item.Messages,
			Next:     item.Next,
		}
	}
	call := fsm.Callbacks{
		"enter_state": func(e *fsm.Event) { f.enterState(e) },
		"leave_state": func(e *fsm.Event) {
			e.Async()
		},
	}
	events := []fsm.EventDesc{}
	for _, item1 := range f.RawRule.Events {
		for _, item2 := range item1.Next {
			events = append(events, fsm.EventDesc{
				Name: fmt.Sprintf("%s-%s", item1.Name, item2),
				Src:  []string{item1.Name},
				Dst:  item2,
			})
		}
	}
	f.FSM = fsm.NewFSM(
		f.BeginWith,
		events,
		call,
	)
}

func (f *FlexTask) Next() error {
	fmt.Println("当前消息", f.Events[f.FSM.Current()].Messages, "可选事件", f.FSM.AvailableTransitions())
	if len(f.FSM.AvailableTransitions()) < 2 {
		err := f.FSM.Event(f.FSM.AvailableTransitions()[0])
		if err != nil {
			fmt.Println(err)
		}
		err = f.FSM.Transition()
		if err != nil {
			fmt.Println(err)
		}
		return nil
	}
	var choose string
	fmt.Scanf("%s", &choose)
	err := f.FSM.Event(choose)
	if err != nil {
		fmt.Println(err)
	}
	err = f.FSM.Transition()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func NewFlexTask() *FlexTask {
	d := &FlexTask{
		RawRule: LoadRule("rule.yml"),
	}
	d.InitFsmFromRawRule()
	// d.FSM = fsm.NewFSM(
	// 	"open",
	// 	fsm.Events{
	// 		{Name: "open-level2", Src: []string{"open"}, Dst: "level2"},
	// 		{Name: "open-level3", Src: []string{"open"}, Dst: "level3"},
	// 		{Name: "level2-level4", Src: []string{"level2"}, Dst: "level4"},
	// 		{Name: "level3-level4", Src: []string{"level3"}, Dst: "level4"},
	// 		{Name: "level4-close", Src: []string{"level4"}, Dst: "close"},
	// 	},
	// 	fsm.Callbacks{
	// 		"enter_state": func(e *fsm.Event) { d.enterState(e) },
	// 		"leave_state": func(e *fsm.Event) {
	// 			e.Async()
	// 		},
	// 	},
	// )
	return d
}

func (f *FlexTask) enterState(e *fsm.Event) {
	fmt.Printf("messages : %v", f.Events[f.FSM.Current()])
	fmt.Printf("enter state, des: %s, src :  %v\n", e.Dst, e.Src)
}

func main() {
	door := NewFlexTask()
	door.Next()
	door.Next()
	//fmt.Println(door)

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
