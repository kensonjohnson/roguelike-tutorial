package component

import "github.com/yohamta/donburi"

type WalletData struct {
	Amount int
}

var Wallet = donburi.NewComponentType[WalletData]()

func (w *WalletData) AddAmount(amount int) {
	w.Amount += amount
}

func (w *WalletData) SubtractAmount(amount int) {
	w.Amount -= amount
}
