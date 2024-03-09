<script>
  import "@fontsource-variable/inter";
  import "../app.css";
  import {settings} from "$lib/store";
</script>

<div class="{$settings?.darkMode ? 'dark' : ''}">
    <div class="bg-tapestry-100 dark:bg-zinc-700 h-screen w-screen">
        <slot />
    </div>
</div>
