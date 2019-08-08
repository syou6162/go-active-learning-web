package command

import (
	"github.com/syou6162/go-active-learning-web/lib/recommend"
	"github.com/syou6162/go-active-learning-web/lib/twitter"
	"github.com/syou6162/go-active-learning-web/lib/update_model"
	"github.com/syou6162/go-active-learning-web/lib/web"
	"github.com/syou6162/go-active-learning/lib/add"
	"github.com/urfave/cli"
)

var Commands = []cli.Command{
	add.CommandAdd,
	update_model.CommandUpdateExampleModel,
	update_model.CommandUpdateTweetModel,
	recommend.CommandRecommend,
	web.CommandServe,
	twitter.CommandSetReferringTweets,
}
