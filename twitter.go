package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/mrjones/oauth"
	"os"
	"os/exec"
	"regexp"
)

const delim = '\n'

/*
I don't want to share my twitter secrets.
export your own consumer key and secret to the environment variables
TWITTER_CONSUMER_KEY and TWITTER_CONSUMER_SECRET respectively.
*/

var consumerkey = os.Getenv("TWITTER_CONSUMER_KEY")
var consumersecret = os.Getenv("TWITTER_CONSUMER_SECRET")
var usertoken = ""
var usersecret = ""
var userid = ""
var screenname = ""
var loggedin = false

func main() {
	clear()
	flag.Parse()

	fmt.Printf("Luke's Awesome Twitter Client in Go: Tweetingo\n\n")

	// Thanks to ChimeraCoder for the anaconda library.
	anaconda.SetConsumerKey(consumerkey)
	anaconda.SetConsumerSecret(consumersecret)
	api := anaconda.NewTwitterApi(usertoken, usersecret)

	// Loop until runnning == false.
	var running = true

	for running {

		//Stuff you can do with out being authenticated.
		menu := map[string]string{
			"S": "earch",
			"C": "lear",
			"L": "ogin",
			"Q": "uit",
		}

		//Stuff that gets appended after you're authenticated.
		if loggedin == true {
			menu["M"] = "y Timeline"
			menu["P"] = "ost"
		}

		showMenu(menu)

		// Get a value from the user.
		var answer string
		fmt.Scanf("%s", &answer)

		// Take a look at the answer and switch accordingly.
		switch answer {
		case "S", "s":
			getHashtag(api)
		case "P", "p":
			postTweet(api)
		case "Q", "q":
			fmt.Printf("Goodbye.\n")
			running = false
		case "C", "c":
			clear()
		case "M", "m":
			getMyTweets(api)

		case "L", "l":
			atoken := login()
			usertoken = atoken.Token
			usersecret = atoken.Secret
			userid = atoken.AdditionalData["user_id"]
			screenname = atoken.AdditionalData["screen_name"]
			api = anaconda.NewTwitterApi(usertoken, usersecret)
			loggedin = true

		}
	}
}

// Just generating the menu string so it's all on one line.
func showMenu(menu map[string]string) {
	var menuString string
	for key, value := range menu {
		menuString = menuString + "[" + key + "]" + value + " "
	}
	if loggedin == true {
		fmt.Printf("Hello %s!, Your UserID is: %s\n", screenname, userid)
	}
	fmt.Printf("%s: ", menuString)
}

// Log in to twitter. Thanks to mrjones for the oauth library.
func login() *oauth.AccessToken {

	// Set the twitter oauth URLs.
	sp := oauth.ServiceProvider{
		RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
		AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
		AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
	}

	// Request a temporary token for the application.
	consumer := oauth.NewConsumer(consumerkey, consumersecret, sp)

	rtoken, url, err := consumer.GetRequestTokenAndUrl("oob")
	if err != nil {
		fmt.Printf("Something went wrong\n")
		fmt.Printf(err.Error())
	}

	// This is lame, but it's how oauth works. Prompt the user to visit this URL in a browser.
	fmt.Printf("\n\nPlease Authorize this App by visiting the OAuth URL in your browser: \n%s \n\n", url)
	fmt.Printf("Please enter the PIN code provided by twitter here: ")

	// Once authenticated, twitter will return a pin. Grab it from the user.
	r := bufio.NewReader(os.Stdin)
	a, err := r.ReadString(delim)

	reg := regexp.MustCompile("[^0-9]")
	pin := reg.ReplaceAllString(a, "")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Request an authorized token.
	atoken, err := consumer.AuthorizeToken(rtoken, pin)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("You're logged in!\n")
	fmt.Println()

	if true == false {
		fmt.Printf("This,%v", atoken)
	}

	// Send back the atoken so you can use it in later requests.
	return atoken
}

// Get the current users timeline
func getMyTweets(api *anaconda.TwitterApi) {
	searchResult, _ := api.GetHomeTimeline(nil)

	for _, tweet := range searchResult {

		fmt.Println("Time: " + tweet.CreatedAt)
		fmt.Println("Author: " + tweet.User.ScreenName)
		fmt.Println("\x1b[31;1mBody: \x1b[0m" + tweet.Text)
		fmt.Println()
	}
}

// Take a hashtag and retrieve the lastest matching tweets.
func getHashtag(api *anaconda.TwitterApi) {

	fmt.Printf("Please enter your search term: #")

	r := bufio.NewReader(os.Stdin)
	line, err := r.ReadString(delim)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	searchResult, _ := api.GetSearch("#"+line, nil)

	for _, tweet := range searchResult {
		fmt.Println("Time: " + tweet.CreatedAt)
		fmt.Println("\x1b[31;1mBody: \x1b[0m" + tweet.Text)
		fmt.Println()
	}
}

// Tweet a string from stdin.
func postTweet(api *anaconda.TwitterApi) {
	fmt.Printf("What do you want to say?\n")

	t := bufio.NewReader(os.Stdin)
	line2, err := t.ReadString(delim)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tweetResult, _ := api.PostTweet(line2, nil)
	fmt.Printf("You posted: %s\n", tweetResult.Text)
}

// Clear the screen
func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
