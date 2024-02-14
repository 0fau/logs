package bot

import (
	"context"
	"github.com/0fau/logs/pkg/database/sql"
	"github.com/bwmarrin/discordgo"
	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
	"log"
	"slices"
	"time"
)

func (b *Bot) handleWhitelist(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.GuildID != b.config.GuildID || r.MessageID != b.config.MessageID {
		return
	}

	if slices.Contains(r.Member.Roles, b.config.RoleID) {
		return
	}

	if err := b.whitelist(context.Background(), r.UserID); err != nil {
		log.Println(errors.Wrap(err, "whitelisting "+r.UserID))
		return
	}

	if err := s.GuildMemberRoleAdd(r.GuildID, r.UserID, b.config.RoleID); err != nil {
		log.Println(errors.Wrap(err, "GuildMemberRoleAdd "+r.UserID))
		return
	}

	log.Println("Whitelisted " + r.Member.User.ID + " [" + r.Member.User.Username + "]")
	_, err := s.ChannelMessageSendEmbed(b.config.LogChannelID, &discordgo.MessageEmbed{
		Title:       "",
		Description: r.Member.User.Username + " signed up for the alpha <:playemon:1199514828453187634>",
		Footer:      &discordgo.MessageEmbedFooter{Text: "Nice!"},
		Timestamp:   time.Now().Format(time.RFC3339),
		Author:      &discordgo.MessageEmbedAuthor{Name: r.Member.User.Username + " (" + r.Member.User.ID + ")", IconURL: "https://cdn.discordapp.com/avatars/" + r.Member.User.ID + "/" + r.Member.User.Avatar + ".png"},
	})
	if err != nil {
		log.Println(errors.Wrap(err, "Failed to send alpha whitelist log message on discord"))
	}
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
