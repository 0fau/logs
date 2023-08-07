<script lang='ts'>
    import Preview from '../lib/Preview.svelte'
</script>

<div class="container flex h-screen justify-center items-center">
    <div class="columns-3xs ...">
      <h1 class="text-3xl text-center font-bold">
        logs.fau.dev
      </h1>
      <p class="text-center">analytics engine for lost ark parses (wip)</p>
      <Preview />
    </div>
</div>

<style lang="postcss">
  :global(html) {
    background-color: theme(colors.gray.100);
  }
</style>