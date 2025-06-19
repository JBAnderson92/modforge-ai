CREATE TABLE users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    firebase_uid TEXT UNIQUE NOT NULL,
    display_name TEXT,
    credits INTEGER DEFAULT 0,
    plan TEXT DEFAULT 'free',
    monthly_jobs_used INTEGER DEFAULT 0,
    monthly_jobs_reset_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE mod_jobs (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    game_type TEXT NOT NULL,
    original_filename TEXT NOT NULL,
    original_file_size INTEGER NOT NULL,
    original_file_url TEXT,
    processed_file_url TEXT,
    preset_type TEXT NOT NULL,
    ai_prompt TEXT,
    ai_response TEXT,
    changelog TEXT,
    tokens_used INTEGER DEFAULT 0,
    credits_used INTEGER DEFAULT 0,
    error_message TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE mod_presets (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    game_type TEXT NOT NULL,
    prompt_template TEXT NOT NULL,
    credit_cost INTEGER DEFAULT 1,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    firebase_token TEXT NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- Insert default presets for Minecraft
INSERT INTO mod_presets (id, name, description, game_type, prompt_template, credit_cost) VALUES
('minecraft_lore_friendly', 'Make Lore-Friendly', 'Rewrite item names and descriptions to be more lore-friendly and immersive', 'minecraft', 'Rewrite the following Minecraft mod JSON to make all item names and descriptions more lore-friendly and immersive. Keep the technical structure intact but make the content feel more natural and engaging: {content}', 1),
('minecraft_balance', 'Balance Items', 'Rebalance item stats, durability, and crafting recipes for better gameplay', 'minecraft', 'Rebalance the following Minecraft mod JSON to improve gameplay balance. Adjust damage values, durability, crafting costs, and rarity appropriately: {content}', 2),
('minecraft_translate', 'Translate Mod', 'Translate mod content to different languages', 'minecraft', 'Translate the following Minecraft mod JSON content to {target_language}. Keep all technical JSON structure and keys intact, only translate the displayable text content: {content}', 1),
('minecraft_expand', 'Add Content', 'Generate additional items, blocks, or recipes that fit the mod theme', 'minecraft', 'Analyze the following Minecraft mod JSON and generate 3-5 additional items, blocks, or recipes that would fit well with the existing mod theme. Maintain the same JSON structure and style: {content}', 3);
