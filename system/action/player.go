package action

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kensonjohnson/roguelike-game-go/archetype"
	"github.com/kensonjohnson/roguelike-game-go/component"
	"github.com/kensonjohnson/roguelike-game-go/event"
	"github.com/kensonjohnson/roguelike-game-go/system/combat"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func TakePlayerAction(ecs *ecs.ECS) bool {
	turnTaken := false

	// Capture any keypresses we care about
	moveX := 0
	moveY := 0

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		moveX -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		moveX += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		moveY -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		moveY += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		turnTaken = true
	}

	if turnTaken {
		return true
	}

	if moveX == 0 && moveY == 0 {
		return false
	}

	// Grab current level
	levelEntry := archetype.LevelTag.MustFirst(ecs.World)
	level := component.Level.Get(levelEntry)

	// Grab player data
	playerEntry := archetype.PlayerTag.MustFirst(ecs.World)
	position := component.Position.Get(playerEntry)
	sprite := component.Sprite.Get(playerEntry)
	vision := component.Fov.Get(playerEntry)

	// TODO: Update so diagonal movement consumes two turns
	// Attempt to move
	tile := level.GetFromXY(position.X+moveX, position.Y+moveY)
	if !tile.Blocked {
		// We can move the player
		// Free up the tile we are currently on and block the one we're going to
		level.GetFromXY(position.X, position.Y).Blocked = false
		tile.Blocked = true
		// Update the player's position
		position.X += moveX
		position.Y += moveY
		sprite.OffestX = -moveX
		sprite.OffestY = -moveY
		sprite.Animating = true
		sprite.SetProgress(0)
		// Update the player's field of view
		vision.VisibleTiles.Compute(level, position.X, position.Y, 8)
		// Update any discoverable entities
		component.Discoverable.Each(ecs.World, func(entry *donburi.Entry) {
			discoverablePosition := component.Position.Get(entry)
			if vision.VisibleTiles.IsVisible(discoverablePosition.X, discoverablePosition.Y) {
				discoverable := component.Discoverable.Get(entry)
				discoverable.SeenByPlayer = true
			}
		})

		query := donburi.NewQuery(
			filter.Contains(
				archetype.ItemTag,
				component.Position,
				component.Name,
			))

		query.Each(ecs.World, func(entry *donburi.Entry) {
			itemPosition := component.Position.Get(entry)
			if position.X == itemPosition.X && position.Y == itemPosition.Y {
				// The character has moved on top of a pickup
				archetype.RemoveItemFromWorld(entry)
				itemName := component.Name.Get(entry)
				playerMessages := component.UserMessage.Get(playerEntry)
				playerMessages.WorldInteractionMessage = fmt.Sprintf("Picked up %s!", itemName.Value)
				// TODO: place in player's inventory
			}
		})

		if tile.TileType == component.STAIR_DOWN {
			// Move to the next level
			event.ProgressLevelEvent.Publish(ecs.World, event.ProgressLevel{})
		}

	} else if tile.TileType != component.WALL {
		// Not a wall, so it must be an enemy. Attack!
		// Find the monster in the direction we're pointing
		enemyPosition := component.PositionData{
			X: position.X + moveX,
			Y: position.Y + moveY,
		}
		var monsterEntry *donburi.Entry
		archetype.MonsterTag.Each(ecs.World, func(entry *donburi.Entry) {
			position := component.Position.Get(entry)
			if position.IsEqual(&enemyPosition) {
				monsterEntry = entry
			}
		})
		combat.AttackSystem(ecs.World, playerEntry, monsterEntry)
	}

	return true
}
