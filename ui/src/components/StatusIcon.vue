<script setup lang="ts">

import {Critical, Info, OK, Unknown, Warning} from "@/types";

const props = defineProps(['status'])
</script>

<template>
  <template v-if="status === OK">
    <span class="icon-text">
      <span class="icon has-text-info">
      <CheckCircleOutline fill-color="green"></CheckCircleOutline>
    </span>
      <span>OK</span>
    </span>
  </template>
  <template v-else-if="status === Info">
    <span class="icon-text">
      <span class="icon has-text-info">
        <InformationOutline fill-color="blue"></InformationOutline>
      </span>
      <span>Info</span>
    </span>
  </template>
  <template v-else-if="status === Warning">
    <span class="icon-text">
      <span class="icon has-text-info">
        <AlertOutline fill-color="blue"></AlertOutline>
      </span>
      <span>Warning</span>
    </span>
  </template>
  <template v-else-if="status === Critical">
    <span class="icon-text">
      <span class="icon has-text-info">
        <AlertOctagram fill-color="red"></AlertOctagram>
      </span>
      <span>Critical</span>
    </span>
  </template>
  <template v-else-if="status === Unknown">
    <span class="icon-text">
      <span class="icon has-text-info">
        <HelpCircleOutline fill-color="grey"></HelpCircleOutline>
      </span>
      <span>Unknown</span>
    </span>
  </template>
  <template v-else>
    <span class="icon-text">
      <span class="icon has-text-info">
        <HelpCircleOutline fill-color="grey"></HelpCircleOutline>
      </span>
      <span>{{ status }}</span>
    </span>
  </template>
</template>

<style scoped>

</style>