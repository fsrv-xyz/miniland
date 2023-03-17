<template>
  <img height="128" width="128" alt="Miniland Logo" src="./assets/logo.png">
  <h1>Miniland</h1>
  <SystemUsagePanel title="System usage"/>
</template>

<script>
import SystemUsagePanel from './components/SystemUsagePanel.vue'

export default {
  name: 'App',
  components: {
    SystemUsagePanel,
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
