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
        <td class="content">{{ memory }}</td>
      </tr>
    </table>
  </div>
</template>

<script>
export default {
  name: 'SystemUsagePanel',
  props: {
    title: String,
  },
  data() {
    return {
      eventClient: null,
      load: "n/a",
      memory: "n/a",
    };
  },
  created: function () {
    this.eventClient = new EventSource("/frontend/sse/usage");
    this.eventClient.onmessage = (event) => {
      const payload = JSON.parse(event.data).message;
      this.load = payload.loadavg;
      this.memory = payload.memused;
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
</style>
