<script lang="ts">
    import * as d3 from 'd3';
    import {onMount} from "svelte";
    import {bossHpMap, getRaid, isGuardian} from "$lib/raids";
    import {formatDamage, formatDuration, formatPercent} from "$lib/print";
    import {shortBossName} from "$lib/raids.js";

    export let encounter;

    const margin = {top: 15, right: 20, bottom: 25, left: 20},
        width = 600,
        height = 280;

    const allowed = ["Nightmare Gehenna", "Nightmare Helkasirs"];
    const disallowed = ["Brand of Subordination"];

    let dataMap = {};
    const hpLog = encounter.data.bossHPLog;
    let bosses = Object.keys(hpLog).filter((b) => hpLog[b].length > 1)
        .filter((b) => !disallowed.includes(b) && (getRaid(b) || isGuardian(b) || allowed.includes(b)))
        .sort((a, b) => hpLog[a].time - hpLog[b].time)
    let end = d3.max(bosses, b => hpLog[b][hpLog[b].length - 1].time)

    const x = d3.scaleTime()
        .domain([0, new Date(0).setSeconds(end)])
        .range([margin.left, width - margin.right]);

    const y = d3.scaleLinear()
        .domain([0, 1])
        .range([height - margin.bottom, margin.top]);

    let tooltip;
    let showTooltip = false;
    let tooltipBosses = [];
    let tooltipPlayers = [];
    let currentX;
    let tooltipX = 0;
    let tooltipY = height;

    function pointermoved(event) {
        tooltipY = height;
        tooltipBosses = [];
        tooltipPlayers = [];

        currentX = x.invert(d3.pointer(event)[0]);

        const bisector = d3.bisector(d => d.time).center;
        for (let [boss, hpLog] of Object.entries(dataMap)) {
            let index = bisector(hpLog, currentX);
            if (index <= 0 || index >= hpLog.length - 1) {
                continue;
            }

            tooltipBosses.push({boss: boss, hp: hpLog[index].hp, percent: hpLog[index].percent})

            tooltipX = event.offsetX;
            tooltipY = Math.min(tooltipY, y(hpLog[index].percent));
        }

        for (let [name, player] of Object.entries(encounter.data.players)) {
            tooltipPlayers.push({name: name, dps: player.dpsLog[Math.floor(currentX / 5000)]})
        }
        tooltipPlayers.sort((a, b) => b.dps - a.dps);

        tooltipX -= tooltip.offsetWidth / 2;
        tooltipY = Math.max(tooltipY, 0)
        tooltipY -= (tooltip.offsetHeight + 50);

        if (tooltipBosses.length > 0) {
            showTooltip = true;
        } else {
            showTooltip = false;
        }
    }

    function pointerleft(event) {
        showTooltip = false;
        // tooltip.style.display = "none";
    }

    onMount(() => {
        const svg = d3.select("#graph")
            .append("svg")
            .attr("viewBox", [0, 0, width, height])
            .attr("width", width)
            .attr("height", height)
            .style("overflow", "visible")
            .attr("style", "max-width: 100%; height: auto; height: intrinsic;")
            .on("pointerenter pointermove", pointermoved)
            .on("pointerleave", pointerleft)
            .append("g")
            .attr("transform", `translate(${margin.left},${margin.top})`);

        const xAxis = svg.append("g")
            .attr("transform", `translate(-${margin.left},${height - margin.bottom - margin.top + 4})`)
            .call(
                d3.axisBottom(x)
                    .ticks(d3.timeMinute)
                    .tickFormat(function (d) {
                        return d3.timeFormat("%-M")(d)
                    })
            )
            .attr("color", "#575279");

        bosses.forEach((b) => {
            const boss = hpLog[b];
            const data = boss.map(log => {
                return {time: new Date(0).setSeconds(log.time), hp: log.hp, percent: log.p}
            })
            dataMap[b] = data;

            const areaGenerator = d3.area()
                .x(function (d) {
                    return x(d.time) - margin.left;
                })
                .y0(y(0) - margin.top)
                .y1(function (d) {
                    return y(d.percent) - margin.top;
                });

            const area = svg.append('g');
            area.append("path")
                .datum(data)
                .attr("class", "myArea")
                .attr("fill", "#eedede")
                .attr("d", areaGenerator)

            const lineGenerator = d3.line()
                .x(d => x(d.time) - margin.left)
                .y(d => y(d.percent) - margin.top);

            const line = svg.append('g');
            line.append("path")
                .datum(data)
                .attr('stroke', '#b96d83')
                .attr('fill', 'none')
                .attr("d", lineGenerator);
        })
    })
</script>

<div class="w-full h-full flex items-center justify-center">
    <div class="relative">
        <div class="absolute pointer-events-none z-[100] flex flex-col items-center justify-center"
             bind:this={tooltip}
             class:block={showTooltip}
             class:hidden={!showTooltip}
             style="transform: translate({tooltipX}px, {tooltipY}px)">
            <div class="flex flex-col items-center justify-center p-2 rounded-md whitespace-nowrap bg-[#F4EDE9] border-[1px] text-sm font-semibold border-[#c58597] text-[#575279]">
                {#each tooltipBosses as boss}
                    <p>{shortBossName(boss.boss)}
                        [{ bossHpMap[boss.boss] ? Math.floor(boss.percent * bossHpMap[boss.boss]) + "x" : formatDamage(boss.hp)}
                        - {boss.percent < 0.01 ? "0" : formatPercent(boss.percent)}%]</p>
                {/each}
                <div class="w-full my-1 h-[1px] bg-[#b96d83] opacity-90 rounded"></div>
                {#each tooltipPlayers as player}
                    <p>{encounter.anonymized ? encounter.players[player.name].class + " " + player.name : player.name}
                        - {formatDamage(player.dps)}</p>
                {/each}
            </div>
            <div class="flex flex-col mt-0.5 w-fit items-center justify-center p-1 rounded-md whitespace-nowrap bg-[#F4EDE9] border-[1px] text-xs font-semibold border-[#c58597] text-[#575279]">
                <p>{tooltipBosses.length > 0 && formatDuration(currentX)}</p>
            </div>
        </div>
        <div id="graph">
        </div>
    </div>
</div>