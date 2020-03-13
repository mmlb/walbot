package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/go-joe/file-memory"
	"github.com/go-joe/joe"
	"github.com/go-joe/slack-adapter"
	"github.com/pkg/errors"
)

type walbot struct {
	*joe.Bot
}

type thinge []string

func main() {
	rand.Seed(time.Now().UnixNano())
	slackToken := os.Getenv("SLACK_TOKEN")
	var j *joe.Bot
	if slackToken != "" {
		j = joe.New("walbot",
			slack.Adapter(slackToken),
			file.Memory("walbot.json"))
	} else {
		j = joe.New("walbot",
			file.Memory("walbot.json"))
	}
	b := &walbot{Bot: j}

	thinges := []string{}
	ok, err := b.Store.Get("thinges", &thinges)
	if err != nil {
		panic(err)
	}
	err = b.Store.Set("thinges", nil)
	if err != nil {
		panic(err)
	}
	if ok {
		for _, thinge := range thinges {
			b.makeThinge(thinge)
		}
	}

	b.Respond(`!make-thinge (.+)`, b.MakeThinge)
	b.Respond(`\(make-thinge (.+)\)`, b.MakeThinge)

	err = b.Run()
	if err != nil {
		b.Logger.Fatal(err.Error())
	}
}

func (b *walbot) MakeThinge(msg joe.Message) error {
	resp, err := b.makeThinge(msg.Matches[0])
	msg.Respond(resp)
	return err
}

func (b *walbot) makeThinge(t string) (string, error) {
	thinges := []string{}
	_, err := b.Store.Get("thinges", &thinges)
	if err != nil {
		return "", errors.New("error getting thinges")
	}
	for _, v := range thinges {
		if v == t {
			return fmt.Sprintf("thinge %s is already defined", t), nil
		}
	}
	thinges = append(thinges, t)
	b.Store.Set("thinges", thinges)

	b.Respond(`\(`+t+`-add (.+)\)`, func(msg joe.Message) error {
		thinge := []string{}
		_, err := b.Store.Get("thinge."+t, &thinge)
		if err != nil {
			return errors.New("error getting thinges")
		}

		thinge = append(thinge, msg.Matches[0])
		err = b.Store.Set("thinge."+t, thinge)
		if err != nil {
			return errors.New("error saving thinge")
		}

		msg.Respond("%s added", t)
		return nil
	})
	b.Respond(`\(`+t+`-del (.+)\)`, func(msg joe.Message) error {
		thinge := []string{}
		_, err := b.Store.Get("thinge."+t, &thinge)
		if err != nil {
			return errors.New("error getting thinges")
		}

		found := false
		i := 0
		v := ""
		for i, v = range thinge {
			if v == msg.Matches[0] {
				found = true
				break
			}
		}
		if found {
			thinge[i] = thinge[len(thinge)-1]
			thinge[len(thinge)-1] = ""
			thinge = thinge[:len(thinge)-1]
			err := b.Store.Set("thinge."+t, thinge)
			if err != nil {
				return errors.New("error saving thinge")
			}
			msg.Respond("%s removed", t)
		} else {
			msg.Respond("%s not found", t)
		}

		return nil
	})
	get := func(msg joe.Message) error {
		thinge := []string{}
		_, err := b.Store.Get("thinge."+t, &thinge)
		if err != nil {
			return errors.New("error getting thinges")
		}

		n := int64(len(thinge))
		if n == 0 {
			msg.Respond("awww 💩, no %ss saved yet", t)
			return nil
		}
		msg.Respond("%s", thinge[rand.Int63n(n)])
		return nil
	}
	b.Respond(`\(`+t+`\)`, get)
	b.Respond(`!`+t, get)
	return "thinge created", nil
}
