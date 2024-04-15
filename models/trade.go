package models

type GoodSymbol string

const (
	PreciousStones          GoodSymbol = "PRECIOUS_STONES"
	QuartzSand              GoodSymbol = "QUARTZ_SAND"
	SiliconCrystals         GoodSymbol = "SILICON_CRYSTALS"
	AmmoniaIce              GoodSymbol = "AMMONIA_ICE"
	LiquidHydrogen          GoodSymbol = "LIQUID_HYDROGEN"
	LiquidNitrogen          GoodSymbol = "LIQUID_NITROGEN"
	IceWater                GoodSymbol = "ICE_WATER"
	ExoticMatter            GoodSymbol = "EXOTIC_MATTER"
	AdvancedCircuitry       GoodSymbol = "ADVANCED_CIRCUITRY"
	GravitonEmitters        GoodSymbol = "GRAVITON_EMITTERS"
	Iron                    GoodSymbol = "IRON"
	IronOre                 GoodSymbol = "IRON_ORE"
	Copper                  GoodSymbol = "COPPER"
	CopperOre               GoodSymbol = "COPPER_ORE"
	Aluminum                GoodSymbol = "ALUMINUM"
	AluminumOre             GoodSymbol = "ALUMINUM_ORE"
	Silver                  GoodSymbol = "SILVER"
	SilverOre               GoodSymbol = "SILVER_ORE"
	Gold                    GoodSymbol = "GOLD"
	GoldOre                 GoodSymbol = "GOLD_ORE"
	Platinum                GoodSymbol = "PLATINUM"
	PlatinumOre             GoodSymbol = "PLATINUM_ORE"
	Diamonds                GoodSymbol = "DIAMONDS"
	Uranite                 GoodSymbol = "URANITE"
	UraniteOre              GoodSymbol = "URANITE_ORE"
	Meritium                GoodSymbol = "MERITIUM"
	MeritiumOre             GoodSymbol = "MERITIUM_ORE"
	Hydrocarbon             GoodSymbol = "HYDROCARBON"
	Antimatter              GoodSymbol = "ANTIMATTER"
	FabMats                 GoodSymbol = "FAB_MATS"
	Fertilizers             GoodSymbol = "FERTILIZERS"
	Fabrics                 GoodSymbol = "FABRICS"
	Food                    GoodSymbol = "FOOD"
	Jewelry                 GoodSymbol = "JEWELRY"
	Machinery               GoodSymbol = "MACHINERY"
	Firearms                GoodSymbol = "FIREARMS"
	AssaultRifles           GoodSymbol = "ASSAULT_RIFLES"
	MilitaryEquipment       GoodSymbol = "MILITARY_EQUIPMENT"
	Explosives              GoodSymbol = "EXPLOSIVES"
	LabInstruments          GoodSymbol = "LAB_INSTRUMENTS"
	Ammunition              GoodSymbol = "AMMUNITION"
	Electronics             GoodSymbol = "ELECTRONICS"
	ShipPlating             GoodSymbol = "SHIP_PLATING"
	ShipParts               GoodSymbol = "SHIP_PARTS"
	Equipment               GoodSymbol = "EQUIPMENT"
	Fuel                    GoodSymbol = "FUEL"
	Medicine                GoodSymbol = "MEDICINE"
	Drugs                   GoodSymbol = "DRUGS"
	Clothing                GoodSymbol = "CLOTHING"
	Microprocessors         GoodSymbol = "MICROPROCESSORS"
	Plastics                GoodSymbol = "PLASTICS"
	Polynucleotides         GoodSymbol = "POLYNUCLEOTIDES"
	Biocomposites           GoodSymbol = "BIOCOMPOSITES"
	QuantumStabilizers      GoodSymbol = "QUANTUM_STABILIZERS"
	Nanobots                GoodSymbol = "NANOBOTS"
	AiMainframes            GoodSymbol = "AI_MAINFRAMES"
	QuantumDrives           GoodSymbol = "QUANTUM_DRIVES"
	RoboticDrones           GoodSymbol = "ROBOTIC_DRONES"
	CyberImplants           GoodSymbol = "CYBER_IMPLANTS"
	GeneTherapeutics        GoodSymbol = "GENE_THERAPEUTICS"
	NeuralChips             GoodSymbol = "NEURAL_CHIPS"
	MoodRegulators          GoodSymbol = "MOOD_REGULATORS"
	ViralAgents             GoodSymbol = "VIRAL_AGENTS"
	MicroFusionGenerators   GoodSymbol = "MICRO_FUSION_GENERATORS"
	Supergrains             GoodSymbol = "SUPERGRAINS"
	LaserRifles             GoodSymbol = "LASER_RIFLES"
	Holographics            GoodSymbol = "HOLOGRAPHICS"
	ShipSalvage             GoodSymbol = "SHIP_SALVAGE"
	RelicTech               GoodSymbol = "RELIC_TECH"
	NovelLifeforms          GoodSymbol = "NOVEL_LIFEFORMS"
	BotanicalSpecimens      GoodSymbol = "BOTANICAL_SPECIMENS"
	CulturalArtifacts       GoodSymbol = "CULTURAL_ARTIFACTS"
	FrameProbe              GoodSymbol = "FRAME_PROBE"
	FrameDrone              GoodSymbol = "FRAME_DRONE"
	FrameInterceptor        GoodSymbol = "FRAME_INTERCEPTOR"
	FrameRacer              GoodSymbol = "FRAME_RACER"
	FrameFighter            GoodSymbol = "FRAME_FIGHTER"
	FrameFrigate            GoodSymbol = "FRAME_FRIGATE"
	FrameShuttle            GoodSymbol = "FRAME_SHUTTLE"
	FrameExplorer           GoodSymbol = "FRAME_EXPLORER"
	FrameMiner              GoodSymbol = "FRAME_MINER"
	FrameLightFreighter     GoodSymbol = "FRAME_LIGHT_FREIGHTER"
	FrameHeavyFreighter     GoodSymbol = "FRAME_HEAVY_FREIGHTER"
	FrameTransport          GoodSymbol = "FRAME_TRANSPORT"
	FrameDestroyer          GoodSymbol = "FRAME_DESTROYER"
	FrameCruiser            GoodSymbol = "FRAME_CRUISER"
	FrameCarrier            GoodSymbol = "FRAME_CARRIER"
	ReactorSolarI           GoodSymbol = "REACTOR_SOLAR_I"
	ReactorFusionI          GoodSymbol = "REACTOR_FUSION_I"
	ReactorFissionI         GoodSymbol = "REACTOR_FISSION_I"
	ReactorChemicalI        GoodSymbol = "REACTOR_CHEMICAL_I"
	ReactorAntimatterI      GoodSymbol = "REACTOR_ANTIMATTER_I"
	EngineImpulseDriveI     GoodSymbol = "ENGINE_IMPULSE_DRIVE_I"
	EngineIonDriveI         GoodSymbol = "ENGINE_ION_DRIVE_I"
	EngineIonDriveII        GoodSymbol = "ENGINE_ION_DRIVE_II"
	EngineHyperDriveI       GoodSymbol = "ENGINE_HYPER_DRIVE_I"
	ModuleMineralProcessorI GoodSymbol = "MODULE_MINERAL_PROCESSOR_I"
	ModuleGasProcessorI     GoodSymbol = "MODULE_GAS_PROCESSOR_I"
	ModuleCargoHoldI        GoodSymbol = "MODULE_CARGO_HOLD_I"
	ModuleCargoHoldII       GoodSymbol = "MODULE_CARGO_HOLD_II"
	ModuleCargoHoldIII      GoodSymbol = "MODULE_CARGO_HOLD_III"
	ModuleCrewQuartersI     GoodSymbol = "MODULE_CREW_QUARTERS_I"
	ModuleEnvoyQuartersI    GoodSymbol = "MODULE_ENVOY_QUARTERS_I"
	ModulePassengerCabinI   GoodSymbol = "MODULE_PASSENGER_CABIN_I"
	ModuleMicroRefineryI    GoodSymbol = "MODULE_MICRO_REFINERY_I"
	ModuleScienceLabI       GoodSymbol = "MODULE_SCIENCE_LAB_I"
	ModuleJumpDriveI        GoodSymbol = "MODULE_JUMP_DRIVE_I"
	ModuleJumpDriveII       GoodSymbol = "MODULE_JUMP_DRIVE_II"
	ModuleJumpDriveIII      GoodSymbol = "MODULE_JUMP_DRIVE_III"
	ModuleWarpDriveI        GoodSymbol = "MODULE_WARP_DRIVE_I"
	ModuleWarpDriveII       GoodSymbol = "MODULE_WARP_DRIVE_II"
	ModuleWarpDriveIII      GoodSymbol = "MODULE_WARP_DRIVE_III"
	ModuleShieldGeneratorI  GoodSymbol = "MODULE_SHIELD_GENERATOR_I"
	ModuleShieldGeneratorII GoodSymbol = "MODULE_SHIELD_GENERATOR_II"
	ModuleOreRefineryI      GoodSymbol = "MODULE_ORE_REFINERY_I"
	ModuleFuelRefineryI     GoodSymbol = "MODULE_FUEL_REFINERY_I"
)

type Good struct {
	Symbol      GoodSymbol `json:"symbol"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
}

type Market struct {
	Symbol       string             `json:"symbol"`
	Exports      []Good             `json:"exports"`
	Imports      []Good             `json:"imports"`
	Exchange     []Good             `json:"exchange"`
	Transactions []Transaction      `json:"transactions"`
	TradeGoods   []MarketTradeGoods `json:"tradeGoods"`
}

type Extraction struct {
	ShipSymbol string `json:"shipSymbol" `
	Yield      Yield  `json:"yield" `
}

type Yield struct {
	Symbol GoodSymbol `json:"symbol" `
	Units  int        `json:"units" `
}

type MarketTradeGoodType string

const (
	Export   MarketTradeGoodType = "EXPORT"
	Import   MarketTradeGoodType = "IMPORT"
	Exchange MarketTradeGoodType = "EXCHANGE"
)

type MarketTradeSupply string

const (
	Scarse   MarketTradeSupply = "SCARSE"
	Limited  MarketTradeSupply = "LIMITED"
	Moderate MarketTradeSupply = "MODERATE"
	High     MarketTradeSupply = "HIGH"
	Abundant MarketTradeSupply = "ABUNDANT"
)

type MarketTradeAvtivity string

const (
	Weak       MarketTradeAvtivity = "WEAK"
	Growing    MarketTradeAvtivity = "GROWING"
	Strong     MarketTradeAvtivity = "STRONG"
	Restricted MarketTradeAvtivity = "RESTRICTED"
)

type MarketTradeGoods struct {
	Symbol        GoodSymbol          `json:"symbol"`
	Type          MarketTradeGoodType `json:"type"`
	TradeValue    int                 `json:"tradeValue"`
	Supply        MarketTradeSupply   `json:"supply"`
	Activity      MarketTradeAvtivity `json:"activity"`
	PurchasePrice int                 `json:"purchasePrice"`
	SellPrice     int                 `json:"sellPrice"`
}
