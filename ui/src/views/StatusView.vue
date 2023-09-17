<script setup lang="ts">

import type {StatusInfo} from "@/types";
import {ref, watch} from "vue";
import {useRoute} from "vue-router";
import StatusPageDisplay from "@/components/StatusPageDisplay.vue";
import RefreshIcon from "@/components/RefreshIcon.vue";

const loading = ref(true);
const status = ref({} as StatusInfo);
const route = useRoute();
const error = ref(false);
const refreshIcon = ref<InstanceType<typeof RefreshIcon> | null>(null);


async function fetchStatus() {
  const response = await fetch(`/api/env/${route.params.env}/status/${route.params.id}`);
  if (!response.ok) {
    error.value = true;
    loading.value = false;
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  status.value = await response.json();
  await refreshIcon.value?.startAnimation();
}

watch(
    () => route.params.id, async (newId) => {
      loading.value = true;
      status.value = {} as StatusInfo;
      error.value = false;
      await fetchStatus();
      loading.value = false;
    }, {immediate: true});

const version = import.meta.env.VITE_GIT_COMMIT_VERSION as string;
</script>

<template>
  <div class="refresh-icon-container">
    <refresh-icon ref="refreshIcon" :callback="fetchStatus" class="image is-16x16"></refresh-icon>
  </div>

  <div class="container" v-if="loading || error">
    <div v-if="loading" class="title">Loading...</div>
    <div v-if="error" class="title">Error loading status</div>
    <span class="subtitle" v-if="error">
      Error loading status: {{ route.params.id }}
    </span>
  </div>

  <status-page-display v-if="!loading && !error" :status="status"></status-page-display>


  <footer class="footer">
    <div class="content has-text-centered">
      <p>
        <router-link to="/">All Status Pages</router-link>
      </p>
      <p>Powered by <a href="https://github.com/gitu/katastasi">katastasi</a></p>
      <p>{{ version }}</p>
    </div>
  </footer>
</template>


<style>

.refresh-icon-container {
  position: absolute;
  top: 0;
  right: 0;
  margin: 10px;
}

</style>