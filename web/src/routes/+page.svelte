<script lang='ts'>
    import Donut from '../lib/Donut.svelte'
</script>

<div class="container flex h-screen justify-center items-center">
    <div class="columns-3xs ...">
      <h1 class="text-3xl text-center font-bold">
        logs.fau.dev
      </h1>
      <p class="text-center">analytics engine for lost ark parses (wip)</p>
      <Donut />
      <div class ="text-center">&gt;&gt; <a class="underline font-bold" href="/me">enter</a> &lt;&lt;</div>
      <br />
      <div class="text-center">built by <a href="https://fau.dev" class="text-sky-600 underline">@faustuwu</a> for friends <img src="https://cdn.discordapp.com/emojis/1090698237410803743.webp?size=240&quality=lossless" alt="<3" class="inline w-6 -translate-y-0.5" /></div>
    </div>
</div>

<style lang="postcss">
  :global(html) {
    background-color: theme(colors.gray.100);
  }
</style>