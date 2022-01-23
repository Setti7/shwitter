package query

import (
	"github.com/Setti7/shwitter/internal/entity"
	"github.com/Setti7/shwitter/internal/signal"
)

func registerShweetCreationCallbacks() {
	// Insert shweet into current user timeline.
	insertIntoTimeLineCallback := func(name string, instance interface{}, args ...interface{}) {
		shweet := instance.(*entity.Shweet)
		InsertShweetIntoLine(shweet.UserID, shweet, entity.TimeLine)
	}
	signal.PostCreate.Connect(entity.Shweet{}, insertIntoTimeLineCallback)

	// Insert shweet into current user userline.
	insertIntoUserLineCallback := func(name string, instance interface{}, args ...interface{}) {
		shweet := instance.(*entity.Shweet)
		InsertShweetIntoLine(shweet.UserID, shweet, entity.UserLine)
	}
	signal.PostCreate.Connect(entity.Shweet{}, insertIntoUserLineCallback)

	// Insert shweet into followers timeline.
	insertIntoFollowersTimelinesCallback := func(name string, instance interface{}, args ...interface{}) {
		shweet := instance.(*entity.Shweet)
		BulkInsertShweetIntoFollowersTimelines(shweet.UserID, shweet)
	}
	signal.PostCreate.Connect(entity.Shweet{}, insertIntoFollowersTimelinesCallback)

}

func init() {
	registerShweetCreationCallbacks()
}
