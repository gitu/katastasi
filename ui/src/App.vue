<script setup lang="ts">
import {RouterView} from 'vue-router'
import {ref} from "vue";
import type {Environment} from "@/types";

let burgerMenuHidden = true;
let envList = ref([] as Environment[]);

fetch('/api/envs')
    .then(response => response.json())
    .then(data => {
      envList.value = data;
    });
</script>

<template>

  <nav class="navbar" role="navigation" aria-label="main navigation">
    <div class="navbar-brand">
      <router-link class="navbar-item" to="/">
        <img src="@/assets/logo.png" width="32" height="32"> katastasi
      </router-link>

      <a role="button" class="navbar-burger" aria-label="menu" aria-expanded="false"
         @click="burgerMenuHidden = !burgerMenuHidden"
         :class="{'is-active':!burgerMenuHidden}"
         data-target="navbarBasicExample">
        <span aria-hidden="true"></span>
        <span aria-hidden="true"></span>
        <span aria-hidden="true"></span>
      </a>
    </div>

    <div id="navbarBasicExample" class="navbar-menu" :class="{'is-active':!burgerMenuHidden}">
      <div class="navbar-start">

        <template v-for="env in envList">
          <div class="navbar-item has-dropdown is-hoverable">
            <div class="navbar-item">
              {{ env.name }}
            </div>

            <div class="navbar-dropdown">
              <router-link class="navbar-item" v-for="(status, statusId) in env.statusPages"
                           v-bind:to="'/env/'+env.name+'/status/'+ statusId">
                {{ status }}
              </router-link>
            </div>
          </div>
        </template>
      </div>
    </div>
  </nav>
  <RouterView/>
</template>
