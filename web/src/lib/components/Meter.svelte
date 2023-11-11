<script lang="ts">
    import {getRaid} from "$lib/raids";
    import DamageView from "$lib/components/DamageView.svelte";
    import PlayerDamageView from "$lib/components/PlayerDamageView.svelte";
    import BuffView from "$lib/components/BuffView.svelte";
    import {writable} from "svelte/store";
    import PlayerBuffView from "$lib/components/PlayerBuffView.svelte";

    export let encounter;
    console.log(encounter)

    let raid = getRaid(encounter.boss)
    console.log(raid)

    enum MeterTab {
        Damage,
        Buff,
    }

    let tab = MeterTab.Buff;
    let focus = writable("");
</script>

<div class="mt-5 rounded-xl min-w-[600px] border-[1px] border-[#efdcc5] shadow-sm bg-[#E9D4DA] p-5">
    {#if tab === MeterTab.Damage}
        {#if $focus === ""}
            <DamageView {encounter} {focus}/>
        {:else}
            <PlayerDamageView {encounter} {focus}/>
        {/if}
    {:else if tab === MeterTab.Buff}
        {#if $focus === ""}
            <BuffView {encounter} {focus}/>
        {:else}
            <PlayerBuffView {encounter} {focus}/>
        {/if}
    {/if}
</div>