package models

type GoodSymbol string

const (
	PreciousStones        GoodSymbol = "PRECIOUS_STONES"
	QuartzSand            GoodSymbol = "QUARTZ_SAND"
	SiliconCrystals       GoodSymbol = "SILICON_CRYSTALS"
	AmmoniaIce            GoodSymbol = "AMMONIA_ICE"
	LiquidHydrogen        GoodSymbol = "LIQUID_HYDROGEN"
	LiquidNitrogen        GoodSymbol = "LIQUID_NITROGEN"
	IceWater              GoodSymbol = "ICE_WATER"
	ExoticMatter          GoodSymbol = "EXOTIC_MATTER"
	AdvancedCircuitry     GoodSymbol = "ADVANCED_CIRCUITRY"
	GravitonEmitters      GoodSymbol = "GRAVITON_EMITTERS"
	Iron                  GoodSymbol = "IRON"
	IronOre               GoodSymbol = "IRON_ORE"
	Copper                GoodSymbol = "COPPER"
	CopperOre             GoodSymbol = "COPPER_ORE"
	Aluminum              GoodSymbol = "ALUMINUM"
	AluminumOre           GoodSymbol = "ALUMINUM_ORE"
	Silver                GoodSymbol = "SILVER"
	SilverOre             GoodSymbol = "SILVER_ORE"
	Gold                  GoodSymbol = "GOLD"
	GoldOre               GoodSymbol = "GOLD_ORE"
	Platinum              GoodSymbol = "PLATINUM"
	PlatinumOre           GoodSymbol = "PLATINUM_ORE"
	Diamonds              GoodSymbol = "DIAMONDS"
	Uranite               GoodSymbol = "URANITE"
	UraniteOre            GoodSymbol = "URANITE_ORE"
	Meritium              GoodSymbol = "MERITIUM"
	MeritiumOre           GoodSymbol = "MERITIUM_ORE"
	Hydrocarbon           GoodSymbol = "HYDROCARBON"
	Antimatter            GoodSymbol = "ANTIMATTER"
	FabMats               GoodSymbol = "FAB_MATS"
	Fertilizers           GoodSymbol = "FERTILIZERS"
	Fabrics               GoodSymbol = "FABRICS"
	Food                  GoodSymbol = "FOOD"
	Jewelry               GoodSymbol = "JEWELRY"
	Machinery             GoodSymbol = "MACHINERY"
	Firearms              GoodSymbol = "FIREARMS"
	AssaultRifles         GoodSymbol = "ASSAULT_RIFLES"
	MilitaryEquipment     GoodSymbol = "MILITARY_EQUIPMENT"
	Explosives            GoodSymbol = "EXPLOSIVES"
	LabInstruments        GoodSymbol = "LAB_INSTRUMENTS"
	Ammunition            GoodSymbol = "AMMUNITION"
	Electronics           GoodSymbol = "ELECTRONICS"
	ShipPlating           GoodSymbol = "SHIP_PLATING"
	ShipParts             GoodSymbol = "SHIP_PARTS"
	Equipment             GoodSymbol = "EQUIPMENT"
	Fuel                  GoodSymbol = "FUEL"
	Medicine              GoodSymbol = "MEDICINE"
	Drugs                 GoodSymbol = "DRUGS"
	Clothing              GoodSymbol = "CLOTHING"
	Microprocessors       GoodSymbol = "MICROPROCESSORS"
	Plastics              GoodSymbol = "PLASTICS"
	Polynucleotides       GoodSymbol = "POLYNUCLEOTIDES"
	Biocomposites         GoodSymbol = "BIOCOMPOSITES"
	QuantumStabilizers    GoodSymbol = "QUANTUM_STABILIZERS"
	Nanobots              GoodSymbol = "NANOBOTS"
	AiMainframes          GoodSymbol = "AI_MAINFRAMES"
	QuantumDrives         GoodSymbol = "QUANTUM_DRIVES"
	RoboticDrones         GoodSymbol = "ROBOTIC_DRONES"
	CyberImplants         GoodSymbol = "CYBER_IMPLANTS"
	GeneTherapeutics      GoodSymbol = "GENE_THERAPEUTICS"
	NeuralChips           GoodSymbol = "NEURAL_CHIPS"
	MoodRegulators        GoodSymbol = "MOOD_REGULATORS"
	ViralAgents           GoodSymbol = "VIRAL_AGENTS"
	MicroFusionGenerators GoodSymbol = "MICRO_FUSION_GENERATORS"
	Supergrains           GoodSymbol = "SUPERGRAINS"
	LaserRifles           GoodSymbol = "LASER_RIFLES"
	Holographics          GoodSymbol = "HOLOGRAPHICS"
	ShipSalvage           GoodSymbol = "SHIP_SALVAGE"
	RelicTech             GoodSymbol = "RELIC_TECH"
	NovelLifeforms        GoodSymbol = "NOVEL_LIFEFORMS"
	BotanicalSpecimens    GoodSymbol = "BOTANICAL_SPECIMENS"
	CulturalArtifacts     GoodSymbol = "CULTURAL_ARTIFACTS"
	// And so on for all allowed values...
)

type Good struct {
	Symbol      GoodSymbol `json:"symbol"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
}

type Market struct {
	Symbol   string `json:"symbol"`
	Exports  []Good `json:"exports"`
	Imports  []Good `json:"imports"`
	Exchange []Good `json:"exchange"`
	// Consider adding other fields as necessary
}
