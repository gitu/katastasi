<script setup lang="ts">
import type {StatusInfo} from "@/types";
import type {PropType} from "vue";
import StatusIcon from "@/components/StatusIcon.vue";

const props = defineProps({
  status: {
    type: Object as PropType<StatusInfo>,
    required: true
  }
})

</script>

<template>
  <div class="container">
    <div class="title">{{ props.status.name }}</div>
    <div class="subtitle"><status-icon :status="props.status.overallStatus"></status-icon></div>
    <div class="subtitle">{{ $filters.timeAgo(props.status.lastUpdate) }}</div>
  </div>

  <div class="container" v-for="service in props.status?.services">
    <div class="title">{{ service.name }}</div>
    <div class="subtitle"><status-icon :status="service.status"></status-icon></div>
    <div class="subtitle">{{ $filters.timeAgo(service.lastUpdate) }}</div>

    <div class="container" v-for="component in service?.serviceComponents">
      <div class="title">{{ component.name }}</div>
      <div class="subtitle"><status-icon :status="component.status"></status-icon></div>
      <div class="subtitle">{{ component.info }}</div>
    </div>
  </div>
</template>

<style scoped>

</style>