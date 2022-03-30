package errnie

import (
	slacker "github.com/slack-go/slack"
	"github.com/spf13/viper"
)

/*
SlackLogger is an output channel that will log errors to Slack.
*/
type SlackLogger struct {
	client *slacker.Client
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

	logger.postSlack(errs[0].Error())
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

	logger.postSlack(errs[0].Error())
	return NewError(errs...)
}

func (logger SlackLogger) postSlack(errStr string) {
	program := viper.GetString("program")

	attachment := slacker.Attachment{
		Pretext: "@channel",
		Text:    "PRODUCTION ERROR",
		Color:   "#FF0055",
		Fields: []slacker.AttachmentField{
			{
				Title: "The following error was detected",
				Value: errStr,
			},
		},
	}

	_, _, err := logger.client.PostMessage(
		viper.GetString(program+".slack.errorChannel"),
		slacker.MsgOptionAttachments(attachment),
	)

	Handles(err).With(NOOP)

}
