package main

var helpDocs = map[string]string{
	"login":    "<Username>: This sets the user to <Username>.",
	"register": "<Username>: This registers the user <Username> and logs that user in.",
	"reset":    "This deletes all users.",
	"users":    "This lists all users.",

	"agg":     "<Duration>: This starts the aggregation. (TODO: EXPLAIN)",
	"addfeed": "<Name> <URL>: This adds the feed.",
	"feeds":   "This lists all feeds for all users.",
	"follow":  "<URL>: This follows a feed. Feed must be added first.",

	"following": "Lists all the feeds the user is following.",
	"unfollow":  "<URL>: This removes the feed the user is following.",

	"help": "This displays this command.",
}
