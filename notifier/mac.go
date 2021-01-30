package notifier

import (
	"fmt"

	gosxnotifier "github.com/deckarep/gosx-notifier"
	"github.com/sirupsen/logrus"
)

type NotificationMessage struct {
	Title    string
	Subtitle string
	Message  string
	Sound    gosxnotifier.Sound
	Url      string

	Hash string
}

const (
	defaultAppIcon    = "../static/gopher.png"
	defaultContentImg = "../static/gopher.png"
)

func Notify(logger logrus.FieldLogger, msg *NotificationMessage) {
	logger = logger.WithField("method", "Notify")
	if msg == nil || msg.Message == "" {
		logger.Error("empty notification invoked!")
	}
	note := gosxnotifier.NewNotification(msg.Message)
	note.Title = msg.Title
	note.Subtitle = msg.Subtitle
	note.Sound = msg.Sound
	note.Link = msg.Url

	note.Group = fmt.Sprintf("com.miko.shop.%s", msg.Hash)

	note.AppIcon = defaultAppIcon
	note.ContentImage = defaultContentImg

	//Then, push the notification
	err := note.Push()

	//If necessary, check error
	if err != nil {
		logger.WithError(err).Error("cannot push notification")
	}
}
