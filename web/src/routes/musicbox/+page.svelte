<script lang="ts">
    import Sidebar from '$lib/components/Sidebar.svelte'
</script>

<svelte:head>
    <title>logs by faust</title>
</svelte:head>

<div class="w-screen h-screen flex justify-center items-center">
    <div class="w-[80%] h-[80%] rounded-xl flex flex-row shadow-lg">
        <div class="w-[20%] h-full bg-[#b4637a] rounded-l flex flex-col">
            <Sidebar selected="musicbox"/>
        </div>
        <div class="w-[35%] h-full bg-[#c68a9b]"></div>
        <div class="w-[45%] h-full bg-[#fcf9f6] rounded-r"></div>
    </div>
</div>

<style lang="postcss">
    :global(body) {
        background-color: #faf4ed;
    }
</style>