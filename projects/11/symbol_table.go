package main

import "fmt"

type VarKind int

const (
	STATICV VarKind = iota
	FIELDV
	ARGV
	VARV
)

type JackType = string

type TableEntry struct {
	jackType JackType
	kind     VarKind
	index    int
}

var CurrentClass string

var classIndexCount = map[VarKind]int{FIELDV: 0, STATICV: 0}

var ClassVars map[string]TableEntry = make(map[string]TableEntry)

var routineIndexCount = map[VarKind]int{ARGV: 0, VARV: 0}

var RoutineVars map[string]TableEntry = make(map[string]TableEntry)

type VarNotFoundError struct{ Name string }

func (e *VarNotFoundError) Error() string {
	return fmt.Sprintf("variable not found: %s", e.Name)
}

func ResetClass() {
	classIndexCount[FIELDV] = 0
	classIndexCount[STATICV] = 0
	clear(ClassVars)
}

func ResetRoutine() {
	routineIndexCount[ARGV] = 0
	routineIndexCount[VARV] = 0
	clear(RoutineVars)
}

func Define(name string, jackType string, kind VarKind) int {
	if kind == STATICV || kind == FIELDV {
		c := classIndexCount[kind]
		ClassVars[name] = TableEntry{kind: kind, index: c, jackType: jackType}
		c++
		classIndexCount[kind] = c
		return c
	}
	c := routineIndexCount[kind]
	RoutineVars[name] = TableEntry{kind: kind, index: c, jackType: jackType}
	c++
	routineIndexCount[kind] = c
	return c
}

func DefineThis() {
	RoutineVars["this"] = TableEntry{kind: ARGV, index: 0, jackType: CurrentClass}
	routineIndexCount[ARGV] = 1
}

func VarCount(kind VarKind) int {
	if kind == STATICV || kind == FIELDV {
		return classIndexCount[kind]
	}
	return routineIndexCount[kind]
}

func KindOf(name string) (VarKind, error) {
	entry, ok := RoutineVars[name]
	if ok {
		return entry.kind, nil
	}
	entry, ok = ClassVars[name]
	if !ok {
		return 0, &VarNotFoundError{Name: name}
	}
	return entry.kind, nil
}

func IndexOf(name string) (int, error) {
	entry, ok := RoutineVars[name]
	if ok {
		return entry.index, nil
	}
	entry, ok = ClassVars[name]
	if !ok {
		return 0, &VarNotFoundError{Name: name}
	}
	return entry.index, nil
}

func TypeOf(name string) (string, error) {
	entry, ok := RoutineVars[name]
	if ok {
		return entry.jackType, nil
	}
	entry, ok = ClassVars[name]
	if !ok {
		return "", &VarNotFoundError{Name: name}
	}
	return entry.jackType, nil
}
