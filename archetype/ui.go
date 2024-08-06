package archetype

import (
	"github.com/kensonjohnson/roguelike-game-go/assets"
	"github.com/kensonjohnson/roguelike-game-go/component"
	"github.com/kensonjohnson/roguelike-game-go/config"
	"github.com/yohamta/donburi"
)

var UITag = donburi.NewTag("ui")

func CreateNewUI(world donburi.World) {
	entity := world.Create(
		UITag,
		component.UI,
	)
	entry := world.Entry(entity)

	ui := component.UIData{
		MessageBox: component.UserMessageBoxData{},
		PlayerHUD:  component.PlayerHUDData{},
	}

	// Configure message box
	messageBoxTopPixel := (config.ScreenHeight - config.UIHeight) * config.TileHeight
	ui.MessageBox.Position = component.PositionData{
		X: 0,
		Y: messageBoxTopPixel,
	}
	ui.MessageBox.FontX = config.FontSize
	ui.MessageBox.FontY = messageBoxTopPixel + 10
	ui.MessageBox.Sprite = assets.UIPanel

	// Configure player HUD
	playerHUDTopPixel := (config.ScreenHeight * config.TileHeight) - 220
	playerEntry := PlayerTag.MustFirst(world)
	ui.PlayerHUD.Position = component.PositionData{
		X: config.ScreenWidth * config.TileWidth / 2,
		Y: playerHUDTopPixel,
	}
	ui.PlayerHUD.FontX = ui.PlayerHUD.Position.X + config.FontSize
	ui.PlayerHUD.FontY = messageBoxTopPixel + 12
	ui.PlayerHUD.Health = component.Health.Get(playerEntry)
	ui.PlayerHUD.Armor = component.Armor.Get(playerEntry)
	ui.PlayerHUD.Weapon = component.Weapon.Get(playerEntry)
	ui.PlayerHUD.Sprite = assets.UIPanelWithMinimap

	component.UI.SetValue(entry, ui)
}
