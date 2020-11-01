package main
import (
  "fmt"
  "os"
  "os/exec"
  "os/signal"
  "syscall"
  "github.com/bwmarrin/discordgo"
  "io/ioutil"
  "encoding/json"
  "strings"
)

// 構造体定義
type Config struct {
    DiscordToken  string  `json:"discordToken"`
}

func main() {
  c := new(Config)
  jsonString, err := ioutil.ReadFile("settings.json")
  if err != nil {
    fmt.Println("error:\n", err)
    return
  }
  err = json.Unmarshal(jsonString, &c)
  if err != nil {
    fmt.Println("error:\n", err)
    return
  }

  dg, err := discordgo.New("Bot " + c.DiscordToken)
  if err != nil {
    fmt.Println("error:start\n", err)
    return
  }
  dg.AddHandler(msgReceived)
  err = dg.Open()
  if err != nil {
    fmt.Println("error:wss\n", err)
    return
  }
  fmt.Println("BOT Running...")

  sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc
    dg.Close()
  }

func msgReceived(s *discordgo.Session, m *discordgo.MessageCreate) {
  if m.Author.Bot {
    return
  }

  nickname := m.Author.Username
  member,err := s.State.Member(m.GuildID, m.Author.ID)

  if err == nil && member.Nick != "" {
    nickname = member.Nick
  }
  fmt.Println(m.Content + " by " + nickname)

  if strings.Contains(m.Content, "!shgei") {
    messageString := string(m.Content)
    fmt.Println(messageString)
    cmdString := strings.Replace(messageString,"!shgei ","",-1)

    out,err := exec.Command("sh","-c", cmdString, "2>&1").Output()
    if err != nil {
      fmt.Println("error:start\n", err)
      s.ChannelMessageSend(m.ChannelID, "エラー")
      s.ChannelMessageSend(m.ChannelID, string(out))
      return
    }
    s.ChannelMessageSend(m.ChannelID, "出力")
    s.ChannelMessageSend(m.ChannelID, "```\n" + string(out) + "\n```")
  }
}
