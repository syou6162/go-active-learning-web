package command

import (
	"github.com/codegangsta/cli"
	"github.com/syou6162/go-active-learning-web/lib/recommend"
	"github.com/syou6162/go-active-learning-web/lib/twitter"
	"github.com/syou6162/go-active-learning-web/lib/web"
	"github.com/syou6162/go-active-learning/lib/add"
)

var Commands = []cli.Command{
	add.CommandAdd,
	recommend.CommandRecommend,
	web.CommandServe,
	twitter.CommandSetReferringTweets,
}
