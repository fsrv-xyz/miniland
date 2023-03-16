<template>
  <div class="hello">
    <h1>{{ msg }}</h1>
    <p>
      {{ load }}
    </p>
  </div>
</template>

<script>
export default {
  name: 'HardwarePanel',
  props: {
    msg: String,
  },
  data() {
    return {
      eventClient: null,
      load: "foo",
    };
  },
  mounted() {
  },
  created: function () {
    this.eventClient = new EventSource("http://100.64.70.93:8080/frontend/sse/load");
    this.eventClient.onmessage = (event) => {
      console.log(event.data);
      this.load = event.data;
    };
  }
}
</script>
