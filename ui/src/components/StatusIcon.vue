<script setup lang="ts">

import {Critical, Info, OK, Unknown, Warning} from "@/types";
import CheckCircleOutline from "vue-material-design-icons/CheckCircleOutline.vue";
import InformationOutline from "vue-material-design-icons/InformationOutline.vue";
import AlertOutline from "vue-material-design-icons/AlertOutline.vue";
import AlertOctagram from "vue-material-design-icons/AlertOctagramOutline.vue";
import HelpCircleOutline from "vue-material-design-icons/HelpCircleOutline.vue";

const props = defineProps(['status'])
</script>

<template>
  <template v-if="status === OK">
    <CheckCircleOutline fill-color="green"></CheckCircleOutline>
  </template>
  <template v-else-if="status === Info">
    <InformationOutline fill-color="blue"></InformationOutline>
  </template>
  <template v-else-if="status === Warning">
    <AlertOutline fill-color="blue"></AlertOutline>
  </template>
  <template v-else-if="status === Critical">
    <AlertOctagram fill-color="red"></AlertOctagram>
  </template>
  <template v-else-if="status === Unknown">
    <HelpCircleOutline fill-color="grey"></HelpCircleOutline>
  </template>
  <template v-else>
    <HelpCircleOutline fill-color="grey"></HelpCircleOutline>
    {{ status }}
  </template>
</template>

<style scoped>

</style>