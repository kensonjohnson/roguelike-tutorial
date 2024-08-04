package archetype

import (
	"github.com/kensonjohnson/roguelike-game-go/assets"
	"github.com/kensonjohnson/roguelike-game-go/component"
	"github.com/kensonjohnson/roguelike-game-go/component/gear"
	"github.com/norendren/go-fov/fov"
	"github.com/yohamta/donburi"
)

var PlayerTag = donburi.NewTag("player")

func CreateNewPlayer(world donburi.World) {
	player := world.Entry(world.Create(
		PlayerTag,
		component.Position,
		component.Sprite,
		component.Name,
		component.Fov,
		component.Armor,
		component.Weapon,
		component.Health,
		component.UserMessage,
	))

	// Grab level
	entry := MustFindDungeon(world)
	level := component.Dungeon.Get(entry).CurrentLevel

	// Set starting position
	startingX, startingY := level.Rooms[0].Center()
	position := component.PositionData{
		X: startingX,
		Y: startingY,
	}
	component.Position.SetValue(player, position)

	// Update player's field of view
	vision := component.FovData{VisibleTiles: fov.New()}
	vision.VisibleTiles.Compute(level, startingX, startingY, 8)
	component.Fov.SetValue(player, vision)

	// Set sprite
	sprite := component.SpriteData{
		Image: assets.Player,
	}
	component.Sprite.SetValue(player, sprite)

	// Set name
	name := component.NameData{Label: "Player"}
	component.Name.SetValue(player, name)

	// Set health
	health := component.HealthData{
		MaxHealth:     30,
		CurrentHealth: 30,
	}
	component.Health.SetValue(player, health)

	// Add gear
	component.Armor.SetValue(player, gear.Armor.PlateArmor)

	component.Weapon.SetValue(player, gear.Weapons.BattleAxe)

	// Set default messages
	component.UserMessage.SetValue(
		player,
		component.UserMessageData{
			AttackMessage:    "",
			DeadMessage:      "",
			GameStateMessage: "",
		},
	)
}
