<template>
  <div>
    <h2>{{ title }}</h2>
    <table class="root_table">
      <tr>
        <td class="head">Load:</td>
        <td class="content">{{ load }}</td>
      </tr>
      <tr>
        <td class="head">Memory:</td>
        <td class="content">{{ memory }} MiB</td>
      </tr>
      <tr>
        <td class="head">Filesystems:</td>
        <td class="content">
          <table>
            <SingleFilesystemUsage v-for="disk in filesystems" :key="disk.path" :data="disk"/>
          </table>
        </td>
      </tr>
    </table>
  </div>

</template>

<script>
import SingleFilesystemUsage from './SingleFilesystemUsage.vue'

export default {
  name: 'SystemUsagePanel',
  components: {
    SingleFilesystemUsage,
  },
  props: {
    title: String,
  },
  data() {
    return {
      eventClient: null,
      load: "n/a",
      memory: "n/a",
      filesystems: [],
    };
  },
  created: function () {
    this.eventClient = new EventSource("/frontend/sse/usage");
    //this.eventClient = new EventSource("https://miniland.wwwtest.org/frontend/sse/usage");
    this.eventClient.onmessage = (event) => {
      const payload = JSON.parse(event.data).message;
      this.load = payload.loadavg;
      this.memory = payload.memused;
      this.filesystems = payload.filesystems.filter(disk => disk.total !== 0 && disk.path !== "/dev");
    };
  }
}
</script>

<style scoped>
.root_table {
  border-collapse: collapse;
  width: 100%;
  position: center;
}
.head {
  width: 50%;
  text-align: right;
  font-weight: bold;
}
.content {
  width: 50%;
  text-align: left;
  padding: 8px;
}
</style>
