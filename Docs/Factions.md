# Faction Operations Guide

This guide covers faction-related operations using the `entities` package.

## Getting Started

```go
import (
    "github.com/jjkirkpatrick/spacetraders-client/client"
    "github.com/jjkirkpatrick/spacetraders-client/entities"
)

// Create client
c, err := client.NewClient(options)
defer c.Close(ctx)

// Get factions
factions, err := entities.ListFactions(c)
```

## Functions

### ListFactions

Fetches all factions in the game.

```go
func ListFactions(c *client.Client) ([]*Faction, error)
```

**Example:**
```go
factions, err := entities.ListFactions(c)
if err != nil {
    log.Fatalf("Failed to list factions: %v", err)
}

for _, faction := range factions {
    fmt.Printf("Faction: %s - %s\n", faction.Symbol, faction.Name)
    fmt.Printf("  Headquarters: %s\n", faction.Headquarters)
    fmt.Printf("  Recruiting: %v\n", faction.IsRecruiting)
}
```

### GetFaction

Retrieves detailed information about a specific faction.

```go
func GetFaction(c *client.Client, symbol string) (*Faction, error)
```

**Example:**
```go
faction, err := entities.GetFaction(c, "COSMIC")
if err != nil {
    log.Fatalf("Failed to get faction: %v", err)
}

fmt.Printf("Name: %s\n", faction.Name)
fmt.Printf("Description: %s\n", faction.Description)
fmt.Printf("Traits:\n")
for _, trait := range faction.Traits {
    fmt.Printf("  - %s: %s\n", trait.Symbol, trait.Name)
}
```

## Faction Structure

| Field | Type | Description |
|-------|------|-------------|
| `Symbol` | string | Unique faction identifier |
| `Name` | string | Display name of the faction |
| `Description` | string | Faction description and lore |
| `Headquarters` | string | Faction headquarters waypoint |
| `Traits` | []Trait | List of faction traits |
| `IsRecruiting` | bool | Whether faction accepts new agents |

## Available Factions

When starting a new game, you can choose from these factions:

| Symbol | Name | Description |
|--------|------|-------------|
| `COSMIC` | Cosmic Engineers | Builders and innovators |
| `VOID` | Voidfarers | Explorers of the unknown |
| `GALACTIC` | Galactic Alliance | Democratic federation |
| `QUANTUM` | Quantum Federation | Scientific researchers |
| `DOMINION` | Stellar Dominion | Expansionist empire |
| `ASTRO` | Astro-Salvage Union | Resource reclaimers |
| `CORSAIRS` | Corsair Collective | Independent traders |
| `OBSIDIAN` | Obsidian Syndicate | Shadow operators |
| `AEGIS` | Aegis Collective | Peacekeepers |
| `UNITED` | United Systems | Economic coalition |
| `SOLITARY` | Solitary Wanderers | Independent agents |
| `COBALT` | Cobalt Industries | Mining conglomerate |
| `OMEGA` | Omega Syndicate | Underground network |
| `ECHO` | Echo Collective | Communication experts |
| `LORDS` | Lords of the Galaxy | Ancient rulers |
| `CULT` | Cult of the Stars | Mystic believers |
| `ANCIENTS` | Ancient Ascendancy | Knowledge seekers |
| `SHADOW` | Shadow Covenant | Covert operatives |
| `ETHEREAL` | Ethereal Enclave | Transcendent beings |

## Faction Traits

Faction traits affect gameplay and interactions:

| Trait | Description |
|-------|-------------|
| `BUREAUCRATIC` | Slower but more reliable |
| `SECRETIVE` | Limited public information |
| `CAPITALISTIC` | Trade-focused economy |
| `INDUSTRIOUS` | Enhanced production |
| `PEACEFUL` | Non-aggressive policies |
| `REBELLIOUS` | Anti-establishment |
| `WELCOMING` | Easy recruitment |
| `SMUGGLERS` | Black market access |
| `SCAVENGERS` | Salvage specialists |
| `INNOVATIVE` | Advanced technology |
| `BOLD` | Risk-taking culture |
| `VISIONARY` | Long-term planning |
| `CURIOUS` | Exploration-focused |
| `DARING` | Combat-ready |
| `EXPLORATORY` | Deep space focus |
| `RESOURCEFUL` | Efficient operations |
| `FLEXIBLE` | Adaptable strategies |
| `COOPERATIVE` | Alliance-building |
| `UNITED` | Strong internal cohesion |
| `STRATEGIC` | Military planning |
| `INTELLIGENT` | Research bonuses |
| `RESEARCH_FOCUSED` | Science priority |
| `COLLABORATIVE` | Team operations |
| `PROGRESSIVE` | Modern approaches |
| `MILITARISTIC` | Combat strength |
| `TECHNOLOGICALLY_ADVANCED` | High-tech equipment |
| `AGGRESSIVE` | Offensive posture |
| `IMPERIALISTIC` | Expansion-focused |
| `TREASURE_HUNTERS` | Rare find bonuses |
| `DEXTEROUS` | Nimble operations |
| `UNPREDICTABLE` | Variable behavior |
| `BRUTAL` | Harsh tactics |
| `FLEETING` | Quick operations |
| `ADAPTABLE` | Flexible responses |
| `SELF_SUFFICIENT` | Independent operations |
| `DEFENSIVE` | Fortified positions |
| `PROUD` | High standards |
| `DIVERSE` | Varied approaches |
| `INDEPENDENT` | Autonomous agents |
| `SELF_INTERESTED` | Individual focus |
| `FRAGMENTED` | Decentralized |
| `COMMERCIAL` | Trade emphasis |
| `FREE_MARKETS` | Open economy |
| `ENTREPRENEURIAL` | Business innovation |
