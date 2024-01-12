<script lang="ts">
    import {getRaid} from "$lib/raids";
    import DamageView from "$lib/components/meter/DamageView.svelte";
    import PlayerDamageView from "$lib/components/meter/PlayerDamageView.svelte";
    import BuffView from "$lib/components/meter/BuffView.svelte";
    import {writable} from "svelte/store";
    import PlayerBuffView from "$lib/components/meter/PlayerBuffView.svelte";
    import SelfBuffView from "$lib/components/meter/SelfBuffView.svelte";
    import PlayerSelfBuffView from "$lib/components/meter/PlayerSelfBuffView.svelte";
    import DamageGraph from "$lib/components/meter/DamageGraph.svelte";

    export let encounter;
    console.log(encounter)

    let raid = getRaid(encounter.boss)
    console.log(raid)

    enum MeterTab {
        Damage = "Damage",
        Buff = "Buff",
        Self = "Self",
    }

    let tabs = [MeterTab.Damage, MeterTab.Buff, MeterTab.Self]

    let current = MeterTab.Damage;
    let focus = writable("");

    function setTab(tab) {
        current = tab
    }
</script>

<div class="mt-2.5">
    {#each tabs as tab}
        <button class="p-1 px-1.5 m-0.5 mb-1 font-medium w-[80px] border-[1px] text-sm border-[#c58597] rounded-lg bg-[#F4EDE9]"
                class:bg-[#c58597]={current === tab}
                class:text-[#c58597]={current !== tab}
                class:text-[#F4EDE9]={current === tab}
                on:click={() => setTab(tab)}>
            {tab}
        </button>
    {/each}
</div>
{#if current === MeterTab.Damage}
    {#if $focus === ""}
        <div class="rounded-xl min-w-[600px] border-[1px] border-[#f7efe5] shadow-sm bg-[#f5ece8] p-2">
            <DamageView {encounter} {focus}/>
        </div>
        {#if encounter.data.bossHPLog}
            <div class="rounded-xl mt-5 min-w-[600px] border-[1px] border-[#f7efe5] shadow-sm bg-[#f5ece8] p-2">
                <DamageGraph {encounter}/>
            </div>
        {/if}
    {:else}
        <div class="rounded-xl min-w-[600px] border-[1px] border-[#f7efe5] shadow-sm bg-[#f5ece8] p-2">
            <PlayerDamageView {encounter} {focus}/>
        </div>
    {/if}
{:else if current === MeterTab.Buff}
    {#if $focus === ""}
        <div class="rounded-xl min-w-[600px] border-[1px] border-[#f7efe5] shadow-sm bg-[#f5ece8] p-2">
            <BuffView {encounter} {focus}/>
        </div>
    {:else}
        <div class="rounded-xl min-w-[600px] border-[1px] border-[#f7efe5] shadow-sm bg-[#f5ece8] p-2">
            <PlayerBuffView {encounter} {focus}/>
        </div>
    {/if}
{:else if current === MeterTab.Self}
    {#if $focus === ""}
        <div class="rounded-xl min-w-[600px] border-[1px] border-[#f7efe5] shadow-sm bg-[#f5ece8] p-2">
            <SelfBuffView {encounter} {focus}/>
        </div>
    {:else}
        <div class="rounded-xl min-w-[600px] border-[1px] border-[#f7efe5] shadow-sm bg-[#f5ece8] p-2">
            <PlayerSelfBuffView {encounter} {focus}/>
        </div>
    {/if}
{/if}