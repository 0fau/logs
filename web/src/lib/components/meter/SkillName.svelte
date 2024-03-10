<script lang="ts">
    import { getSkillIcon } from "$lib/game";

    export let info;

</script>

<td class="relative z-10 pl-1 py-0.5">
    <img
        alt={info.name}
        src={getSkillIcon(info.icon)}
        class="inline h-6 w-6 opacity-95" />
</td>
<td colspan="2" class="relative z-10">
    <div class="truncate">
        {info.name}
    </div>
</td>
