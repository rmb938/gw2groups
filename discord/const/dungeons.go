package _const

var DungeonIDs = []string{
	"ascalonian_catacombs",
	"caudecus_manor",
	"twilight_arbor",
	"sorrows_embrace",
	"citadel_of_flame",
	"honor_of_the_waves",
	"crucible_of_eternity",
	"ruined_city_of_arah",
}

var DungeonsIDsToName = map[string]string{
	"ascalonian_catacombs": "Ascalonian Catacombs",
	"caudecus_manor":       "Caudecus's Manor",
	"twilight_arbor":       "Twilight Arbor",
	"sorrows_embrace":      "Sorrow's Embrace",
	"citadel_of_flame":     "Citadel of Flame",
	"honor_of_the_waves":   "Honor of the Waves",
	"crucible_of_eternity": "Crucible of Eternity",
	"ruined_city_of_arah":  "The Ruined City of Arah",
}

var DungeonsPaths = map[string][]string{
	"ascalonian_catacombs": {"Hodgins", "Detha", "tzark"},
	"caudecus_manor":       {"Asura", "Seraph", "Butler"},
	"twilight_arbor":       {"Leurent", "Vevina", "Aetherpath"},
	"sorrows_embrace":      {"Fergg", "Rasolov", "Koptev"},
	"citadel_of_flame":     {"Ferrah", "Magg", "Rhiannon"},
	"honor_of_the_waves":   {"Butcher", "Plunderer", "Zealot"},
	"crucible_of_eternity": {"Submarine", "Teleporter", "Front door"},
	"ruined_city_of_arah":  {"Jotun", "Mursaat", "Forgotten", "Seer"},
}
