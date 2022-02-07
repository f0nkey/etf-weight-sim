<script>
  import spyData from './assets/spy.json'

  // todo: add correlation color display on hover
  // todo: add cumulative weight
  // todo: add extra buttons for +0.1, +0.5, +1

  function addSimWeightProp(data) {
    data.map((security) => {
      if(!security.hasOwnProperty("change")) {
        security.change = 0;
      }
    })

    return data
  }

  $: etfData = addSimWeightProp(spyData)
  $: etfChange = sum(etfData)

  function useRecentChange() {
    etfData = etfData.map((a) => {
      a.change = a.daily_change
      return a;
    })
  }

  function sum(a) {
    let total = 0
    a.forEach((b) => {
      total += (b.change / 100) * b.weight
    })
    return Math.round((total + Number.EPSILON) * 100) / 100 // prevents floating point rounding error
  }
</script>

<main>
  <h1>ETF Weighting Simulator</h1>
  <p>Determine the impact one stock has on a whole fund.</p>
  <button on:click={useRecentChange}>Use Recent Gain/Loss Data</button>
  <div class="layout">
    <table>
      <tr><th>Name</th><th>Ticker</th><th>Weight</th><th>Gain/Loss</th></tr>
      {#each etfData as row}
        <tr>
          <td title={row.name}>{row.name}</td>
          <td title={row.name}>{row.ticker}</td>
          <td>{row.weight}</td>
          <td><input id="number" on:change={()=>{etfData = etfData}} step="0.1" type="number" bind:value={row.change}>%</td>
        </tr>
      {/each}
    </table>
    <div style="position: relative">
      <div class="etf-panel">
        <h2>Overall SPY Change:</h2>
        <h1>{etfChange}%</h1>
      </div>
    </div>
  </div>
</main>

<style>
  :root {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen,
      Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  }

  main {
    text-align: center;
    padding: 1em;
    margin: 0 auto;
  }

  .layout {
    display: grid;
    grid-template-columns: max-content 1fr;
  }

  .etf-panel {
    position:sticky;
    top:50%;
  }

  table {
    table-layout: fixed;
    width: 500px;
  }

  tr > td:nth-child(1), tr > th:nth-child(1) {
    text-align: left;
    width: 12em;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  td > input {
    width: 60px
  }
</style>
