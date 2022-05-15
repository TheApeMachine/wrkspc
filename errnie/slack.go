package errnie

import (
	"github.com/acarl005/stripansi"
	slacker "github.com/slack-go/slack"
	"github.com/spf13/viper"
)

var healthy bool

/*
SlackLogger is an output channel that will log errors to Slack.
*/
type SlackLogger struct {
	client *slacker.Client
}

func ensureLogger() {
	if len(ambctx.loggers) == 1 {
		program := viper.GetViper().GetString("program")
		stage := viper.GetViper().GetString(program + ".stage")
		token := viper.GetViper().GetString(program + ".slack.token")

		if stage == "production" {
			ambctx.loggers = append(ambctx.loggers, NewLogger(&SlackLogger{
				client: slacker.New(token, slacker.OptionDebug(true)),
			}))
		}

		healthy = true
	}
}

func (logger SlackLogger) Info(events ...interface{})    {}
func (logger SlackLogger) Debug(events ...interface{})   {}
func (logger SlackLogger) Inspect(events ...interface{}) {}

func (logger SlackLogger) Warning(events ...interface{}) *Error {
	if len(events) == 0 {
		return nil
	}

	var errs []error

	for _, err := range events {
		if err == nil {
			break
		}

		errs = append(errs, err.(error))
	}

	if len(errs) == 0 {
		return nil
	}

	logger.postSlack(errs[0].Error(), Traces(true, true))
	return NewError(errs...)
}

func (logger SlackLogger) Error(events ...interface{}) *Error {
	if len(events) == 0 {
		return nil
	}

	var errs []error

	for _, err := range events {
		if err == nil {
			break
		}

		errs = append(errs, err.(error))
	}

	if len(errs) == 0 {
		return nil
	}

	logger.postSlack(errs[0].Error(), Traces(true, true))
	return NewError(errs...)
}

func (logger SlackLogger) postSlack(errStr string, trace string) {
	if !healthy {
		return
	}

	program := viper.GetString("program")
	cleanTrace := stripansi.Strip(trace)

	attachment := slacker.Attachment{
		Pretext: viper.GetViper().GetString("program"),
		Text:    "PRODUCTION ERROR",
		Color:   "#FF0055",
		Fields: []slacker.AttachmentField{
			{
				Title: "The following error was detected",
				Value: errStr,
			},
			{
				Title: "Stacktrace",
				Value: cleanTrace,
			},
		},
	}

	_, _, err := logger.client.PostMessage(
		viper.GetString(program+".slack.errorChannel"),
		slacker.MsgOptionAttachments(attachment),
	)

	if err := Handles(err).With(NOOP).ERR.First(); err != nil {
		healthy = false
	}
}
