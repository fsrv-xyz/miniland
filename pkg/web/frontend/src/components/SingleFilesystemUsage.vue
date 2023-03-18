<template>
  <div>
    {{ disk.path }}
    {{ disk.used }}
    {{ disk.total }}
  </div>
</template>

<script>

export default {
  name: 'SystemUsagePanel',
  props: {
    data: {},
  },
  data() {
    return {
      disk: {
        "path": "n/a",
        "total": "n/a",
        "used": "n/a",
        "percent": "n/a",
      },
    };
  },
  created: function () {
    this.disk.path = this.data.Path;
    this.disk.used = this.data.Used;
    this.disk.total = this.data.Total;
    this.disk.percent = this.disk.Used / this.disk.Total * 100;
  }
}
</script>