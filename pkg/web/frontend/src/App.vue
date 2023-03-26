<template>
  <div class="navbar">
    <a :href=git_repo_link><img src="./assets/logo.png" alt="Logo" class="navbar-logo"></a>
    <h1>Miniland</h1>
    <a style="margin-left: auto" :href=build_info.pipeline_url>Build: {{ build_info.version_string }}</a>
  </div>
  <hr>
  <div class="panel">
    <SystemUsagePanel title="System usage"/>
  </div>
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
  get git_repo_link() {
    return process.env.VUE_APP_SRC_REPO_LINK || '#';
  }

  get build_info() {
    return {
      version_string: process.env.VUE_APP_BUILD_INFO || 'unknown',
      pipeline_url: process.env.VUE_APP_PIPELINE_URL || '#',
    };
  }
}
</script>

<style>
html {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  margin-top: 2vw;
  margin-left: 2vw;
  margin-right: 2vw;
  color: #2c3e50;
}

hr {
  border: 0;
  height: 1px;
  background: #eee;
}

@media only screen and (max-width: 600px){
  .navbar-logo {
    height: 50px;
    margin-right: 16px;
  }
}

@media only screen and (min-width: 600px){
  .navbar-logo {
    height: 100px;
    margin-right: 16px;
  }
}

.panel {
  padding-inline: 10%;
  border: 1px solid #ddd;
  border-radius: 5px;
  margin: 0 auto;
  display: flex;
}

.navbar {
  display: flex;
  align-items: center;
  padding: 8px;
}

.navbar a {
  color: slategray;
  text-decoration: none;
  padding: 8px;
}

</style>
