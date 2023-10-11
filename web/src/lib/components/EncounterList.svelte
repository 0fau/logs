<script lang="ts">
    import IconRefresh from '~icons/ri/refresh-line';

    export let encounters;
</script>

<div class="flex flex-row m-auto mt-10 justify-center items-center flex flex-col">
    {#each encounters as encounter}
        <div class="mb-6 h-[70px] rounded-lg w-4/5 bg-[#faf4ed] opacity-80">
        </div>
    {/each}
</div>