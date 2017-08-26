package main

import (
  "github.com/bwmarrin/discordgo"
  "os"
  "os/signal"
  "fmt"
  "syscall"
)

var bot_token string
var discord *discordgo.Session

func init() {
  // Ensure the bot token is present:
  bot_token := os.Getenv("BOT_TOKEN")
  if len(bot_token) == 0 {
    fmt.Fprintf( os.Stderr, "[CRITICAL] Bot token not set:\n" )
    os.Exit(1)
  }
  // Create the DiscordGo session object
  var session_err error
  discord, session_err = discordgo.New("Bot "+bot_token)
  if session_err != nil {
    fmt.Fprintf(
      os.Stderr,
      "[CRITICAL] Unable to create Discord API session: %s\n",
      session_err.Error(),
    )
    os.Exit(2)
  }
}

func main() {

  discord.AddHandler(onMessage)

  discord_open_err := discord.Open()
  if discord_open_err != nil {
    fmt.Fprintf(
      os.Stderr,
      "[CRITICAL] Unable to open connection to Discord API gateway: %s",
      discord_open_err.Error(),
    )
    return
  }

  fmt.Println("Bot online.")
  bot_user, bot_user_err := discord.User("@me")
  if bot_user_err != nil {
    fmt.Printf(
      "[CRITICAL] Unable to get bot user's information - Please double check "+
      "your credentials.\n%s\n",
      bot_user_err.Error(),
    )
  } else {
    fmt.Printf(
      "Invite URL: "+
      "https://discordapp.com/oauth2/authorize"+
      "?&client_id=%s&scope=bot&permissions=0\n", bot_user.ID)
  }
  fmt.Println("Press CTRL-C to exit.")
  sc := make(chan os.Signal, 1)
  signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
  <-sc
  discord.Close()
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
  // Log all messages
  var msg_prefix string
  if isDirectMessage(s, m) == true {
    msg_prefix = "MESSAGE_PRIVATE"
  } else {
    msg_prefix = "MESSAGE_PUBLIC"
  }

  fmt.Fprintf(os.Stdout,
    "[%s] %s\n",
    msg_prefix,
    m.Content,
  )
  // Ignore all messages created by the bot itself
  // This isn't required in this specific example but it's a good practice.
  if m.Author.ID == s.State.User.ID {
    return
  }
  // If the message is "ping" reply with "Pong!"
  if m.Content == "ping" {
    s.ChannelMessageSend(m.ChannelID, "Pong!")
  }

  // If the message is "pong" reply with "Ping!"
  if m.Content == "pong" {
    s.ChannelMessageSend(m.ChannelID, "Ping!")
  }
}

func isDirectMessage(s *discordgo.Session, m *discordgo.MessageCreate) bool {
    channel, err := s.State.Channel(m.ChannelID)
    if err != nil {
      fmt.Printf(
        "[ERROR] An error occurred in isDirectMessage(): %s\n",
        err.Error(),
      )
      return false
    }
    return channel.IsPrivate
}