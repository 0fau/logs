<script>
    import {formatDamage, formatPercentFlat} from "$lib/components/meter/print.ts";

    let partyTextColors = ["text-[#a2596d]", "text-[#5e5277]"]
    let partyBgColors = ["bg-[#a2596d]", "bg-[#5e5277]"]

    export let focused;
    let tab = 0;

    function setTab(num) {
        tab = num
    }

    function getParties(encounter) {
        if (encounter.broken || encounter.parties.length === 0) {
            return [Object.keys(encounter.players)]
        }
        return encounter.parties
    }

    function sortByDPS(encounter, players) {
        players.sort((a, b) => {
            return encounter.players[b].dps - encounter.players[a].dps
        })
        return players
    }

    function getDamagePercent(player) {
        let percent = focused.players[player].damage / focused.max;
        if (percent < 0.08) {
            return "0"
        }

        return formatPercentFlat(percent)
    }
</script>


<div class="h-[400px] overflow-x-scroll px-2 mb-2 flex justify-center items-center border-[0.5px] border-[#c58597] shadow-sm rounded-md w-[94%] bg-[#f7f2ef]">
    <div class="w-[100%] h-full flex flex-col items-center justify-evenly">
        {#each getParties(focused) as party, i}
            <div class="grid grid-cols-2 gap-2 w-full {partyTextColors[i]}">
                {#each sortByDPS(focused, party) as player}
                    <div class="h-[88px] bg-[#F4EDE9] rounded-sm flex flex-col justify-center items-center">
                        <div class="self-start mx-auto mb-auto w-[50%] h-[1.5px]">
                            <div style="width: {getDamagePercent(player)}%"
                                 class="rounded-b-[0.1rem] mx-auto h-full {partyBgColors[i]}"></div>
                        </div>
                        <div class="flex flex-col items-center justify-evenly">
                            <p class="text-sm font-medium">{player}</p>
                            <p class="text-xs">{focused.players[player].class}</p>
                            <p class="font-medium">{formatDamage(focused.players[player].dps)}</p>
                        </div>
                        <div class="self-start mx-auto mb-auto w-[50%] h-[1.5px]">
                        </div>
                    </div>
                {/each}
            </div>
        {/each}
    </div>
</div>
<div class="mb-2 flex flex-row items-center justify-center">
    {#each {length: 3} as _, i}
        <button class="px-1" on:click={() => setTab(i)}>
            <div class="w-2 h-2 rounded-3xl shadow-sm border-[#a2596d] border-[0.5px] {tab === i ? 'bg-[#a2596d]' : 'bg-[#f7f2ef]'}"></div>
        </button>
    {/each}
</div>