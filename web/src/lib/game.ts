const classIds = {
    "Unknown": 0,
    "Warrior (Male)": 101,
    "Berserker": 102,
    "Destroyer": 103,
    "Gunlancer": 104,
    "Paladin": 105,
    "Female Warrior": 111,
    "Slayer": 112,
    "Mage": 201,
    "Arcanist": 202,
    "Summoner": 203,
    "Bard": 204,
    "Sorceress": 205,
    "Martial Artist (Female)": 301,
    "Wardancer": 302,
    "Scrapper": 303,
    "Soulfist": 304,
    "Glaivier": 305,
    "Martial Artist (Male)": 311,
    "Striker": 312,
    "Assassin": 401,
    "Deathblade": 402,
    "Shadowhunter": 403,
    "Reaper": 404,
    "Gunner (Male)": 501,
    "Sharpshooter": 502,
    "Deadeye": 503,
    "Artillerist": 504,
    "Machinist": 505,
    "Gunner (Female)": 511,
    "Gunslinger": 512,
    "Specialist": 601,
    "Artist": 602,
    "Aeromancer": 603,
    "Alchemist": 604
}

export function getClassIcon(c: string): string {
    return '/icons/classes/' + (classIds[c] ?? classIds['Unknown']) + '.png'
}

export function getSkillIcon(icon: string): string {
    return '/icons/skills/' + (icon ? icon : 'unknown.png');
}