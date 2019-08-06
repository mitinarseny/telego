package bot

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

func configureBot()

func Start(token string){
    b, err := tb.NewBot(tb.Settings{
        Token:    token,
        Poller:   nil,
        Reporter: nil,
    })
}
