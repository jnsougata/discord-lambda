package main


func main() {
	app := Client{
		DiscordToken: "token",
		DiscordPublicKey: "public key",
		DiscordApplicationID: "bot id",
	}
	app.Run(8080)
}