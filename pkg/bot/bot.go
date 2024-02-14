package bot

import (
	"context"
	"github.com/0fau/logs/pkg/database"
	"github.com/bwmarrin/discordgo"
	"github.com/cockroachdb/errors"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type DiscordConfig struct {
	Token string

	GuildID   string
	MessageID string
	RoleID    string

	LogChannelID string
}

type Config struct {
	DiscordConfig

	DatabaseURL string
}

type Bot struct {
	config *Config
	db     *database.DB
	sesh   *discordgo.Session

	commands map[string]*Command
}

func NewBot(config *Config) *Bot {
	return &Bot{config: config}
}

type InteractionHandler func(*discordgo.Session, *discordgo.InteractionCreate)

type Command struct {
	Command *discordgo.ApplicationCommand
	Handler InteractionHandler
}

func (b *Bot) Run(ctx context.Context) error {
	var err error
	b.sesh, err = discordgo.New("Bot " + b.config.Token)
	if err != nil {
		return errors.Wrap(err, "new discord session")
	}

	b.sesh.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})

	b.db, err = database.Connect(ctx, b.config.DatabaseURL, "logs_bot", false)
	if err != nil {
		return errors.Wrap(err, "connecting to database")
	}

	b.sesh.AddHandler(b.handleWhitelist)

	if err := b.sesh.Open(); err != nil {
		return errors.Wrap(err, "opening discord session")
	}
	defer b.sesh.Close()

	if err := b.RegisterCommands(
		b.friendCommand(),
		b.friendUserCommand(),
		b.unfriendCommand(),
		b.friendsCommand(),
	); err != nil {
		return errors.Wrap(err, "registering commands")
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	for _, cmd := range b.commands {
		if err := b.sesh.ApplicationCommandDelete(
			b.sesh.State.User.ID, b.config.GuildID, cmd.Command.ID,
		); err != nil {
			return errors.Wrapf(err, "deleting command %s", cmd.Command.Name)
		}
	}

	return nil
}

func (b *Bot) RegisterCommands(cmds ...*Command) error {
	b.commands = make(map[string]*Command)
	for _, cmd := range cmds {
		var err error
		cmd.Command, err = b.sesh.ApplicationCommandCreate(
			b.sesh.State.User.ID, b.config.GuildID, cmd.Command,
		)
		if err != nil {
			return errors.Wrap(err, "registering command")
		}
		b.commands[cmd.Command.Name] = cmd
	}

	b.sesh.AddHandler(func(s *discordgo.Session, ic *discordgo.InteractionCreate) {
		if cmd, ok := b.commands[ic.ApplicationCommandData().Name]; ok {
			cmd.Handler(s, ic)
		}
	})

	return nil
}
