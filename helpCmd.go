package main

var helpDocs = map[string]string{
	"login":    "<Username>: This sets the user to <username>.",
	"register": "<Username>: This registers the user <username> and logs that user in.",
	"reset":    "This deletes all users.",
	"users":    "This lists all users.",

	"agg":     "UNFINISHED",
	"addfeed": "<Name> <URL>: This adds the feed.",
	"feeds":   "List all feeds for all users.",
	"follow":  "<URL>: This follows a feed. Feed must be added first.",

	"following": "Lists all the feeds the user is following.",

	"help": "This displays this command.",
}
