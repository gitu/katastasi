<script setup lang="ts">
import {RouterView} from 'vue-router'
import {ref} from "vue";
import type {StatusPage} from "@/types";

let burgerMenuHidden = true;
let statusList = ref([] as StatusPage[]);

fetch('/api/status')
    .then(response => response.json())
    .then(data => {
      statusList.value = data;
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

        <div class="navbar-item has-dropdown is-hoverable">
          <router-link class="navbar-item" v-for="status in statusList" v-bind:to="'/status/'+ status.id">
            {{ status.name }}
          </router-link>

        </div>
      </div>
    </div>
  </nav>
  <RouterView/>
</template>
