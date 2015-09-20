/*

followtheleader - follow the presidential candidates with Go and Twitter

Copyright (c) 2015 RapidLoop

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package main

import (
	"log"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

const (
	// Create an app from https://apps.twitter.com/, create an access token
	// and copy-paste the values from there.
	CONSUMER_KEY        = "_REPLACE_THIS_"
	CONSUMER_SECRET     = "_REPLACE_THIS_"
	ACCESS_TOKEN        = "_REPLACE_THIS_"
	ACCESS_TOKEN_SECRET = "_REPLACE_THIS_"

	// The interface and port the web server should listen on.
	LISTEN_ADDRESS = "0.0.0.0:8080"
)

// Info from one twitter account.
type TwitterInfo struct {
	Name      string
	Image     string
	Followers int
	Tweets    int64
	Democrat  bool
}

// Stats about all candidates.
type Stats struct {
	sync.Mutex
	At      time.Time
	Twitter []TwitterInfo
}

func (s *Stats) Put(t []TwitterInfo) {
	s.Lock()
	defer s.Unlock()
	s.At = time.Now().In(TZ_ET)
	s.Twitter = t
}

func (s *Stats) Get() (time.Time, []TwitterInfo) {
	s.Lock()
	defer s.Unlock()
	return s.At, s.Twitter
}

// This is uh, how you sort in Go..
type ByFollowers []TwitterInfo

func (a ByFollowers) Len() int           { return len(a) }
func (a ByFollowers) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFollowers) Less(i, j int) bool { return a[i].Followers > a[j].Followers }

var stats = &Stats{}

var TZ_ET *time.Location

var democrats = []string{
	"LincolnChafee",  // Lincoln Chafee
	"HillaryClinton", // Hillary Clinton
	"MartinOMalley",  // Martin O'Malley
	"BernieSanders",  // Bernie Sanders
	"JimWebbUSA",     // Jim Webb
}

var republicans = []string{
	"JebBush",         // Jeb Bush
	"realbencarson",   // Dr. Ben Carson
	"ChrisChristie",   // Chris Christie
	"tedcruz",         // Ted Cruz
	"CarlyFiorina",    // Clary Fiorina
	"gov_gilmore",     // James Gilmore
	"LindseyGrahamSC", // Lindsey Graham
	"GovMikeHuckabee", // Mike Huckabee
	"BobbyJindal",     // Bobby Jindal
	"JohnKasich",      // John Kasich
	"GovernorPataki",  // George Pataki
	"RandPaul",        // Dr. Rand Paul
	"marcorubio",      // Marco Rubio
	"RickSantorum",    // Rick Santorum
	"realDonaldTrump", // Donald Trump
	"ScottWalker",     // Scott Walker
}

// Fetch stats from Twitter and store it for display by the web interface.
func fetchStats(api *anaconda.TwitterApi) error {
	handles := strings.Join(democrats, ",") + "," + strings.Join(republicans, ",")

	users, err := api.GetUsersLookup(handles, url.Values{})
	if err != nil {
		return err
	}

	t := make([]TwitterInfo, len(users))
	for i, u := range users {
		t[i].Name = u.Name
		t[i].Image = u.ProfileImageURL
		t[i].Followers = u.FollowersCount
		t[i].Tweets = u.StatusesCount
		for _, d := range democrats {
			if u.ScreenName == d {
				t[i].Democrat = true
				break
			}
		}
	}
	sort.Sort(ByFollowers(t))
	stats.Put(t)

	return nil
}

func main() {
	log.SetFlags(0)

	// Load ET time zone
	TZ_ET, _ = time.LoadLocation("America/New_York")

	// Init the Twitter API
	anaconda.SetConsumerKey(CONSUMER_KEY)
	anaconda.SetConsumerSecret(CONSUMER_SECRET)
	api := anaconda.NewTwitterApi(ACCESS_TOKEN, ACCESS_TOKEN_SECRET)

	// Fetch it once
	if err := fetchStats(api); err != nil {
		log.Fatal(err)
	}

	// Start the web interface
	go startWeb()

	// Keep updating the stats every minute
	for {
		time.Sleep(time.Minute)
		fetchStats(api)
	}
}
