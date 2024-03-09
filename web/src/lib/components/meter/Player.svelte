<script lang="ts">
    import {getClassIcon} from "$lib/game";
    import IconSkull from "~icons/lucide/ghost";
    import {formatDuration} from "$lib/print";

    export let player;
    export let anonymized;
    export let difficulty;

    let hovered;
</script>

<img alt={player.class}
     src="{getClassIcon(player.class)}"
     class="h-6 w-6 mr-1.5 z-50 inline"
/>
<span class="z-50">{["Inferno", "Trial"].includes(difficulty) ? "" : Math.floor(player.gearScore) + " " }{anonymized ? player.class + " " + player.name : player.name}</span>
{#if player.dead}
    <div class="ml-0.5 flex items-center justify-center"
         on:mouseover={() => hovered = true}
         on:mouseleave={() => hovered = false}>
        <IconSkull/>
        {#if hovered}
            <div class="absolute flex flex-row items-center justify-center p-1 z-50 rounded-lg whitespace-nowrap bg-bouquet-50 border-[1px] border-[#c58597] -translate-y-[calc(100%-0.1rem)] text-[#575279]">
                <IconSkull/>
                <p class="ml-1">Dead for {formatDuration(player.deadFor)}</p>
            </div>
        {/if}
    </div>
{/if}