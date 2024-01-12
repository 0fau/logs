<script>
    import {formatDamage, formatPercentFlat} from "$lib/print.ts";

    let partyTextColors = ["text-[#a2596d]", "text-[#5e5277]"]
    let partyBgColors = ["bg-[#a2596d]", "bg-[#5e5277]"]

    export let focused;
    let tab = 0;

    function setTab(num) {
        tab = num
    }

    let parties = [];
    focused.parties.forEach(p => {
        let party = [];
        p.forEach((name) => {
            if (focused.players[name]) {
                party.push(name);
            }
        })

        if (party.length > 0) {
            parties.push(party);
        }
    })

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

<div class="w-[94%] overflow-hidden">
    {#if tab === 0}
        <div class="h-[400px] px-2 mb-2 flex justify-center items-center border-[0.5px] border-[#c58597] shadow-sm rounded-md w-full bg-[#f7f2ef]">
            <div class="w-full h-full flex flex-col items-center justify-evenly">
                {#each parties as party, i}
                    <div class="grid grid-cols-2 gap-2 w-full {partyTextColors[i]}">
                        {#each sortByDPS(focused, party) as player}
                            {#if focused.players[player]}
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
                            {/if}
                        {/each}
                    </div>
                {/each}
            </div>
        </div>
    {:else if tab === 1}
        <div class="h-[400px] px-2 mb-2 flex justify-center items-center border-[0.5px] border-[#c58597] shadow-sm rounded-md w-full bg-[#f7f2ef]">
            <div class="w-full h-full flex flex-col items-center justify-center">
                <p class="text-[#a2596d] text-sm font-semibold mb-3">*meow*</p>
                <img alt="meow"
                     class="h-16"
                     src="https://cdn.discordapp.com/emojis/667829524486160424.gif?size=240&quality=lossless"/>
            </div>
        </div>
    {:else if tab === 2}
        <div class="h-[400px] px-2 mb-2 flex justify-center items-center border-[0.5px] border-[#c58597] shadow-sm rounded-md w-full bg-[#f7f2ef]">
            <div class="w-full h-full flex flex-col items-center justify-center">
                <p class="text-[#a2596d] text-sm font-semibold mb-3">*mwah*</p>
                <img alt="mwah"
                     class="h-16"
                     src="https://cdn.discordapp.com/emojis/749356125475962940.gif?size=240&quality=lossless"/>
            </div>
        </div>
    {/if}
</div>
<div class="mb-2 flex flex-row items-center justify-center">
    {#each {length: 3} as _, i}
        <button class="px-1" on:click={() => setTab(i)}>
            <div class="w-2.5 h-2 rounded-3xl shadow-sm border-[#a2596d] border-[0.5px] {tab === i ? 'bg-[#b96d83]' : 'bg-[#f7f2ef]'}"></div>
        </button>
    {/each}
</div>