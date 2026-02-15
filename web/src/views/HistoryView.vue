<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Clock } from 'lucide-vue-next'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import type { LookupHistory } from '@/types'

const router = useRouter()

const history = ref<LookupHistory[]>([
  {
    id: '1',
    ip: '8.8.8.8',
    organization: 'Google LLC',
    location: 'Mountain View, US',
    timestamp: new Date(Date.now() - 2 * 60 * 1000),
  },
  {
    id: '2',
    ip: '1.1.1.1',
    organization: 'Cloudflare Inc',
    location: 'San Francisco, US',
    timestamp: new Date(Date.now() - 5 * 60 * 1000),
  },
  {
    id: '3',
    ip: '8.8.4.4',
    organization: 'Google LLC',
    location: 'Mountain View, US',
    timestamp: new Date(Date.now() - 10 * 60 * 1000),
  },
])

function formatTime(date: Date): string {
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / 60000)

  if (minutes < 1) return 'Just now'
  if (minutes < 60) return `${minutes}m ago`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}h ago`
  return date.toLocaleDateString()
}

function lookupAgain(ip: string) {
  router.push({ path: '/', query: { ip } })
}
</script>

<template>
  <div class="space-y-8">
    <div class="flex flex-col items-center justify-center space-y-4">
      <h1 class="text-3xl font-bold">Recent Lookups</h1>
      <p class="text-muted-foreground">Your recent IP lookups</p>
    </div>

    <Card>
      <CardContent class="p-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>IP Address</TableHead>
              <TableHead>Organization</TableHead>
              <TableHead>Location</TableHead>
              <TableHead>Time</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow
              v-for="item in history"
              :key="item.id"
              class="cursor-pointer"
              @click="lookupAgain(item.ip)"
            >
              <TableCell class="font-medium">{{ item.ip }}</TableCell>
              <TableCell>{{ item.organization }}</TableCell>
              <TableCell>{{ item.location }}</TableCell>
              <TableCell>
                <Badge variant="outline" class="flex items-center gap-1">
                  <Clock class="w-3 h-3" />
                  {{ formatTime(item.timestamp) }}
                </Badge>
              </TableCell>
            </TableRow>
            <TableRow v-if="history.length === 0">
              <TableCell colspan="4" class="text-center text-muted-foreground py-8">
                No recent lookups
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  </div>
</template>
