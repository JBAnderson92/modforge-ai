package mods

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
)

// GameType represents supported game types
type GameType string

const (
	GameTypeMinecraft GameType = "minecraft"
	GameTypeSkyrim    GameType = "skyrim"
	GameTypeLua       GameType = "lua"
	GameTypeUnknown   GameType = "unknown"
)

// FileInfo contains information about an uploaded mod file
type FileInfo struct {
	Filename    string
	Size        int64
	GameType    GameType
	ContentType string
	Metadata    map[string]interface{}
}

// DetectGameType detects the game type from a file
func DetectGameType(header *multipart.FileHeader, content []byte) GameType {
	filename := strings.ToLower(header.Filename)
	ext := filepath.Ext(filename)

	// Detect by file extension first
	switch ext {
	case ".json":
		if isMinecraftJSON(content) {
			return GameTypeMinecraft
		}
	case ".esp", ".esm":
		return GameTypeSkyrim
	case ".lua":
		return GameTypeLua
	}

	// Detect by content analysis
	contentStr := string(content)

	// Minecraft detection patterns
	if strings.Contains(contentStr, "minecraft:") ||
		strings.Contains(contentStr, "modid") ||
		strings.Contains(contentStr, "block") && strings.Contains(contentStr, "item") {
		return GameTypeMinecraft
	}

	// Skyrim detection patterns
	if strings.Contains(contentStr, "GRUP") ||
		strings.Contains(contentStr, "TES4") ||
		strings.Contains(contentStr, "Skyrim.esm") {
		return GameTypeSkyrim
	}

	// Lua detection patterns
	if strings.Contains(contentStr, "function") ||
		strings.Contains(contentStr, "local") ||
		strings.Contains(contentStr, "require") {
		return GameTypeLua
	}

	return GameTypeUnknown
}

// isMinecraftJSON checks if the JSON content appears to be a Minecraft mod
func isMinecraftJSON(content []byte) bool {
	contentStr := string(content)

	// Look for common Minecraft mod JSON patterns
	minecraftPatterns := []string{
		"minecraft:",
		"modid",
		"forge",
		"fabric",
		"item",
		"block",
		"recipe",
		"texture",
	}

	for _, pattern := range minecraftPatterns {
		if strings.Contains(strings.ToLower(contentStr), pattern) {
			return true
		}
	}

	return false
}

// ExtractMetadata extracts metadata from mod content
func ExtractMetadata(content []byte, gameType GameType) map[string]interface{} {
	metadata := make(map[string]interface{})

	switch gameType {
	case GameTypeMinecraft:
		return extractMinecraftMetadata(content)
	case GameTypeSkyrim:
		return extractSkyrimMetadata(content)
	case GameTypeLua:
		return extractLuaMetadata(content)
	}

	return metadata
}

// extractMinecraftMetadata extracts metadata from Minecraft JSON
func extractMinecraftMetadata(content []byte) map[string]interface{} {
	metadata := make(map[string]interface{})
	contentStr := string(content)

	// Simple string-based extraction (in production, use proper JSON parsing)
	if strings.Contains(contentStr, "modid") {
		metadata["type"] = "forge_mod"
	}
	if strings.Contains(contentStr, "fabric") {
		metadata["type"] = "fabric_mod"
	}
	if strings.Contains(contentStr, "recipe") {
		metadata["has_recipes"] = true
	}
	if strings.Contains(contentStr, "item") {
		metadata["has_items"] = true
	}
	if strings.Contains(contentStr, "block") {
		metadata["has_blocks"] = true
	}

	metadata["format"] = "json"
	return metadata
}

// extractSkyrimMetadata extracts metadata from Skyrim ESP files
func extractSkyrimMetadata(content []byte) map[string]interface{} {
	metadata := make(map[string]interface{})
	metadata["format"] = "esp"
	metadata["game"] = "skyrim"
	return metadata
}

// extractLuaMetadata extracts metadata from Lua scripts
func extractLuaMetadata(content []byte) map[string]interface{} {
	metadata := make(map[string]interface{})
	metadata["format"] = "lua"
	metadata["type"] = "script"
	return metadata
}

// ValidateFile validates an uploaded mod file
func ValidateFile(header *multipart.FileHeader, content []byte) error {
	// Check file size (max 50MB for MVP)
	maxSize := int64(50 * 1024 * 1024) // 50MB
	if header.Size > maxSize {
		return fmt.Errorf("file size exceeds maximum limit of 50MB")
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExtensions := []string{".json", ".esp", ".esm", ".lua", ".txt"}

	isAllowed := false
	for _, allowed := range allowedExtensions {
		if ext == allowed {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return fmt.Errorf("unsupported file type: %s", ext)
	}

	// Detect game type
	gameType := DetectGameType(header, content)
	if gameType == GameTypeUnknown {
		return fmt.Errorf("unable to detect game type from file content")
	}

	return nil
}

// DetectModType detects the mod type from content and filename (simpler version for handlers)
func DetectModType(content []byte, filename string) (string, error) {
	gameType := DetectGameType(&multipart.FileHeader{Filename: filename}, content)
	return string(gameType), nil
}
