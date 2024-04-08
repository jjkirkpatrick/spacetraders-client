package models

type Chart struct {
	WaypointSymbol string `json:"waypointSymbol"`
	SubmittedBy    string `json:"submittedBy"`
	SubmittedOn    string `json:"submittedOn"`
}

type WaypointType string

const (
	GasGiant              WaypointType = "GAS_GIANT"
	Planet                WaypointType = "PLANET"
	Moon                  WaypointType = "MOON"
	OrbitalStation        WaypointType = "ORBITAL_STATION"
	JumpGateWaypoint      WaypointType = "JUMP_GATE"
	AsteroidField         WaypointType = "ASTEROID_FIELD"
	Asteroid              WaypointType = "ASTEROID"
	EngineeredAsteroid    WaypointType = "ENGINEERED_ASTEROID"
	AsteroidBase          WaypointType = "ASTEROID_BASE"
	Nebular               WaypointType = "NEBULAR"
	DebrisField           WaypointType = "DEBRIS_FIELD"
	GravityWell           WaypointType = "GRAVITY_WELL"
	ArtificialGravityWell WaypointType = "ARTIFICIAL_GRAVITY_WELL"
	FuelStation           WaypointType = "FUEL_STATION"
)

type WaypointTrait string

const (
	TraitUncharted             WaypointTrait = "UNCHARTED"
	TraitUnderConstruction     WaypointTrait = "UNDER_CONSTRUCTION"
	TraitMarketplace           WaypointTrait = "MARKETPLACE"
	TraitShipyard              WaypointTrait = "SHIPYARD"
	TraitOutpost               WaypointTrait = "OUTPOST"
	TraitScatteredSettlements  WaypointTrait = "SCATTERED_SETTLEMENTS"
	TraitSprawlingCities       WaypointTrait = "SPRAWLING_CITIES"
	TraitMegaStructures        WaypointTrait = "MEGA_STRUCTURES"
	TraitPirateBase            WaypointTrait = "PIRATE_BASE"
	TraitOvercrowded           WaypointTrait = "OVERCROWDED"
	TraitHighTech              WaypointTrait = "HIGH_TECH"
	TraitCorrupt               WaypointTrait = "CORRUPT"
	TraitBureaucratic          WaypointTrait = "BUREAUCRATIC"
	TraitTradingHub            WaypointTrait = "TRADING_HUB"
	TraitIndustrial            WaypointTrait = "INDUSTRIAL"
	TraitBlackMarket           WaypointTrait = "BLACK_MARKET"
	TraitResearchFacility      WaypointTrait = "RESEARCH_FACILITY"
	TraitMilitaryBase          WaypointTrait = "MILITARY_BASE"
	TraitSurveillanceOutpost   WaypointTrait = "SURVEILLANCE_OUTPOST"
	TraitExplorationOutpost    WaypointTrait = "EXPLORATION_OUTPOST"
	TraitMineralDeposits       WaypointTrait = "MINERAL_DEPOSITS"
	TraitCommonMetalDeposits   WaypointTrait = "COMMON_METAL_DEPOSITS"
	TraitPreciousMetalDeposits WaypointTrait = "PRECIOUS_METAL_DEPOSITS"
	TraitRareMetalDeposits     WaypointTrait = "RARE_METAL_DEPOSITS"
	TraitMethanePools          WaypointTrait = "METHANE_POOLS"
	TraitIceCrystals           WaypointTrait = "ICE_CRYSTALS"
	TraitExplosiveGases        WaypointTrait = "EXPLOSIVE_GASES"
	TraitStrongMagnetosphere   WaypointTrait = "STRONG_MAGNETOSPHERE"
	TraitVibrantAuroras        WaypointTrait = "VIBRANT_AURORAS"
	TraitSaltFlats             WaypointTrait = "SALT_FLATS"
	TraitCanyons               WaypointTrait = "CANYONS"
	TraitPerpetualDaylight     WaypointTrait = "PERPETUAL_DAYLIGHT"
	TraitPerpetualOvercast     WaypointTrait = "PERPETUAL_OVERCAST"
	TraitDrySeabeds            WaypointTrait = "DRY_SEABEDS"
	TraitMagmaSeas             WaypointTrait = "MAGMA_SEAS"
	TraitSupervolcanoes        WaypointTrait = "SUPERVOLCANOES"
	TraitAshClouds             WaypointTrait = "ASH_CLOUDS"
	TraitVastRuins             WaypointTrait = "VAST_RUINS"
	TraitMutatedFlora          WaypointTrait = "MUTATED_FLORA"
	TraitTerraformed           WaypointTrait = "TERRAFORMED"
	TraitExtremeTemperatures   WaypointTrait = "EXTREME_TEMPERATURES"
	TraitExtremePressure       WaypointTrait = "EXTREME_PRESSURE"
	TraitDiverseLife           WaypointTrait = "DIVERSE_LIFE"
	TraitScarceLife            WaypointTrait = "SCARCE_LIFE"
	TraitFossils               WaypointTrait = "FOSSILS"
	TraitWeakGravity           WaypointTrait = "WEAK_GRAVITY"
	TraitStrongGravity         WaypointTrait = "STRONG_GRAVITY"
	TraitCrushingGravity       WaypointTrait = "CRUSHING_GRAVITY"
	TraitToxicAtmosphere       WaypointTrait = "TOXIC_ATMOSPHERE"
	TraitCorrosiveAtmosphere   WaypointTrait = "CORROSIVE_ATMOSPHERE"
	TraitBreathableAtmosphere  WaypointTrait = "BREATHABLE_ATMOSPHERE"
	TraitThinAtmosphere        WaypointTrait = "THIN_ATMOSPHERE"
	TraitJovian                WaypointTrait = "JOVIAN"
	TraitRocky                 WaypointTrait = "ROCKY"
	TraitVolcanic              WaypointTrait = "VOLCANIC"
	TraitFrozen                WaypointTrait = "FROZEN"
	TraitSwamp                 WaypointTrait = "SWAMP"
	TraitBarren                WaypointTrait = "BARREN"
	TraitTemperate             WaypointTrait = "TEMPERATE"
	TraitJungle                WaypointTrait = "JUNGLE"
	TraitOcean                 WaypointTrait = "OCEAN"
	TraitRadioactive           WaypointTrait = "RADIOACTIVE"
	TraitMicroGravityAnomalies WaypointTrait = "MICRO_GRAVITY_ANOMALIES"
	TraitDebrisCluster         WaypointTrait = "DEBRIS_CLUSTER"
	TraitDeepCraters           WaypointTrait = "DEEP_CRATERS"
	TraitShallowCraters        WaypointTrait = "SHALLOW_CRATERS"
	TraitUnstableComposition   WaypointTrait = "UNSTABLE_COMPOSITION"
	TraitHollowedInterior      WaypointTrait = "HOLLOWED_INTERIOR"
	TraitStripped              WaypointTrait = "STRIPPED"
)

type WaypoinTraits struct {
	Symbol      WaypointTrait `json:"symbol"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
}
