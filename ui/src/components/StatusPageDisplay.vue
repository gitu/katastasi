<script setup lang="ts">
import type {StatusInfo} from "@/types";
import {OK} from "@/types";
import type {PropType} from "vue";
import {computed, ref} from "vue";
import StatusIcon from "@/components/StatusIcon.vue";
import StatusServiceDisplay from "@/components/StatusServiceDisplay.vue";
import ArrowExpand from "vue-material-design-icons/ArrowExpand.vue";

const props = defineProps({
  status: {
    type: Object as PropType<StatusInfo>,
    required: true
  }
});

const showDetails = ref(false);

const hasOkServices = computed(() => {
  return props.status?.services.some(service => service.status === OK);
});

</script>

<template>
  <section class="columns  p-6">
    <div class="column">
      <div class="container block">
        <div class="title">{{ props.status.name }}</div>
        <div class="content">
          <status-icon :status="props.status.overallStatus"></status-icon>
        </div>
      </div>
    </div>
    <div class="column is-2 has-text-right">
      <div class="content" v-if="hasOkServices && !showDetails">
        <button class="button is-small is-inverted is-info"
                @click="showDetails=true">
      <span class="icon is-small">
        <ArrowExpand></ArrowExpand></span>
          <span>show details</span>
        </button>
      </div>
    </div>

  </section>

  <template v-for="service in props.status?.services">
    <status-service-display :service="service"
                            :showDetails="showDetails"
                            v-if="showDetails || service.status!==OK"></status-service-display>
  </template>

</template>

<style scoped>

</style>