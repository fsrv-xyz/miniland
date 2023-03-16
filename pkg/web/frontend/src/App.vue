<template>
  <img height="128" width="128" alt="Miniland Logo" src="./assets/logo.png">
  <h1>Miniland</h1>
  <HardwarePanel msg="changeme"/>
</template>

<script>
import HardwarePanel from './components/HardwarePanel.vue'

export default {
  name: 'App',
  components: {
    HardwarePanel,
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
