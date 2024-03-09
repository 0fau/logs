const icons = {
    "Valtan": "/icons/raids/valtan.png",
    "Vykas": "/icons/raids/vykas.png",
    "Kakul Saydon": "/icons/raids/clown.png",
    "Brelshaza": "/icons/raids/brelshaza.png",
    "Kayangel": "/icons/raids/kayangel.png",
    "Akkan": "/icons/raids/akkan.png",
    "Ivory": "/icons/raids/tower.png",
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
    "Ivory": [
        ["Kaltaya, the Blooming Chaos"],
        ["Rakathus, the Lurking Arrogance"],
        ["Firehorn, Trampler of Earth"],
        ["Lazaram, the Trailblazer", "Subordinated Vertus", "Subordinated Calventus", "Subordinated Legoros", "Brand of Subordination"],
    ],
}

export function shortBossName(boss: string): string {
    if (boss === "Phantom Legion Commander Brelshaza") {
        return "Phantom LC Brelshaza"
    }

    if (boss === "Covetous Legion Commander Vykas") {
        return "Covetous LC Vykas"
    }

    if (boss === "Plague Legion Commander Akkan") {
        return "Plague LC Akkan"
    }

    if (boss === "Lazaram, the Trailblazer") {
        return "Lazaram"
    }

    if (boss === "Firehorn, Trampler of Earth") {
        return "Firehorn"
    }

    if (boss === "Rakathus, the Lurking Arrogance") {
        return "Rakathus"
    }

    if (boss === "Kaltaya, the Blooming Chaos") {
        return "Kaltaya"
    }

    return boss
}

const guardians = [
    "Gargadeth", "Sonavel", "Hanumatan", "Caliligos", "Deskaluda",
]

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

export function isGuardian(boss: string) {
    return guardians.includes(boss)
}

export const bossHpMap = {
    "Dark Mountain Predator": 50,
    "Destroyer Lucas": 50,
    "Leader Lugaru": 50,
    "Demon Beast Commander Valtan": 160,
    "Ravaged Tyrant of Beasts": 40,
    "Incubus Morphe": 60,
    "Nightmarish Morphe": 60,
    "Covetous Devourer Vykas": 160,
    "Covetous Legion Commander Vykas": 180,
    "Saydon": 160,
    "Kakul": 140,
    "Kakul-Saydon": 180,
    "Encore-Desiring Kakul-Saydon": 77,
    "Gehenna Helkasirs": 120,
    "Ashtarot": 170,
    "Primordial Nightmare": 190,
    "Brelshaza, Monarch of Nightmares": 200,
    "Imagined Primordial Nightmare": 20,
    "Pseudospace Primordial Nightmare": 20,
    "Phantom Legion Commander Brelshaza": 250,
    "Griefbringer Maurug": 150,
    "Evolved Maurug": 30,
    "Lord of Degradation Akkan": 190,
    "Plague Legion Commander Akkan": 220,
    "Lord of Kartheon Akkan": 300,
    "Tienis": 110,
    "Celestial Sentinel" : 60,
    "Prunya": 90,
    "Lauriel": 200,
    "Kaltaya, the Blooming Chaos": 120,
    "Rakathus, the Lurking Arrogance": 160,
    "Firehorn, Trampler of Earth": 160,
    "Lazaram, the Trailblazer": 200
};