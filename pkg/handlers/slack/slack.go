/*
Copyright 2016 Skippbox, Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package slack

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/slack-go/slack"

	"github.com/marvasgit/kubestatewatch/config"
	"github.com/marvasgit/kubestatewatch/pkg/event"
	"github.com/marvasgit/kubestatewatch/pkg/message"
)

var slackColors = map[string]string{
	"Normal":  "good",
	"Warning": "warning",
	"Danger":  "danger",
}

var slackErrMsg = `
%s

You need to set both slack token and channel for slack notify,
using "--token/-t" and "--channel/-c", or using environment variables:

export KW_SLACK_TOKEN=slack_token
export KW_SLACK_CHANNEL=slack_channel

Command line flags will override environment variables

`

// Slack handler implements handler.Handler interface,
// Notify event to slack channel
type Slack struct {
	Token   string
	Channel string
	Title   string
}

// Init prepares slack configuration
func (s *Slack) Init(c *config.Config) error {
	token := c.Handler.Slack.Token
	channel := c.Handler.Slack.Channel

	if token == "" {
		token = os.Getenv("KW_SLACK_TOKEN")
	}

	if channel == "" {
		channel = os.Getenv("KW_SLACK_CHANNEL")
	}

	s.Token = token
	s.Channel = channel
	s.Title = message.GetTitle(c.Message.Title, "KW_SLACK_TITLE")

	return checkMissingSlackVars(s)
}

// Handle handles the notification.
func (s *Slack) Handle(e event.StatemonitorEvent) {
	api := slack.New(s.Token)
	attachment := prepareSlackAttachment(e, s)

	channelID, timestamp, err := api.PostMessage(s.Channel,
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true))
	if err != nil {
		logrus.Errorf("failed to send slack message %s\n", err)
		return
	}

	logrus.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}

func checkMissingSlackVars(s *Slack) error {
	if s.Token == "" || s.Channel == "" {
		return fmt.Errorf(slackErrMsg, "Missing slack token or channel")
	}

	return nil
}

func prepareSlackAttachment(e event.StatemonitorEvent, s *Slack) slack.Attachment {

	attachment := slack.Attachment{
		Fields: []slack.AttachmentField{
			{
				Title: s.Title,
				Value: e.Message(),
			},
		},
	}

	if color, ok := slackColors[e.Status]; ok {
		attachment.Color = color
	}

	attachment.MarkdownIn = []string{"fields"}

	return attachment
}
