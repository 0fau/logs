<script lang="ts">
    import { getClassIcon } from "$lib/game";
    import { formatDuration } from "$lib/print";
    import IconSkull from "~icons/lucide/ghost";

    export let player;
    export let difficulty;
    export let focus;

    let hovered;
</script>

<td on:click={() => focus.set(player.name)} class="relative z-10 py-0.5 pl-1">
    <img alt={player.class} src={getClassIcon(player.class)} class="z-50 inline h-6 w-6" />
</td>
<td colspan="2" on:click={() => focus.set(player.name)} class="z-10">
    <div class="flex justify-start">
        <span class="z-10 truncate"
            >{["Inferno", "Trial"].includes(difficulty)
                ? ""
                : Math.floor(player.gearScore) + " "}{player.name}</span>
        {#if player.dead}
            <div
                class="z-10 ml-0.5 flex items-center justify-center"
                on:mouseover={() => (hovered = true)}
                on:mouseleave={() => (hovered = false)}>
                <IconSkull />
                {#if hovered}
                    <div
                        class="absolute z-50 flex -translate-y-[calc(100%-0.1rem)] flex-row items-center justify-center whitespace-nowrap rounded-lg border border-tapestry-300 bg-bouquet-50 p-1 text-gray-700">
                        <IconSkull />
                        <p class="ml-1">Dead for {formatDuration(player.deadFor)}</p>
                    </div>
                {/if}
            </div>
        {/if}
    </div>
</td>
