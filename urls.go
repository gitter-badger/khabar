package main

import (
	"github.com/bulletind/khabar/handlers"
	"gopkg.in/simversity/gottp.v2"
)

func registerHandlers() {
	gottp.NewUrl("notifications", "^/notifications/?$",
		new(handlers.Notifications))

	gottp.NewUrl("stats", "^/notifications/stats/?$",
		new(handlers.Stats))

	gottp.NewUrl("notification", "^/notifications/(?P<_id>\\w+)/?$",
		new(handlers.Notification))

	gottp.NewUrl("channel", "^/channels/(?P<ident>\\w+)/?$",
		new(handlers.Gully))

	gottp.NewUrl("topic_channel",
		"^/topics/(?P<ident>\\w+)/channels/(?P<channel>\\w+)/?$",
		new(handlers.TopicChannel))

	gottp.NewUrl("topic", "^/topics/(?P<ident>\\w+)/?$",
		new(handlers.Topic))

	gottp.NewUrl("user_locale", "^/locales/(?P<user>\\w+)/?$",
		new(handlers.UserLocale))

	gottp.NewUrl("channels", "^/channels/?$", new(handlers.Gullys))

	gottp.NewUrl("topics", "^/topics/?$", new(handlers.Topics))

	gottp.NewUrl("bounceNotifications", "^/bounce/?$", new(handlers.Bounce))

	gottp.NewUrl("complaintNotifications", "^/complaint/?$", new(handlers.Complaint))
}
