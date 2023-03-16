<template>
  <div class="hello">
    <h1>{{ msg }}</h1>
  </div>
</template>

<script>
export default {
  name: 'HardwarePanel',
  props: {
    msg: String
  },
  mounted() {
    const source = new EventSource("http://100.64.70.93:8080/sse");
    source.onmessage = (event) => {
      console.log(event.data);
    };
  }
}
</script>
