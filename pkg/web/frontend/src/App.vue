<template>
  <img alt="Miniland Logo" height="128" src="./assets/logo.png" width="128">
  <h1>Miniland</h1>
  <SystemUsagePanel title="System usage"/>
</template>

<script lang="ts">
import {Options, Vue} from 'vue-class-component';
import SystemUsagePanel from "@/components/SystemUsagePanel.vue";

@Options({
  components: {
    SystemUsagePanel,
  },
})
export default class App extends Vue {
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
