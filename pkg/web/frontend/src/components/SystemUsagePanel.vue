<template>
  <div>
    <h2>{{ title }}</h2>
    <table>
      <tr>
        <td class="head">Load:</td>
        <td class="content">{{ load }}</td>
      </tr>
      <tr>
        <td class="head">Memory:</td>
        <td class="content">{{ memory }} MiB</td>
      </tr>
      <tr>
        <td class="head">Filesystem:</td>
        <td class="content">
          <SingleFilesystemUsage v-for="disk in disks" :key="disk.path" :data="disk"/>
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
      disks: [],
    };
  },
  created: function () {
    this.eventClient = new EventSource("https://miniland.wwwtest.org/frontend/sse/usage");
    this.eventClient.onmessage = (event) => {
      const payload = JSON.parse(event.data).message;
      this.load = payload.loadavg;
      this.memory = payload.memused;
      this.disks = payload.disks.filter(disk => disk.Total !== 0 && disk.Path !== "/dev");
    };
  }
}
</script>

<style scoped>
table {
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
progress {
  margin-left: 2%;
}
</style>
