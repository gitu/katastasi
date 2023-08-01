<script setup lang="ts">
import type {Service} from "@/types";
import {OK} from "@/types";
import type {PropType} from "vue";
import StatusIcon from "@/components/StatusIcon.vue";


const props = defineProps({
  service: {
    type: Object as PropType<Service>,
    required: true
  },
  showDetails: {
    type: Boolean,
    required: false,
    default: false
  }
});
</script>

<template>
  <section class="section container box">
    <div class="title">{{ service.name }}</div>
    <div class="content">
      <status-icon :status="service.status"></status-icon>
    </div>
    <div class="content" v-if="showDetails || service.status!==OK">
      <div class="columns">
        <template v-for="component in service?.serviceComponents">
          <div class="column" v-if="showDetails || component.status!==OK">
            <div class="title is-5">{{ component.name }}</div>
            <div class="content">
              <status-icon :status="component.status"></status-icon>
            </div>
            <div class="notification" v-if="component.info">{{ component.info }}</div>
          </div>
        </template>
      </div>
    </div>
  </section>

</template>

<style scoped>

</style>