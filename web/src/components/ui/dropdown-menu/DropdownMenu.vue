<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps<{
  class?: string
}>()

const delegatedProps = computed(() => {
  const { class: _, ...delegated } = props
  return delegated
})
</script>

<template>
  <div :class="cn('relative inline-block text-left', props.class)">
    <slot />
  </div>
</template>
