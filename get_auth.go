package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"code.google.com/p/goauth2/oauth"
	drive "code.google.com/p/google-api-go-client/drive/v2"
	oauth2 "code.google.com/p/google-api-go-client/oauth2/v2"
)

var (
	clientId     = flag.String("id", "", "Client ID")
	clientSecret = flag.String("secret", "", "Client Secret")
	apiURL       = flag.String("api", strings.Join([]string{drive.DriveScope, drive.DriveFileScope, drive.DriveAppdataScope, oauth2.UserinfoEmailScope, oauth2.UserinfoProfileScope, "https://spreadsheets.google.com/feeds"}, " "), "Scope")
	code         = flag.String("code", "", "Authorization Code")
	cachefile    = flag.String("cache", "auth.json", "Token cache file")
)

const usageMsg = `
To obtain a request token you must specify both -id and -secret.

To obtain Client ID and Secret, see the "OAuth 2 Credentials" section under
the "API Access" tab on this page: https://code.google.com/apis/console/

Once you have completed the OAuth flow, the credentials should be stored inside
the file specified by -cache and you may run without the -id and -secret flags.
`

func writeToFile(data []byte, filePath string) error {
    fo, err := os.Create(filePath)
    if err != nil { return err }
    defer func() {
	if err:= fo.Close(); err != nil { panic(err) }
    }()
    if _, err := fo.Write(data); err != nil { return err }
    return nil
}

func oauthreq(config *oauth.Config) *oauth.Transport{

	// Set up a Transport using the config.
	transport := &oauth.Transport{Config: config}

	// Try to pull the token from the cache; if this fails, we need to get one.
	token, err := config.TokenCache.Token()
	if err != nil {
		if *clientId == "" || *clientSecret == "" {
			flag.Usage()
			fmt.Fprint(os.Stderr, usageMsg)
			os.Exit(2)
		}
		if *code == "" {
			// Get an authorization code from the data provider.
			// ("Please ask the user if I can access this resource.")
			url := config.AuthCodeURL("statss")
			fmt.Println("Visit this URL to get a code, then run again with -code=YOUR_CODE\n")
			fmt.Println(url)
			os.Exit(2)
		}
		// Exchange the authorization code for an access token.
		// ("Here's the code you gave the user, now give me a token!")
		token, err = transport.Exchange(*code)
		if err != nil {
			log.Fatal("Exchange:", err)
		}
		// (The Exchange method will automatically cache the token.)
		fmt.Printf("Token is cached in %v\n", config.TokenCache)
	}

	// Make the actual request using the cached token to authenticate.
	// ("Here's the token, let me in!")
	transport.Token = token
	return transport
}


func About(d *drive.Service) (*drive.About, error) {
    result, err := d.About.Get().Do()
    if err != nil {
	fmt.Printf("An error occurred: %v\n", err)
	return nil, err
    }
    return result, nil
}


func main() {
	flag.Parse()

	// Set up a configuration.
	config := &oauth.Config{
		ClientId:     *clientId,
		ClientSecret: *clientSecret,
		Scope:        *apiURL,
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
		TokenCache:   oauth.CacheFile(*cachefile),
		AccessType:   "offline",
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
	}

	// Set up a Transport using the config.
	transport := oauthreq(config)
	// Make the request.

	d, err := drive.New(transport.Client())
	if err != nil {
	    log.Fatal("drive.New:", err)
	}
	about, err := About(d)
	if err != nil {
	    log.Fatal("drive.About:", err)
	}
	log.Print("Hello ", about.Name)
	fmt.Println("Done")
}
