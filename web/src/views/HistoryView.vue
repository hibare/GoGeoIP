<template>
  <div class="space-y-8">
    <div class="flex items-center justify-between">
      <div class="flex flex-col space-y-1">
        <h1 class="text-3xl font-bold">Recent Lookups</h1>
        <p class="text-muted-foreground">Your recent IP lookups</p>
      </div>
      <Button
        v-if="history.length > 0"
        variant="outline"
        size="sm"
        @click="clearAll"
      >
        <Trash2Icon class="w-4 h-4 mr-2" />
        Clear All
      </Button>
    </div>

    <DataTable
      v-if="history.length > 0"
      :columns="columns"
      :data="history"
      :page-size="10"
    />

    <div v-else class="text-center py-12 text-muted-foreground">
      No recent lookups. Start by looking up an IP address.
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, h } from "vue";
import { useRouter } from "vue-router";
import { useTimeAgo } from "@vueuse/core";
import { ClockIcon, Trash2Icon } from "lucide-vue-next";
import { createColumnHelper } from "@tanstack/vue-table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import DataTable from "@/components/DataTable.vue";
import { useHistoryStore, type LookupHistory } from "@/store/history";

const router = useRouter();
const historyStore = useHistoryStore();

const history = computed(() => historyStore.history);

const columnHelper = createColumnHelper<LookupHistory>();

const columns = computed(() => [
  columnHelper.accessor("ip", {
    header: "IP Address",
    cell: ({ row }) => h("span", { class: "font-medium" }, row.original.ip),
  }),
  columnHelper.accessor("organization", {
    header: "Organization",
    cell: ({ row }) => row.original.organization || "-",
  }),
  columnHelper.accessor("location", {
    header: "Location",
    cell: ({ row }) => row.original.location || "-",
  }),
  columnHelper.accessor("first_lookup", {
    header: "First Lookup",
    cell: ({ row }) => {
      const timestamp = row.original.first_lookup;
      const timeAgo = useTimeAgo(timestamp);
      return h(
        Badge,
        { variant: "secondary", class: "flex items-center gap-1" },
        () => [h(ClockIcon, { class: "w-3 h-3" }), timeAgo.value],
      );
    },
  }),
  columnHelper.accessor("timestamp", {
    header: "Last Lookup",
    cell: ({ row }) => {
      const timestamp = row.original.timestamp;
      const timeAgo = useTimeAgo(timestamp);
      return h(
        Badge,
        { variant: "outline", class: "flex items-center gap-1" },
        () => [h(ClockIcon, { class: "w-3 h-3" }), timeAgo.value],
      );
    },
  }),
  columnHelper.accessor("times_looked", {
    header: "Count",
    cell: ({ row }) => {
      const count = row.original.times_looked;
      return h(
        Badge,
        { variant: count > 1 ? "default" : "outline" },
        () => count,
      );
    },
  }),
  columnHelper.display({
    id: "actions",
    header: "",
    cell: ({ row }) => {
      return h(
        Button,
        {
          variant: "ghost",
          size: "sm",
          class: "h-8 w-8 p-0",
          onClick: (e: Event) => {
            e.stopPropagation();
            historyStore.removeEntry(row.original.id);
          },
        },
        () => h(Trash2Icon, { class: "w-4 h-4" }),
      );
    },
  }),
]);

function lookupAgain(ip: string) {
  router.push({ path: "/", query: { ip } });
}

function clearAll() {
  historyStore.clearHistory();
}
</script>
