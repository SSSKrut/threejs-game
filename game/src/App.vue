<template>
  <div>
    <Game v-if="showGame" />
    <button @click="showGame = !showGame">Toggle Game</button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import Game from './logic/Game.vue';

const showGame = ref(false);
</script>

<style lang="scss">
* {
  box-sizing: border-box;
  font-weight: 400;
  padding: 0;
  margin: 0;
}

#app {
  background-color: black;
}

body {
  background-color: black;
}
</style>