<script lang="ts">
    import * as d3 from 'd3';
    import {onMount} from "svelte";

    export let encounter;

    onMount(() => {
        const margin = {top: 15, right: 20, bottom: 25, left: 20},
            width = 600 - margin.left - margin.right,
            height = 280 - margin.top - margin.bottom;

        const hpLog = encounter.data.bossHPLog;
        let bosses = Object.keys(hpLog).filter((b) => hpLog[b].length > 1)
            .sort((a, b) => hpLog[a].time - hpLog[b].time)
        let end = d3.max(bosses, b => hpLog[b][hpLog[b].length - 1].time)

        const svg = d3.select("#graph")
            .append("svg")
            .attr("width", width + margin.left + margin.right)
            .attr("height", height + margin.top + margin.bottom)
            .append("g")
            .attr("transform", `translate(${margin.left},${margin.top})`);

        const x = d3.scaleTime()
            .domain([0, new Date(0).setSeconds(end)])
            .range([0, width]);
        const xAxis = svg.append("g")
            .attr("transform", `translate(0,${height + 4})`)
            .call(
                d3.axisBottom(x)
                    .ticks(d3.timeMinute)
                    .tickFormat(function (d) {
                        return d3.timeFormat("%-M")(d)
                    })
            )
            .attr("color", "#575279");
        const y = d3.scaleLinear()
            .domain([0, 1])
            .range([height, 0]);

        const bisector = d3.bisector(d => d.time).center;

        bosses.forEach((b) => {
            const boss = hpLog[b];

            const data = boss.map(log => {
                return {time: new Date(0).setSeconds(log.time), percent: log.p}
            })
            const areaGenerator = d3.area()
                .x(function (d) {
                    return x(d.time)
                })
                .y0(y(0))
                .y1(function (d) {
                    return y(d.percent)
                });

            const area = svg.append('g');
            area.append("path")
                .datum(data)
                .attr("class", "myArea")
                .attr("fill", "#eedede")
                .attr("d", areaGenerator)

            const lineGenerator = d3.line()
                .x(d => x(d.time))
                .y(d => y(d.percent));

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
    <div id="graph"></div>
</div>