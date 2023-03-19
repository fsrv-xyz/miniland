<template>
  <tr>
    <td><b>{{ filesystem.path }}</b></td>
    <td><progress class="progress" :value="filesystem.percent" max="100" :title="filesystem.percent+'%'">{{
        filesystem.percent
      }}%</progress></td>
    <td>{{ filesystem.used }} / {{ filesystem.total }} MiB</td>
  </tr>
</template>

<script>

export default {
  name: 'SystemUsagePanel',
  props: {
    data: {},
  },
  data() {
    return {
      filesystem: {
        "path": "n/a",
        "total": "n/a",
        "used": "n/a",
        "percent": 0,
      },
    };
  },
  created: function () {
    ({Path: this.filesystem.path, Used: this.filesystem.used, Total: this.filesystem.total} = this.data);
    this.filesystem.percent = (this.filesystem.used / this.filesystem.total * 100).toFixed(2);
  }
}
</script>

<style scoped>
td {
  padding-left: 5px;
  padding-right: 5px;
}
</style>