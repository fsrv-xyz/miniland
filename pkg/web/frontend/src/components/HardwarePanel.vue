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
      load: "0",
    };
  },
  created: function () {
    this.eventClient = new EventSource("/frontend/sse/load");
    this.eventClient.onmessage = (event) => {
      this.load = JSON.parse(event.data).message;
    };
  }
}
</script>
