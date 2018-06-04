package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/Bullpeen/stox"
	log "github.com/Sirupsen/logrus"
	"github.com/jirwin/quadlek/plugins/archive"
	"github.com/jirwin/quadlek/plugins/comics"
	"github.com/jirwin/quadlek/plugins/echo"
	"github.com/jirwin/quadlek/plugins/eslogs"
	"github.com/jirwin/quadlek/plugins/karma"
	"github.com/jirwin/quadlek/plugins/nextep"
	"github.com/jirwin/quadlek/plugins/random"
	"github.com/jirwin/quadlek/plugins/spotify"
	"github.com/jirwin/quadlek/plugins/twitter"
	"github.com/jirwin/quadlek/quadlek"
	cointip "github.com/morgabra/cointip/quadlek"

	"fmt"

	"github.com/Bullpeen/infobot"
	"github.com/urfave/cli"
)

const Version = "0.0.1"

func run(c *cli.Context) error {
	var apiToken string
	if c.IsSet("api-key") {
		apiToken = c.String("api-key")
	} else {
		cli.ShowAppHelp(c)
		return cli.NewExitError("Missing --api-key arg.", 1)
	}

	var verificationToken string
	if c.IsSet("verification-token") {
		verificationToken = c.String("verification-token")
	} else {
		cli.ShowAppHelp(c)
		return cli.NewExitError("Missing --verification-token arg.", 1)
	}

	dbPath := c.String("db-path")

	bot, err := quadlek.NewBot(context.Background(), apiToken, verificationToken, dbPath)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("error creating bot")
		return nil
	}

	err = bot.RegisterPlugin(echo.Register())
	if err != nil {
		fmt.Printf("error registering echo plugin: %s\n", err.Error())
		return nil
	}

	err = bot.RegisterPlugin(karma.Register())
	if err != nil {
		fmt.Printf("error registering karma plugin: %s\n", err.Error())
		return nil
	}

	err = bot.RegisterPlugin(random.Register())
	if err != nil {
		fmt.Printf("error registering random plugin: %s\n", err.Error())
		return nil
	}

	err = bot.RegisterPlugin(spotify.Register())
	if err != nil {
		fmt.Printf("error registering spotify plugin: %s\n", err.Error())
		return nil
	}

	err = bot.RegisterPlugin(archive.Register())
	if err != nil {
		fmt.Printf("error registering archive plugin: %s\n", err.Error())
		return nil
	}

	if c.IsSet("tvdb-key") {
		tvdbKey := c.String("tvdb-key")

		err = bot.RegisterPlugin(nextep.Register(tvdbKey))
		if err != nil {
			fmt.Printf("error registering nextep plugin: %s\n", err.Error())
			return nil
		}
	}

	err = bot.RegisterPlugin(infobot.Register())
	if err != nil {
		fmt.Printf("error registering infobot plugin: %s\n", err.Error())
		return nil
	}

	err = bot.RegisterPlugin(twitter.Register(
		c.String("twitter-consumer-key"),
		c.String("twitter-consumer-secret"),
		c.String("twitter-access-token"),
		c.String("twitter-access-secret"),
		// These must be twitter user ids, not names. https://tweeterid.com/ for easy conversion between the two.
		map[string]string{
			"25073877":           "politics", // @realDonaldTrump
			"830896623689547776": "politics", // @PresVillain
			"976366106561490944": "artfolio", // @DrawnDavidsOff
			"921111554371682304": "artfolio", // @DrawnDavidson
			"1581511":            "wwdc",     //@macrumorslive
			//added for retweet testing
			"778682": "quadlek-chat", // @jirwin

		},
	))

	coinbasePlugin := cointip.Register(
		c.String("coinbase-key"),
		c.String("coinbase-secret"),
		c.String("coinbase-account"),
	)
	if coinbasePlugin != nil {
		err = bot.RegisterPlugin(coinbasePlugin)
		if err != nil {
			panic(err)
		}
	}

	err = bot.RegisterPlugin(comics.Register(
		c.String("imgur-client-id"),
		"/opt/quad-assets/Arial.ttf",
	))
	if err != nil {
		fmt.Printf("error registering comic plugin: %s\n", err.Error())
		return err
	}
	esPlugin, err := eslogs.Register(c.String("es-endpoint"), c.String("es-index"))
	if err != nil {
		fmt.Printf("Error creating eslogs plugin: %s\n", err.Error())
		return err
	}
	err = bot.RegisterPlugin(esPlugin)
	if err != nil {
		fmt.Printf("Error registering eslogs plugin: %s\n", err.Error())
		return err
	}

	stoxPlugin := stox.Register(c.String("av-key"))
	err = bot.RegisterPlugin(stoxPlugin)
	if err != nil {
		fmt.Printf("Error registering stox plugin: %s\n", err.Error())
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	bot.Start()
	<-signals
	bot.Stop()

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "quadlek"
	app.Version = Version
	app.Usage = "a slack bot"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "api-key",
			Usage:  "The slack api token for the bot",
			EnvVar: "QUADLEK_API_TOKEN",
		},
		cli.StringFlag{
			Name:   "verification-token",
			Usage:  "The slack webhook verification token.",
			EnvVar: "QUADLEK_VERIFICATION_TOKEN",
		},
		cli.StringFlag{
			Name:   "db-path",
			Usage:  "The path where the database is stored.",
			Value:  "quadlek.db",
			EnvVar: "QUADLEK_DB_PATH",
		},
		cli.StringFlag{
			Name:   "tvdb-key",
			Usage:  "The TVDB api key for the bot, used by the nextep command",
			EnvVar: "QUADLEK_TVDB_TOKEN",
		},
		cli.StringFlag{
			Name:   "twitter-consumer-key",
			Usage:  "The consumer key for the twitter api",
			EnvVar: "QUADLEK_TWITTER_CONSUMER_KEY",
		},
		cli.StringFlag{
			Name:   "twitter-consumer-secret",
			Usage:  "The consumer secret for the twitter api",
			EnvVar: "QUADLEK_TWITTER_CONSUMER_SECRET",
		},
		cli.StringFlag{
			Name:   "twitter-access-token",
			Usage:  "The access key for the twitter api",
			EnvVar: "QUADLEK_TWITTER_ACCESS_TOKEN",
		},
		cli.StringFlag{
			Name:   "twitter-access-secret",
			Usage:  "The access secret for the twitter api",
			EnvVar: "QUADLEK_TWITTER_ACCESS_SECRET",
		},
		cli.StringFlag{
			Name:   "coinbase-key",
			Usage:  "The access key for the coinbase api",
			EnvVar: "QUADLEK_COINBASE_KEY",
		},
		cli.StringFlag{
			Name:   "coinbase-secret",
			Usage:  "The access secret for the coinbase api",
			EnvVar: "QUADLEK_COINBASE_SECRET",
		},
		cli.StringFlag{
			Name:   "coinbase-account",
			Usage:  "The bank account for the coinbase api",
			EnvVar: "QUADLEK_COINBASE_BANK_ACCOUNT",
		},
		cli.StringFlag{
			Name:   "imgur-client-id",
			Usage:  "The imgur client id",
			EnvVar: "QUADLEK_IMGUR_CLIENT_ID",
		},
		cli.StringFlag{
			Name:   "es-endpoint",
			Usage:  "The elastic search endpoint",
			EnvVar: "QUADLEK_ES_ENDPOINT",
		},
		cli.StringFlag{
			Name:   "es-index",
			Usage:  "The elastic search index",
			EnvVar: "QUADLEK_ES_INDEX",
		},
		cli.StringFlag{
			Name:   "av-key",
			Usage:  "Alpha Vantage Api Key",
			EnvVar: "QUADLEK_AV_KEY",
		},
	}

	app.Run(os.Args)
}
