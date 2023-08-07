<script lang="ts">
  import * as d3 from 'd3';

  let svg;
  let width = 500;
  let height = 350;
  let margin = 20;

  const radius = Math.min(width, height) / 2 - margin

  let data = {arcana: 9, sharpshooter: 20, gunlancer:30, paladin:8, deathblade:12, bard:3, wardancer:7, gunslinger:14}

  const color = d3.scaleOrdinal()
    .domain(["a", "b", "c", "d", "e", "f", "g", "h"])
    .range(d3.schemeDark2);

  const pie = d3.pie()
    .sort(null)
    .value(d => d[1])
  const data_ready = pie(Object.entries(data))

  const arc = d3.arc()
    .innerRadius(radius * 0.5)
    .outerRadius(radius * 0.8)

  const outerArc = d3.arc()
    .innerRadius(radius * 0.9)
    .outerRadius(radius * 0.9)

  function polylinePoints(d) {
    var posA = arc.centroid(d) // line insertion in the slice
    var posB = outerArc.centroid(d) // line break: we use the other arc generator that has been built only for that
    var posC = outerArc.centroid(d); // Label position = almost the same as posB
    var midangle = d.startAngle + (d.endAngle - d.startAngle) / 2 // we need the angle to see if the X position will be at the extreme right or extreme left
    posC[0] = radius * 0.95 * (midangle < Math.PI ? 1 : -1); // multiply by 1 or -1 to put it on the right or on the left
    return [posA, posB, posC]
  }

  function labelTransform(d) {
    var pos = outerArc.centroid(d);
    var midangle = d.startAngle + (d.endAngle - d.startAngle) / 2
    pos[0] = radius * 0.99 * (midangle < Math.PI ? 1 : -1);
    return 'translate(' + pos + ')';
  }

  function labelTextAnchor(d) {
    var midangle = d.startAngle + (d.endAngle - d.startAngle) / 2
    return (midangle < Math.PI ? 'start' : 'end')
  }
</script>

<svg bind:this={svg} {width} {height}>
  <g transform="translate({width / 2}, {height / 2})">
    {#each data_ready as d}
      <path d={arc(d)} fill={color(d.data[1])} stroke=white stroke-width=2px opacity=0.7 />
      <polyline stroke={color(d.data[1])} fill=none stroke-width=1 points={polylinePoints(d)} />
      <text transform={labelTransform(d)} text-anchor={labelTextAnchor(d)}>{d.data[0]}</text>
    {/each}
  </g>
</svg>
