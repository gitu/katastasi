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
  <div class="container block">
    <div class="title">{{ props.status.name }}</div>
    <div class="content"><status-icon :status="props.status.overallStatus"></status-icon></div>
  </div>



  <div class="container box" v-for="service in props.status?.services">
    <div class="title">{{ service.name }}</div>
    <div class="content"><status-icon :status="service.status"></status-icon></div>

    <div class="columns">
      <div class="column" v-for="component in service?.serviceComponents">
        <div class="title is-5">{{ component.name }}</div>
        <div class="content"><status-icon :status="component.status"></status-icon></div>
        <div class="notification" v-if="component.info">{{component.info}}</div>
      </div>
      </div>
  </div>
</template>

<style scoped>

</style>