<script lang="ts">
    import { getBossIcon } from "$lib/raids";
    import { formatDamage, formatDate, formatDateSolid, formatDuration } from "$lib/print";

    const difficultyColors = {
        Inferno: "#9a3148",
        Trial: "#9a3148",
        Challenge: "#625f77",
        Hard: "#b9982e",
        Normal: "#625f77"
    };

    function shortBossName(boss: string): string {
        if (boss === "Phantom Legion Commander Brelshaza") {
            return "Phantom LC Brelshaza";
        }

        if (boss === "Covetous Legion Commander Vykas") {
            return "Covetous LC Vykas";
        }

        if (boss === "Plague Legion Commander Akkan") {
            return "Plague LC Akkan";
        }

        if (boss === "Lazaram, the Trailblazer") {
            return "Lazaram";
        }

        if (boss === "Firehorn, Trampler of Earth") {
            return "Firehorn";
        }

        if (boss === "Rakathus, the Lurking Arrogance") {
            return "Rakathus";
        }

        if (boss === "Kaltaya, the Blooming Chaos") {
            return "Kaltaya";
        }

        return boss;
    }

    export let encounter;
    export let screenshot = false;

    $: player = encounter.players[encounter.localPlayer];
</script>

<div
    class="flex h-[80px] rounded-md border border-tapestry-400 bg-tapestry-50 shadow-sm"
    class:screenshot>
    <div class="flex h-full w-full flex-row items-center pl-3 pr-2">
        <div>
            <div class="self-start text-left text-gray-600">
                <div class="flex items-center space-x-1">
                    <img alt={encounter.boss} src={getBossIcon(encounter.boss)} class="size-6" />
                    <span class="truncate font-medium tracking-tight">{shortBossName(encounter.boss)}</span>
                </div>
                <p class="text-sm">
                    <span class="font-medium">{formatDamage(encounter.damage)}</span> damage dealt in {formatDuration(
                        encounter.duration
                    )}
                </p>
                <p class="text-xs tracking-tight text-gray-500">
                    <span class="font-medium">#{encounter.id} |</span>
                    <span
                        class="font-semibold"
                        style="color: {difficultyColors[encounter.difficulty]}"
                        >{encounter.difficulty}</span>
                    {screenshot
                        ? formatDateSolid(encounter.date)
                        : "cleared " + formatDate(encounter.date + encounter.duration)}
                </p>
            </div>
        </div>
        <div class="ml-auto flex h-full flex-col self-end rounded-r-md py-1 text-white">
            <span
                class="mr-0.5 mt-1.5 self-end rounded-sm bg-tapestry-500 px-1 py-0.5 text-center text-xs font-medium"
                >{encounter.anonymized
                    ? player.class + " " + encounter.localPlayer
                    : encounter.localPlayer}</span>
            <span class="mr-1 mt-0.5 self-end text-right text-xs font-medium text-tapestry-500"
                >{player.class}</span>
            <span class="my-auto mr-1 text-right text-lg font-medium text-gray-600"
                >{formatDamage(player.dps)}</span>
        </div>
    </div>
</div>
