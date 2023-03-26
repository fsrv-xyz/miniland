<template>
  <tr>
    <td><b>{{ data.path !== "" ? data.path : "n/a" }}</b></td>
    <td class="progress-column">
      <progress :title="(percent !== 0 ? percent : 'n/a') + '%'" :value="percent" class="progress" max="100"></progress>
    </td>
    <td>{{ data.used }} / {{ data.total }} MiB</td>
  </tr>
</template>

<script lang="ts">
import { defineComponent, computed } from 'vue';

interface Filesystem {
  path: string;
  total: number;
  used: number;
}

export default defineComponent({
  props: {
    data: {
      type: Object as () => Filesystem,
      required: true,
    },
  },
  setup(props) {
    const percent = computed(() => {
      if (!props.data.total || props.data.total === 0) {
        return 0;
      }
      return +(((props.data.used || 0) / props.data.total) * 100).toFixed(2);
    });

    return {
      percent,
    };
  },
});
</script>

<style scoped>
td {
  padding-left: 5px;
  padding-right: 5px;
}
@media screen and (max-width: 600px){
  .progress-column{
    display: none;
  }
}
</style>
