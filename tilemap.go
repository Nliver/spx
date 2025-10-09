/*
 * Copyright (c) 2021 The XGo Authors (xgo.dev). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package spx

import (
	"fmt"
	"sort"

	spxfs "github.com/goplus/spx/v2/fs"
	"github.com/goplus/spx/v2/internal/engine"
	tm "github.com/goplus/spx/v2/internal/tilemap"

	"github.com/goplus/spbase/mathf"
)

type tilemapMgr struct {
	g     *Game
	datas *tm.TscnMapData
}

func (p *tilemapMgr) init(g *Game, fs spxfs.Dir, path string) {
	p.g = g
	if path == "" {
		return
	}
	var data tm.TscnMapData
	err := loadJson(&data, fs, path)
	if err != nil {
		panic(fmt.Sprintf("Failed to load tilemap JSON file %s: %v", path, err))
	}
	p.datas = &data
	tm.ConvertData(&data)
}
func (p *tilemapMgr) hasData() bool {
	return p.datas != nil
}

func (p *tilemapMgr) loadTilemaps(datas *tm.TscnMapData) {
	tm.LoadTilemaps(datas, p.g.setTileInfo__1, p.g.setTileMapLayerIndex, p.g.PlaceTiles__1)
}
func (p *tilemapMgr) loadDecorators(datas *tm.TscnMapData) {
	const headingOffset = -90.0
	for _, item := range datas.Decorators {
		position := item.Position.ToVec2()
		pivot := item.Pivot.ToVec2()
		assetPath := engine.ToAssetPath("tilemaps/" + item.Path)
		texSize := resMgr.GetImageSize(assetPath)
		colliderPivot := item.ColliderPivot.ToVec2().Add(pivot)
		pivot = pivot.Sub(texSize.Divf(2))
		p.g.createStaticSprite("tilemaps/"+item.Path, position, item.Ratation+headingOffset,
			item.Scale.ToVec2(), int64(item.ZIndex), pivot, item.ColliderType, colliderPivot, item.ColliderParams)
	}
}

func (p *tilemapMgr) loadSprites(datas *tm.TscnMapData) {

	sort.Slice(datas.Sprites, func(i, j int) bool {
		return datas.Sprites[i].Path < datas.Sprites[j].Path
	})

	for _, item := range datas.Sprites {
		sp, ok := p.g.sprs[item.Path]
		if ok {
			x, y := item.Position.X, item.Position.Y
			doClone(sp, nil, true, func(sprite *SpriteImpl) {
				sprite.SetXYpos(x, y)
				sprite.Show()
			})
		}
	}
}

func (p *tilemapMgr) parseTilemap() {
	if p.datas == nil {
		return
	}
	p.loadTilemaps(p.datas)
	p.loadDecorators(p.datas)
	//p.loadSprites(p.datas)

	// Update world size based on actual tilemap content
	p.calcWorldSize()
}

// calcWorldSize calculates and updates world size based on actual tile distribution in tilemap
func (p *tilemapMgr) calcWorldSize() {
	if p.datas == nil || len(p.datas.TileMap.Layers) == 0 {
		fmt.Println("[TILEMAP DEBUG] No tilemap data or layers, skipping world size update")
		return
	}

	tileSizeX := int(p.datas.TileMap.TileSize.Width)
	tileSizeY := int(p.datas.TileMap.TileSize.Height)

	var minX, maxX, minY, maxY int32 = 0, 0, 0, 0
	hasAnyTiles := false
	totalTiles := 0

	for _, layer := range p.datas.TileMap.Layers {
		tiles := p.parseTileDataForBounds(layer.TileData)
		totalTiles += len(tiles)
		for _, tile := range tiles {
			if !hasAnyTiles {
				minX, maxX = tile.X, tile.X
				minY, maxY = tile.Y, tile.Y
				hasAnyTiles = true
			} else {
				if tile.X < minX {
					minX = tile.X
				}
				if tile.X > maxX {
					maxX = tile.X
				}
				if tile.Y < minY {
					minY = tile.Y
				}
				if tile.Y > maxY {
					maxY = tile.Y
				}
			}
		}
	}

	if hasAnyTiles {
		minWorldX := int((minX) * int32(tileSizeX))
		maxWorldX := int((maxX + 1) * int32(tileSizeX)) // +1 to include the full size of the last tile
		minWorldY := int((minY - 1) * int32(tileSizeY)) // -1 to include the full size of the last tile
		maxWorldY := int((maxY) * int32(tileSizeY))

		worldWidth := maxWorldX - minWorldX
		worldHeight := maxWorldY - minWorldY

		p.g.minWorldX_ = minWorldX
		p.g.minWorldY_ = minWorldY
		p.g.worldWidth_ = worldWidth
		p.g.worldHeight_ = worldHeight

	} else {
		fmt.Println("[TILEMAP DEBUG] No tiles found in any layer")
	}
}

// parseTileDataForBounds parses tile data for boundary calculation (copied logic from internal/tilemap package)
func (p *tilemapMgr) parseTileDataForBounds(tileData []int32) []mathf.Vec2i {
	tileCount := len(tileData) / 5
	tiles := make([]mathf.Vec2i, 0, tileCount)

	for i := 0; i < len(tileData); i += 5 {
		if i+4 >= len(tileData) {
			break
		}

		tileX := tileData[i+1]
		tileY := tileData[i+2]

		tile := mathf.Vec2i{
			X: tileX,
			Y: tileY,
		}

		tiles = append(tiles, tile)
	}

	return tiles
}
