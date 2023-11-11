const icons = {
    "Valtan": "/icons/raids/valtan.png",
    "Vykas": "/icons/raids/vykas.png",
    "Kakul Saydon": "/icons/raids/clown.png",
    "Brelshaza": "/icons/raids/brelshaza.png",
    "Kayangel": "/icons/raids/kayangel.png",
    "Akkan": "/icons/raids/akkan.png",
    "": "/icons/raids/surprised.png"
}

const raids = {
    "Valtan": [
        ["Dark Mountain Predator", "Destroyer Lucas", "Leader Lugaru"],
        ["Demon Beast Commander Valtan", "Ravaged Tyrant of Beasts"],
    ],
    "Vykas": [
        ["Incubus Morphe", "Nightmarish Morphe"],
        ["Covetous Devourer Vykas"],
        ["Covetous Legion Commander Vykas"],
    ],
    "Kakul Saydon": [
        ["Saydon"],
        ["Kakul"],
        ["Kakul-Saydon", "Encore-Desiring Kakul-Saydon"],
    ],
    "Brelshaza": [
        ["Gehenna Helkasirs"],
        ["Prokel", "Prokel's Spiritual Echo", "Ashtarot"],
        ["Primordial Nightmare"],
        ["Phantom Legion Commander Brelshaza"],
        ["Brelshaza, Monarch of Nightmares", "Imagined Primordial Nightmare", "Pseudospace Primordial Nightmare"],
        ["Phantom Legion Commander Brelshaza"],
    ],
    "Kayangel": [
        ["Tienis"],
        ["Prunya"],
        ["Lauriel"],
    ],
    "Akkan": [
        ["Griefbringer Maurug", "Evolved Maurug"],
        ["Lord of Degradation Akkan"],
        ["Plague Legion Commander Akkan", "Lord of Kartheon Akkan"],
    ],
}

const raidLookup = {}
for (const [raid, gates] of Object.entries(raids)) {
    for (let gate = 0; gate < gates.length; gate++) {
        for (const boss of gates[gate]) {
            raidLookup[boss] = {
                raid: raid,
                gate: gate + 1,
            }
        }
    }
}

export function getBossIcon(boss: string): string {
    return icons[raidLookup[boss] ? raidLookup[boss].raid : ""]
}

export function getRaidIcon(raid: string): string {
    return icons[raid] ? icons[raid] : icons[""]
}

export function getRaid(boss: string) {
    return raidLookup[boss]
}