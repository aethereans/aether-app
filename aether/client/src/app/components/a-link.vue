<template>
  <span>
    <template v-if="isAetherLink">
      <router-link class="a-link" :to="processedAetherRoute">
        {{ linktext }}
      </router-link>
    </template>
    <template v-else>
      <a class="a-link" :href="sanitisedLink" @click.stop>
        {{ linktext }}
      </a>
    </template>
  </span>
</template>

<script lang="ts">
export default {
  name: 'a-link',
  props: {
    link: {
      type: String,
      default: '',
    },
    linktext: {
      type: String,
      default: '',
    },
  },
  data() {
    return {}
  },
  computed: {
    isAetherLink(this: any) {
      if (typeof this.link !== 'string') {
        return false
      }
      return this.link.startsWith('aether://')
    },
    processedAetherRoute(this: any) {
      if (!this.isAetherLink) {
        return ''
      }
      return this.link.replace('aether://', '/')
      // ^ replaces only the first instance.
    },
    sanitisedLink(this: any) {
      if (typeof this.link === 'undefined') {
        return ''
      }
      if (this.link.length === 0) {
        return ''
      }
      if (
        this.link.substring(0, 8) === 'https://' ||
        this.link.substring(0, 7) === 'http://'
      ) {
        return this.link
      }
      return 'http://' + this.link
    },
  },
}
</script>

<style lang="scss" scoped></style>
