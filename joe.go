package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-joe/file-memory"
	"github.com/go-joe/joe"
	"github.com/go-joe/slack-adapter"
	"github.com/pkg/errors"
)

type walbot struct {
	*joe.Bot
}

func main() {
	rand.Seed(time.Now().UnixNano())
	slackToken := os.Getenv("SLACK_TOKEN")
	var j *joe.Bot
	if slackToken != "" {
		j = joe.New("walbot",
			slack.Adapter(slackToken, slack.WithListenPassive()),
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

	b.lispBang(`list-thinges`, func(msg joe.Message) error {
		thinges := []string{}
		ok, err := b.Store.Get("thinges", &thinges)
		if err != nil {
			return err
		}
		if !ok {
			msg.Respond("no thinges defined yet... wat!")
			return nil
		}
		sort.Strings(thinges)
		msg.Respond(fmt.Sprintf("known thinges:\n%s", strings.Join(thinges, "\n")))
		return nil
	})
	b.lispBang(`magic8ball .*\?`, randomizer(magic8ball))
	b.lispBang(`make-thinge (.+)`, b.MakeThinge)
	b.lispBang(`overlord`, randomizer(overlord))
	b.lispBang(`ferengi`, randomizer(ferengi))
	b.lispBang(`rand ([0-9]+)(?: ([0-9]+))?`, roll)
	b.lispBang(`roll ([0-9]+)(?: ([0-9]+))?`, roll)

	err = b.Run()
	if err != nil {
		b.Logger.Fatal(err.Error())
	}
}

func roll(msg joe.Message) error {
	max, err := strconv.Atoi(msg.Matches[0])
	if err != nil {
		msg.Respond("could not parse stop value")
		return err
	}
	if max < 1 {
		msg.Respond("max must be greater than 0")
		return nil
	}

	count := 1
	if msg.Matches[1] != "" {
		count, err = strconv.Atoi(msg.Matches[1])
		if err != nil {
			msg.Respond("could not parse count value")
			return err
		}
	}
	for i := 0; i < count; i++ {
		msg.Respond("%d", rand.Intn(max)+1)
	}
	return nil
}

func randomizer(items []string) func(joe.Message) error {
	n := int64(len(items))
	return func(msg joe.Message) error {
		msg.Respond("%s", items[rand.Int63n(n)])
		return nil
	}
}

var errThingeExists = fmt.Errorf("thinge already exists")

func (b *walbot) MakeThinge(msg joe.Message) error {
	thinge := msg.Matches[0]

	resp, err := b.makeThinge(thinge)
	if errors.Is(err, errThingeExists) {
		msg.Respond(fmt.Sprintf("thinge %q already exists", thinge))
		return nil
	}

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
			return "", errThingeExists
		}
	}
	thinges = append(thinges, t)
	b.Store.Set("thinges", thinges)

	b.lispBang(t+`-add \s*(.+)\s*`, func(msg joe.Message) error {
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
	b.lispBang(t+`-del \s*(.+)\s*`, func(msg joe.Message) error {
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
	b.lispBang(t, func(msg joe.Message) error {
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
	})
	return "thinge created", nil
}

func (b *walbot) lispBang(pattern string, fn func(joe.Message) error) {
	b.Respond(`!`+pattern, fn)
	b.Respond(`\(`+pattern+`\)`, fn)
}
