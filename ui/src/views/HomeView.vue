<script setup lang="ts">
import {ref} from "vue";
import type {Environment} from "@/types";

let envList = ref([] as Environment[]);

fetch('/api/envs')
    .then(response => response.json())
    .then(data => {
      envList.value = data;
    });
</script>

<template>
  <div class="container">
    <div class="title py-5">Status Pages</div>


    <template v-for="env in envList">
      <div class="container block">
        <div class="subtitle">
          {{ env.name }}
        </div>

        <div class="content">
          <div class="buttons">
            <router-link class="button is-link is-outlined"
                         v-for="(status, statusId) in env.statusPages"
                         v-bind:to="'/env/'+env.name+'/status/'+ statusId">
              {{ status }}
            </router-link>
          </div>
        </div>
      </div>
    </template>


  </div>
  <footer class="footer">
    <div class="content has-text-centered">
      <p>Powered by <a href="https://github.com/gitu/katastasi">katastasi</a></p>
    </div>
  </footer>
</template>
