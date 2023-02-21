package nostr

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip04"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/spf13/viper"
)

type Entry struct {
	Key      string
	Password string
	nostrKey string
}

func AddEntry(entry *Entry) {
	result, _ := json.Marshal(&entry)

	sk := viper.GetString("KEY")
	pub, _ := nostr.GetPublicKey(sk)

	shared, _ := nip04.ComputeSharedSecret(pub, sk)
	msg, _ := nip04.Encrypt(string(result), shared)

	ev := nostr.Event{
		PubKey:    pub,
		CreatedAt: time.Now(),
		Kind:      8,
		Tags:      nostr.Tags{nostr.Tag{"p", pub}},
		Content:   msg,
	}

	// calling Sign sets the event ID field and the event Sig field
	ev.Sign(sk)

	publishEvent(&ev)
}

func ListEntries() []*Entry {
	evs := fetchEvents()
	entries := eventsToEntries(evs)

	result := make([]*Entry, 0, len(entries))

	for _, entry := range entries {
		result = append(result, entry)
	}

	return result
}

func GetEntry(key string) *Entry {
	evs := fetchEvents()
	entries := eventsToEntries(evs)

	if entries[key] == nil {
		return nil
	}

	return entries[key]
}

func DeleteEntry(key string) {
	evs := fetchEvents()
	entries := eventsToEntries(evs)

	sk := viper.GetString("KEY")
	pub, _ := nostr.GetPublicKey(sk)

	if entries[key] == nil {
		return
	}

	ev := nostr.Event{
		Kind:    5,
		PubKey:  pub,
		Tags:    nostr.Tags{nostr.Tag{"e", entries[key].nostrKey}},
		Content: "Deleting " + entries[key].nostrKey,
	}

	ev.Sign(sk)
	publishEvent(&ev)
}

func eventsToEntries(events map[string]*nostr.Event) map[string]*Entry {
	entries := make(map[string]*Entry)
	for _, ev := range events {
		shared, _ := nip04.ComputeSharedSecret(ev.PubKey, viper.GetString("KEY"))
		msg, _ := nip04.Decrypt(ev.Content, shared)

		var entry Entry
		json.Unmarshal([]byte(msg), &entry)
		entry.nostrKey = ev.ID

		if (entry.Key == "") || (entry.Password == "") {
			continue
		}

		entries[entry.Key] = &entry
	}

	return entries
}

func fetchEvents() map[string]*nostr.Event {
	events := make(map[string]*nostr.Event)

	for _, url := range strings.Split(viper.GetString("RELAYS"), ",") {
		relay, err := nostr.RelayConnect(context.Background(), strings.TrimSpace(url))
		if err != nil {
			continue
		}

		pub, _ := nostr.GetPublicKey(viper.GetString("KEY"))
		npub, _ := nip19.EncodePublicKey(pub)

		var filter nostr.Filter
		if _, v, err := nip19.Decode(npub); err == nil {
			pub := v.(string)
			filter = nostr.Filter{
				Kinds:   []int{8},
				Authors: []string{pub},
			}
		} else {
			fmt.Println(err)
		}

		ctx, _ := context.WithCancel(context.Background())
		evs := relay.QuerySync(ctx, filter)

		for _, ev := range evs {
			events[ev.ID] = ev
		}
	}

	return events
}

func publishEvent(ev *nostr.Event) {
	for _, url := range strings.Split(viper.GetString("RELAYS"), ",") {
		url := strings.TrimSpace(url)

		relay, e := nostr.RelayConnect(context.Background(), url)
		if e != nil {
			fmt.Println(e)
			continue
		}
		fmt.Println("published to ", url, relay.Publish(context.Background(), *ev))
	}
}
