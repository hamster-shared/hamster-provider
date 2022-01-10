package events

type Event struct {
	Tag  OperationTag
	Data string
}

type OperationTag int

const OPCreated OperationTag = 1
const OPUpdated OperationTag = 2
const OPDeleted OperationTag = 3
