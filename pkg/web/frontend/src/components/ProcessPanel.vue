<template>
  <button @click="startProcess">Start Process</button>
  <div class="hello">
    <h1>{{ msg }}</h1>
  </div>
</template>

<script>
export default {
  name: 'ProcessPanel',
  props: {
    msg: String,
  },
  data() {
    return {
      processes: [],
    };
  },
  created() {
    this.pullProcesses();
  },
  methods: {
    pullProcesses: function () {
      fetch("http://100.64.70.93:8080/api/processes")
        .then(response => response.json())
        .then(data => {
          this.processes = data;
        });
    },
  },
}
</script>
