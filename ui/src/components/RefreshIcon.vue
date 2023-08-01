<script setup lang="ts">
import {onMounted, ref} from 'vue'

const props = defineProps({
  callback: {
    type: Function,
    required: true
  }
})

const animEnabled = ref(false)
const startAnimation = (): void => {
  animEnabled.value = true;

  setTimeout(() => {
    animEnabled.value = false;
    props.callback();
  }, 6000);
}

onMounted(() => {
})

defineExpose({startAnimation})
</script>

<template>
  <svg id="bars" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"
       viewBox="0 0 300 300" class="anim" v-if="animEnabled"
       shape-rendering="geometricPrecision" text-rendering="geometricPrecision">
    <rect ref="bar5" id="bars-u-5" width="300" height="50" rx="0" ry="0" transform="translate(0 250)" fill="#cbddff"
          stroke-width="0" stroke-linejoin="round"/>
    <rect id="bars-u-4" width="300" height="50" rx="0" ry="0" transform="translate(0 200)" fill="#cbddff"
          stroke-width="0" stroke-linejoin="round"/>
    <rect id="bars-u-3" width="300" height="50" rx="0" ry="0" transform="translate(0 150)" fill="#cbddff"
          stroke-width="0" stroke-linejoin="round"/>
    <rect id="bars-u-2" width="300" height="50" rx="0" ry="0" transform="translate(0 100)" fill="#cbddff"
          stroke-width="0" stroke-linejoin="round"/>
    <rect id="bars-u-1" width="300" height="50" rx="0" ry="0" transform="translate(0 50)" fill="#cbddff"
          stroke-width="0" stroke-linejoin="round"/>
    <rect id="bars-u-0" width="300" height="50" rx="0" ry="0" fill="#cbddff"
          stroke-width="0" stroke-linejoin="round"/>
  </svg>

</template>

<style scoped>
.anim #bars-u-5 {
  animation: bars-u-5_c_o 6000ms linear 1 normal forwards
}

@keyframes bars-u-5_c_o {
  0% {
    opacity: 1
  }
  83.333333% {
    opacity: 1;
    animation-timing-function: cubic-bezier(0.42, 0, 0.58, 1)
  }
  100% {
    opacity: 0
  }
}

.anim #bars-u-4 {
  animation: bars-u-4_c_o 6000ms linear 1 normal forwards
}

@keyframes bars-u-4_c_o {
  0% {
    opacity: 1
  }
  66.666667% {
    opacity: 1;
    animation-timing-function: cubic-bezier(0.42, 0, 0.58, 1)
  }
  83.333333% {
    opacity: 0
  }
  100% {
    opacity: 0
  }
}

.anim #bars-u-3 {
  animation: bars-u-3_c_o 6000ms linear 1 normal forwards
}

@keyframes bars-u-3_c_o {
  0% {
    opacity: 1
  }
  50% {
    opacity: 1
  }
  65% {
    opacity: 0
  }
  100% {
    opacity: 0
  }
}

.anim #bars-u-2 {
  animation: bars-u-2_c_o 6000ms linear 1 normal forwards
}

@keyframes bars-u-2_c_o {
  0% {
    opacity: 1
  }
  33.333333% {
    opacity: 1;
    animation-timing-function: cubic-bezier(0.42, 0, 0.58, 1)
  }
  50% {
    opacity: 0
  }
  100% {
    opacity: 0
  }
}

.anim #bars-u-1 {
  animation: bars-u-1_c_o 6000ms linear 1 normal forwards
}

@keyframes bars-u-1_c_o {
  0% {
    opacity: 1
  }
  16.666667% {
    opacity: 1;
    animation-timing-function: cubic-bezier(0.42, 0, 0.58, 1)
  }
  33.333333% {
    opacity: 0
  }
  100% {
    opacity: 0
  }
}

.anim #bars-u-0 {
  animation: bars-u-0_c_o 6000ms linear 1 normal forwards
}

@keyframes bars-u-0_c_o {
  0% {
    opacity: 1;
    animation-timing-function: cubic-bezier(0.42, 0, 0.58, 1)
  }
  16.666667% {
    opacity: 0
  }
  100% {
    opacity: 0
  }
}

</style>