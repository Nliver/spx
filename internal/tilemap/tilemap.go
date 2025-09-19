package tilemap

import (
	"sort"
)

type vec2i struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

// WorldPoint represents a 2D coordinate in world space (pixels)
type vec2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// tileSize represents the dimensions of a tile
type tileSize struct {
	Width  int32 `json:"width"`
	Height int32 `json:"height"`
}

// physicsData represents physics properties of a tile
type physicsData struct {
	CollisionPoints []vec2 `json:"collision_points,omitempty"`
	// other properties
}

// tileInfo represents information about a single tile in the tileset
type tileInfo struct {
	AtlasCoords vec2i       `json:"atlas_coords"`
	Physics     physicsData `json:"physics,omitempty"`
}

// tileSource represents a tileset source
type tileSource struct {
	ID          int32      `json:"id"`
	TexturePath string     `json:"texture_path"`
	Tiles       []tileInfo `json:"tiles"`
}

// tileSet represents the complete tileset information
type tileSet struct {
	Sources []tileSource `json:"sources"`
	// Other properties
}

// tileInstance represents a placed tile in the map
type tileInstance struct {
	TileCoords  vec2i `json:"tile_coords"`
	WorldCoords vec2  `json:"world_coords"`
	SourceID    int32 `json:"source_id"`
	AtlasCoords vec2i `json:"atlas_coords"`
}

// tilemapLayer represents a tilemap layer with compact tile data format
type tilemapLayer struct {
	ID       int32   `json:"id"`
	Name     string  `json:"name"`
	TileData []int32 `json:"tile_data"`
}

// tileMapData represents the complete tilemap data
type tileMapData struct {
	Format   int32          `json:"format"`
	TileSize tileSize       `json:"tile_size"`
	TileSet  tileSet        `json:"tileset"`
	Layers   []tilemapLayer `json:"layers"`
}

// decoratorNode represents a Sprite2D node in the scene
type decoratorNode struct {
	Name        string `json:"name"`
	Parent      string `json:"parent"`
	Position    vec2   `json:"position"`
	TexturePath string `json:"texture_path"`
	ZIndex      int32  `json:"z_index,omitempty"`
}

// spriteNode represents an instantiated prefab node in the scene
type spriteNode struct {
	Name       string                 `json:"name"`
	Parent     string                 `json:"parent"`
	Position   vec2                   `json:"position"`
	PrefabPath string                 `json:"prefab_path"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// TscnMapData represents the root structure for JSON output
type TscnMapData struct {
	TileMap    tileMapData     `json:"tilemap"`
	Decorators []decoratorNode `json:"decorators"`
	Sprites    []spriteNode    `json:"prefabs"`
}

// Runtime utilities for parsing tile data
func LoadTilemaps(datas *TscnMapData, funcSetTile func(texturePath string, isCollision bool), funcSetLayer func(layerIndex int64),
	funcPlaceTiles func(positions []float64, texturePath string, layerIndex int64)) {
	paths := make(map[int32]string)
	for _, item := range datas.TileMap.TileSet.Sources {
		paths[item.ID] = item.TexturePath
		hasCollision := false
		for _, tile := range item.Tiles {
			hasCollision = hasCollision || tile.Physics.CollisionPoints != nil
		}
		funcSetTile(item.TexturePath, hasCollision)
	}
	for idx, layer := range datas.TileMap.Layers {
		layerId := int64(idx)
		funcSetLayer(layerId)
		tileData := layer.TileData
		tiles := parseTileData(tileData, datas.TileMap.TileSize)
		sort.Slice(tiles, func(i, j int) bool {
			return tiles[i].SourceID < tiles[j].SourceID
		})
		lastId := int32(-1)
		path := ""
		positions := make([]float64, 0, len(tiles)*2)
		for _, tile := range tiles {
			if lastId != tile.SourceID {
				if len(positions) > 0 {
					funcPlaceTiles(positions, path, layerId)
				}
				positions = positions[:0]
				lastId = tile.SourceID
				path = paths[tile.SourceID]
			}
			y, x := tile.WorldCoords.Y, tile.WorldCoords.X
			positions = append(positions, x, -y)
		}
		if len(positions) > 0 {
			funcPlaceTiles(positions, path, layerId)
		}
	}
}

// TileMapParser provides utilities for parsing compact tile data
// ParseTileData converts compact tile data array to tile instances
// The tile_data format is: [x, source_id, atlas_coords_encoded, x2, source_id2, atlas_coords_encoded2, ...]
// Where atlas_coords_encoded combines atlas X and Y coordinates
func parseTileData(tileData []int32, tileSize tileSize) []tileInstance {
	var tiles []tileInstance

	for i := 0; i < len(tileData); i += 3 {
		if i+2 >= len(tileData) {
			break
		}

		tilePos := tileData[i]
		sourceID := tileData[i+1]
		atlasEncoded := tileData[i+2]

		// Decode tile position (Godot uses a specific encoding)
		tileX := tilePos & 0xFFFF
		if tileX >= 0x8000 {
			tileX -= 0x10000 // Handle negative coordinates
		}
		tileY := (tilePos >> 16) & 0xFFFF
		if tileY >= 0x8000 {
			tileY -= 0x10000 // Handle negative coordinates
		}

		// Decode atlas coordinates (usually just X and Y)
		atlasX := atlasEncoded & 0xFFFF
		atlasY := (atlasEncoded >> 16) & 0xFFFF

		tile := tileInstance{
			TileCoords: vec2i{X: tileX, Y: tileY},
			WorldCoords: vec2{
				X: float64(tileX * tileSize.Width),
				Y: float64(tileY * tileSize.Height),
			},
			SourceID:    sourceID,
			AtlasCoords: vec2i{X: atlasX, Y: atlasY},
		}

		tiles = append(tiles, tile)
	}

	return tiles
}
