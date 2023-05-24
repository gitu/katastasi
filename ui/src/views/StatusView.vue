<script setup lang="ts">

import type {StatusInfo} from "@/types";
import {ref, watch} from "vue";
import {useRoute} from "vue-router";
import StatusPageDisplay from "@/components/StatusPageDisplay.vue";

const loading = ref(true);
const status = ref({} as StatusInfo);
const route = useRoute();
const error = ref(false);

async function fetchStatus(newId: string) {
  const response = await fetch(`/api/status/${newId}`);
  if (!response.ok) {
    error.value = true;
    loading.value = false;
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  status.value = await response.json();
}

watch(() => route.params.id, async (newId) => {
  loading.value = true;
  status.value = {} as StatusInfo;
  error.value = false;
  await fetchStatus(newId as string);
  loading.value = false;
}, {immediate: true});
</script>

<template>
  <div class="container" v-if="loading || error">
    <div v-if="loading" class="title">Loading...</div>
    <div v-if="error" class="title">Error loading status</div>
    <span class="subtitle" v-if="error">
      Error loading status: {{ route.params.id }}
    </span>
  </div>

  <status-page-display v-if="!loading && !error" :status="status"></status-page-display>
</template>


