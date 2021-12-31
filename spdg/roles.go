package spdg

/*
ContextRole wraps the role field into an enum.
*/
type ContextRole string

const (
	// NULLGRAM is used when the return value is void basically.
	NULLGRAM ContextRole = "null"
	// BASEGRAM is used with a minimally (incomplete) initialized datagram.
	BASEGRAM ContextRole = "base"
	// ANONYMOUS is used as a `generic` type.
	ANONYMOUS ContextRole = "anonymous"
	// ERROR is used as a `generic` error type.
	ERROR ContextRole = "error"
	// DATAPOINT is usually something you want to store somehow.
	DATAPOINT ContextRole = "datapoint"
	// COMMAND is to trigger processes in the data pipelines.
	COMMAND ContextRole = "command"
	// QUESTION is to search for data.
	QUESTION ContextRole = "question"
	// TOPIC is a key that targets data at a Question.
	TOPIC ContextRole = "topic"
)
