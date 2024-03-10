package bot

import (
	"context"
	"github.com/0fau/logs/pkg/database/sql"
	"github.com/bwmarrin/discordgo"
	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
	"log"
	"strings"
)

func (b *Bot) friendUserCommand() *Command {
	return &Command{
		Command: &discordgo.ApplicationCommand{
			Name: "Friend",
			Type: discordgo.UserApplicationCommand,
		},
		Handler: func(session *discordgo.Session, create *discordgo.InteractionCreate) {
			b.friend(session, create)
		},
	}
}

func (b *Bot) friendCommand() *Command {
	return &Command{
		Command: &discordgo.ApplicationCommand{
			Name:        "friend",
			Description: "Add a friend",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionMentionable,
					Required:    true,
					Name:        "username",
					Description: "The username of the user to add as a friend",
				},
			},
		},
		Handler: func(session *discordgo.Session, create *discordgo.InteractionCreate) {
			b.friend(session, create)
		},
	}
}

func (b *Bot) unfriendCommand() *Command {
	return &Command{
		Command: &discordgo.ApplicationCommand{
			Name:        "unfriend",
			Description: "Remove a friend",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
					Name:        "username",
					Description: "The username of the user to remove as a friend",
				},
			},
		},
		Handler: func(session *discordgo.Session, ic *discordgo.InteractionCreate) {
			user1 := ic.Interaction.Member.User.ID
			user2 := ic.ApplicationCommandData().Options[0].Value.(string)

			if err := crdbpgx.ExecuteTx(context.Background(), b.db.Pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
				qtx := b.db.Queries.WithTx(tx)
				ctx := context.Background()

				u, err := qtx.GetUserByDiscordID(ctx, user1)
				if err != nil {
					if errors.Is(err, pgx.ErrNoRows) {
						respond(session, ic, "You are not signed up on the website.")
						return nil
					}
					return errors.Wrap(err, "Failed to check if user is signed up")
				}

				u2, err := qtx.GetUser(ctx, user2)
				if err != nil {
					if errors.Is(err, pgx.ErrNoRows) {
						respond(session, ic, "`@"+user2+"` is not signed up on the website.")
						return nil
					}
					return errors.Wrap(err, "Failed to check if user is signed up")
				}

				if u.ID == u2.ID {
					respond(session, ic, "You can't remove yourself as a friend.")
					return nil
				}

				yes, err := qtx.HasFriendRequest(ctx, sql.HasFriendRequestParams{
					User1: user1, User2: u2.DiscordID,
				})
				if err != nil {
					return errors.Wrap(err, "Failed to check if users are friends")
				}
				if yes {
					if err := qtx.DeleteFriendRequest(ctx, sql.DeleteFriendRequestParams{
						User1: user1, User2: u2.DiscordID,
					}); err != nil {
						return errors.Wrap(err, "Failed to delete friend request")
					}

					respond(session, ic, "You removed the friend request to `@"+user2+"`.")
					return nil
				}

				yes, err = qtx.HasFriendRequest(ctx, sql.HasFriendRequestParams{
					User1: u2.DiscordID, User2: user1,
				})
				if err != nil {
					return errors.Wrap(err, "Failed to check if users are friends")
				}
				if yes {
					if err := qtx.DeleteFriendRequest(ctx, sql.DeleteFriendRequestParams{
						User1: u2.DiscordID, User2: user1,
					}); err != nil {
						return errors.Wrap(err, "Failed to delete friend request")
					}

					respond(session, ic, "You declined the friend request from `@"+user2+"`.")
					return nil
				}

				yes, err = qtx.AreFriends(ctx, sql.AreFriendsParams{
					User1: u.ID, User2: u2.ID,
				})
				if err != nil {
					return errors.Wrap(err, "Failed to check if users are friends")
				}
				if !yes {
					respond(session, ic, "You are not friends with @"+user2+".")
					return nil
				}

				if err := qtx.DeleteFriend(ctx, sql.DeleteFriendParams{
					User1: u.ID, User2: u2.ID,
				}); err != nil {
					return errors.Wrap(err, "Failed to delete friend")
				}

				respond(session, ic, "Removed `@"+user2+"` as a friend.")
				return nil
			}); err != nil {
				respond(session, ic, "Something went boom.")
				log.Println(errors.Wrap(err, "Transaction failed"))
				return
			}
		},
	}
}

func (b *Bot) friend(
	session *discordgo.Session,
	ic *discordgo.InteractionCreate,
) {
	user1 := ic.Interaction.Member.User.ID
	var user2 string
	if ic.ApplicationCommandData().TargetID != "" {
		user2 = ic.ApplicationCommandData().TargetID
	} else {
		user2 = ic.ApplicationCommandData().Options[0].Value.(string)
	}
	username2 := ic.ApplicationCommandData().Resolved.Users[user2].Username

	if user1 == user2 {
		respond(session, ic, "You can't add yourself as a friend.")
		return
	}

	if err := crdbpgx.ExecuteTx(context.Background(), b.db.Pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		qtx := b.db.Queries.WithTx(tx)

		ctx := context.Background()

		u, err := qtx.GetUserByDiscordID(ctx, user1)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				respond(session, ic, "You are not signed up on the website.")
				return nil
			}
			return errors.Wrap(err, "Failed to check if user is signed up")
		}

		u2, err := qtx.GetUserByDiscordID(ctx, user2)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				respond(session, ic, "`@"+username2+"` is not signed up on the website.")
				return nil
			}
			return errors.Wrap(err, "Failed to check if user is signed up")
		}

		yes, err := qtx.AreFriends(ctx, sql.AreFriendsParams{
			User1: u.ID, User2: u2.ID,
		})
		if err != nil {
			return errors.Wrap(err, "Failed to check if users are friends")
		}
		if yes {
			respond(session, ic, "You are already friends with `@"+username2+"`.")
			return nil
		}

		yes, err = qtx.HasFriendRequest(ctx, sql.HasFriendRequestParams{
			User1: user1, User2: user2,
		})
		if err != nil {
			return errors.Wrap(err, "Failed to check if users have a friend request")
		}
		if yes {
			respond(session, ic, "You already sent a friend request to `@"+username2+"`.")
			return nil
		}

		yes, err = qtx.HasFriendRequest(ctx, sql.HasFriendRequestParams{
			User1: user2, User2: user1,
		})
		if err != nil {
			return errors.Wrap(err, "Failed to check if users have a friend request")
		}
		if yes {
			if err := qtx.DeleteFriendRequest(ctx, sql.DeleteFriendRequestParams{
				User1: user2, User2: user1,
			}); err != nil {
				return errors.Wrap(err, "Failed to delete friend request")
			}

			if err := qtx.CreateFriend(ctx, sql.CreateFriendParams{
				User1: u.ID, User2: u2.ID,
			}); err != nil {
				return errors.Wrap(err, "Failed to create friend")
			}

			respond(session, ic, "Added `@"+username2+"` as a friend.")
		} else {
			if err := qtx.SendFriendRequest(ctx, sql.SendFriendRequestParams{
				User1: user1, User2: user2,
			}); err != nil {
				return errors.Wrap(err, "Failed to send friend request")
			}

			respond(session, ic, "Sent friend request to `@"+username2+"`.")
		}

		return nil
	}); err != nil {
		respond(session, ic, "Something went boom.")
		log.Println(errors.Wrap(err, "Transaction failed"))
		return
	}
}

func (b *Bot) friendsCommand() *Command {
	return &Command{
		Command: &discordgo.ApplicationCommand{
			Name:        "friends",
			Description: "List your friends",
		}, Handler: func(session *discordgo.Session, ic *discordgo.InteractionCreate) {
			ctx := context.Background()
			user1 := ic.Interaction.Member.User.ID
			u, err := b.db.Queries.GetUserByDiscordID(ctx, user1)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					respond(session, ic, "You are not signed up on the website.")
					return
				}
				log.Println(errors.Wrap(err, "Failed to check if user is signed up"))
				return
			}

			friends, err := b.db.Queries.ListFriends(ctx, u.ID)
			if err != nil && !errors.Is(err, pgx.ErrNoRows) {
				log.Println(errors.Wrap(err, "Failed to get friends"))
				return
			}

			received, err := b.db.Queries.ListReceivedFriendRequests(ctx, user1)
			if err != nil && !errors.Is(err, pgx.ErrNoRows) {
				log.Println(errors.Wrap(err, "Failed to get received friend requests"))
				return
			}

			sent, err := b.db.Queries.ListSentFriendRequests(ctx, user1)
			if err != nil && !errors.Is(err, pgx.ErrNoRows) {
				log.Println(errors.Wrap(err, "Failed to get sent friend requests"))
				return
			}

			var buf strings.Builder
			buf.WriteString("```\nFriends:\n--------")
			for _, friend := range friends {
				if friend.DiscordTag == u.DiscordTag {
					continue
				}

				buf.WriteString("\n@")
				buf.WriteString(friend.DiscordTag)

				if friend.Username != nil {
					buf.WriteString(" (" + *friend.Username + ")")
				}
			}
			if len(friends) == 0 {
				buf.WriteString("\nNone")
			}

			if len(received) > 0 {
				buf.WriteString("\n\nReceived:\n---------")
				for _, user := range received {
					buf.WriteString("\n@")
					buf.WriteString(user.DiscordTag)

					if user.Username != nil {
						buf.WriteString(" (" + *user.Username + ")")
					}
				}
			}

			if len(sent) > 0 {
				buf.WriteString("\n\nSent:\n-----")
				for _, user := range sent {
					buf.WriteString("\n@")
					buf.WriteString(user.DiscordTag)

					if user.Username != nil {
						buf.WriteString(" (" + *user.Username + ")")
					}
				}
			}
			buf.WriteString("```")

			respond(session, ic, buf.String())
		},
	}
}

func respond(
	session *discordgo.Session,
	ic *discordgo.InteractionCreate,
	content string,
) {
	if err := session.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}); err != nil {
		log.Println(errors.Wrap(err, "Failed to respond to interaction"))
	}
}
