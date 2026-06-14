package main

var helpDocs = map[string]string{
	"login":    "<Username>: This sets the user to <Username>.",
	"register": "<Username>: This registers the user <Username> and logs that user in.",
	"reset":    "Prompts you if you want to delete all users.",
	"users":    "This lists all users.",

	"agg":     "<Duration>: This starts the aggregation process. Leave this running in the background to keep your system up to date on the latest posts.",
	"addfeed": "<Name> <URL>: This adds the feed, and makes it avalible to follow.",
	"feeds":   "This lists all feeds for all users.",
	"follow":  "<URL>: This follows a feed. Feed must be added first.",

	"following": "Lists all the feeds the user is following.",
	"unfollow":  "<URL>: This removes the feed the user is following.",
	"browse":    "<PostCount (Default=2)>: Shows posts in order of most to least recent.",

	"help": "This displays this command.",
}
