package bot

import (
	"context"
	"github.com/0fau/logs/pkg/database"
	"github.com/0fau/logs/pkg/database/sql"
	"github.com/bwmarrin/discordgo"
	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
	"os/signal"
	"slices"
	"syscall"
)

type DiscordConfig struct {
	Token string

	GuildID   string
	MessageID string
	RoleID    string
}

type Bot struct {
	DiscordConfig DiscordConfig

	DatabaseURL string

	db *database.DB
}

func (b *Bot) Run(ctx context.Context) error {
	sesh, err := discordgo.New("Bot " + b.DiscordConfig.Token)
	if err != nil {
		return errors.Wrap(err, "new discord session")
	}

	sesh.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})

	b.db, err = database.Connect(ctx, b.DatabaseURL)
	if err != nil {
		return errors.Wrap(err, "connecting to database")
	}

	sesh.AddHandler(b.handleWhitelist)

	if err := sesh.Open(); err != nil {
		return errors.Wrap(err, "opening discord session")
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	if err := sesh.Close(); err != nil {
		return errors.Wrap(err, "closing discord session")
	}

	return nil
}

func (b *Bot) handleWhitelist(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.GuildID != b.DiscordConfig.GuildID || r.MessageID != b.DiscordConfig.MessageID {
		return
	}

	if slices.Contains(r.Member.Roles, b.DiscordConfig.RoleID) {
		return
	}

	if err := b.whitelist(context.Background(), r.UserID); err != nil {
		log.Println(errors.Wrap(err, "whitelisting "+r.UserID))
		return
	}

	if err := s.GuildMemberRoleAdd(r.GuildID, r.UserID, b.DiscordConfig.RoleID); err != nil {
		log.Println(errors.Wrap(err, "GuildMemberRoleAdd "+r.UserID))
		return
	}

	log.Println("Whitelisted " + r.Member.User.ID + " (" + r.Member.User.Username + ")")
}

func (b *Bot) whitelist(ctx context.Context, discordID string) error {
	if err := crdbpgx.ExecuteTx(ctx, b.db.Pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		qtx := b.db.Queries.WithTx(tx)

		row, err := qtx.GetRolesByDiscordID(ctx, discordID)
		if err == nil {
			if slices.Contains(row.Roles, "alpha") || slices.Contains(row.Roles, "trusted") {
				return nil
			}

			if err := qtx.UpdateRoles(ctx, sql.UpdateRolesParams{
				ID:    row.ID,
				Roles: append(row.Roles, "alpha"),
			}); err != nil {
				return errors.Wrap(err, "updating roles")
			}
		} else if errors.Is(err, pgx.ErrNoRows) {
			if err := qtx.Whitelist(ctx, sql.WhitelistParams{
				Discord: discordID,
				Role:    "alpha",
			}); err != nil {
				return errors.Wrap(err, "whitelist")
			}
		} else {
			return errors.Wrap(err, "getting roles")
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "executing transaction")
	}
	return nil
}
