<script setup lang="ts">
import type {StatusInfo} from "@/types";
import type {PropType} from "vue";

const props = defineProps({
  status: {
    type: Object as PropType<StatusInfo>,
    required: true
  }
})
</script>

<template>
  <div class="container">
    <div class="title">{{ props.status.statusPage.name }}</div>
    <div class="subtitle">{{ props.status.overallStatus }}</div>
    <div class="subtitle">{{ props.status.lastUpdate }}</div>
  </div>

  <div class="container" v-for="service in props.status?.services">
    <div class="title">{{ service.name }}</div>
    <div class="subtitle">{{ service.status }}</div>
    <div class="subtitle">{{ service.lastUpdate }}</div>
  </div>
</template>

<style scoped>

</style>