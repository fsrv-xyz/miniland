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
      <tr v-if="filteredFilesystems.length !== 0">
        <td class="head">Filesystems:</td>
        <td class="content">
          <table>
            <SingleFilesystemUsage
                v-for="filesystem in filteredFilesystems"
                :key="filesystem.path"
                :data="filesystem"
            />
          </table>
        </td>
      </tr>
    </table>
  </div>
</template>

<script lang="ts">
import {Options, Vue} from 'vue-class-component';
import SingleFilesystemUsage from './SingleFilesystemUsage.vue';

interface Filesystem {
  path: string;
  total: number;
  used: number;
}

interface EventPayload {
  loadavg: number;
  memused: number;
  filesystems: Filesystem[];
}

@Options({
  components: {
    SingleFilesystemUsage,
  },
  props: {
    title: String,
  },
})
export default class SystemUsagePanel extends Vue {
  eventClient: EventSource | null = null;
  load: number | string = 'n/a';
  memory: number | string = 'n/a';
  filesystems: Filesystem[] = [];

  get filteredFilesystems() {
    if (this.filesystems.length === 0) {
      return [];
    }
    return this.filesystems.filter(
        (filesystem: Filesystem) => filesystem.total !== 0 && filesystem.path !== '/dev'
    );
  }

  created() {
    //this.eventClient = new EventSource('https://miniland.wwwtest.org/frontend/sse/usage');
    this.eventClient = new EventSource('/frontend/sse/usage');
    this.eventClient.onmessage = this.handleMessage;
  }

  handleMessage(event: MessageEvent) {
    const payload : EventPayload = JSON.parse(event.data).message;
    this.load = payload.loadavg;
    this.memory = payload.memused;
    this.filesystems = payload.filesystems;
  }
}
</script>

<style scoped>
.root_table {
  border-collapse: collapse;
  width: 100%;
  margin: 0 auto;
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
