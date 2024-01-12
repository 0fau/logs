<script lang="ts">
    import {getBossIcon} from "$lib/raids";
    import {formatDamage, formatDate, formatDuration} from "$lib/print";

    const difficultyColors = {
        "Inferno": "#9a3148",
        "Trial": "#9a3148",
        "Challenge": "#625f77",
        "Hard": "#b9982e",
        "Normal": "#625f77",
    };

    function shortBossName(boss: string): string {
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

    export let encounter;
    export let width;
</script>

<div class="h-[80px] {width} flex border-[0.5px] border-[#c58597] shadow-sm rounded-md bg-[#F4EDE9]">
    <div class="w-full h-full flex flex-row ml-5 items-center">
        <div>
            <div class="self-start text-left text-[#575279]">
                <div>
                    <span class="font-medium">[#{encounter.id}]</span>
                    <img alt={encounter.boss} src={getBossIcon(encounter.boss)}
                         class="inline w-6 h-6 -translate-y-0.5"/>
                    <span class="font-medium">{shortBossName(encounter.boss)}</span>
                </div>
                <p class="text-sm">{formatDamage(encounter.damage)} damage dealt
                    in {formatDuration(encounter.duration)}</p>
                <p class="text-xs text-[#5d5978]"><span class="font-semibold"
                                                        style="color: {difficultyColors[encounter.difficulty]}">{encounter.difficulty}</span> {formatDate(encounter.date)}
                </p>
            </div>
        </div>
        <div class="py-1 px-1.5 h-full ml-auto self-end flex flex-col rounded-r-md text-white">
            <span class="text-xs text-center self-end text-[#F4EDE9] p-0.5 px-1 mr-0.5 mt-1.5 rounded-sm bg-[#b4637a] font-medium">{encounter.anonymized ? encounter.players[encounter.localPlayer].class + " " + encounter.localPlayer : encounter.localPlayer}</span>
            <span class="text-xs text-[#b4637a] self-end text-right mr-0.5 mt-0.5 font-medium">{encounter.players[encounter.localPlayer].class}</span>
            <span class="text-[#575279] text-right mr-1 my-auto text-lg font-medium">{formatDamage(encounter.players[encounter.localPlayer].dps)}</span>
        </div>
    </div>
</div>