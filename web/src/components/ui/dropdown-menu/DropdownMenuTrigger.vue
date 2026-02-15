<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps<{
  asChild?: boolean
  class?: string
}>()
</script>

<template>
  <slot :class="cn('relative', props.class)" />
</template>
