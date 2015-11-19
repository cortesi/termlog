
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/cortesi/termlog)

# termlog: Logging for interactive terminals

![screenshot](_demo/screenshot.png "termlog in action")

# Basic usage

    l := termlog.NewLog()
    l.Say("Log")
    l.Notice("Notice!")
    l.Warn("Warn!")
    l.Shout("Error!")

Each log entry gets a timestamp.


# Groups

Entries can be grouped together under one timestamp, with subsequent lines
indented:

    g = l.Group()
    g.Say("This line gets a timestamp")
    g.Say("This line will be indented with no timestamp")
    g.Done()

Groups must be marked as .Done() before output is produced - a good use for
defer. Termlog ensures that all grouped entries appear together in output.
