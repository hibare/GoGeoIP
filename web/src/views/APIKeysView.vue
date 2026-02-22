<template>
  <div class="space-y-8">
    <div class="flex flex-col space-y-1">
      <h1 class="text-3xl font-bold">API Keys</h1>
      <p class="text-muted-foreground">Manage your API keys for programmatic access</p>
    </div>

    <div class="flex justify-end">
      <Button
        :disabled="apiKeysStore.createLoading"
        @click="showCreateDialog = true"
      >
        Generate New Key
      </Button>
    </div>

    <div v-if="apiKeysStore.listLoading" class="flex justify-center py-8">
      <Loader2Icon class="h-8 w-8 animate-spin" />
    </div>

    <DataTable v-else :columns="columns" :data="apiKeys" :page-size="10" />

    <Dialog v-model:open="showDeleteDialog">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Delete API Key</DialogTitle>
          <DialogDescription>
            Are you sure you want to delete this API key? This action cannot be
            undone.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" @click="showDeleteDialog = false"
            >Cancel</Button
          >
          <Button
            variant="destructive"
            :disabled="apiKeysStore.deleteLoading"
            @click="confirmDelete"
          >
            <Loader2Icon
              v-if="apiKeysStore.deleteLoading"
              class="mr-2 h-4 w-4 animate-spin"
            />
            Delete
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <Dialog v-model:open="showRevokeDialog">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Revoke API Key</DialogTitle>
          <DialogDescription>
            Are you sure you want to revoke this API key? The key will no longer
            be usable.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" @click="showRevokeDialog = false"
            >Cancel</Button
          >
          <Button :disabled="apiKeysStore.revokeLoading" @click="confirmRevoke">
            <Loader2Icon
              v-if="apiKeysStore.revokeLoading"
              class="mr-2 h-4 w-4 animate-spin"
            />
            Revoke
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <Dialog v-model:open="showCreateDialog">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create API Key</DialogTitle>
          <DialogDescription>
            Generate a new API key for programmatic access.
          </DialogDescription>
        </DialogHeader>
        <div class="space-y-4 py-4">
          <div class="space-y-2">
            <Label for="key-name">Name</Label>
            <Input
              id="key-name"
              v-model="newKeyName"
              placeholder="My API Key"
            />
          </div>
          <div class="space-y-2">
            <Label for="key-expires">Expires At (optional)</Label>
            <Input
              id="key-expires"
              type="datetime-local"
              v-model="newKeyExpiresAt"
            />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="showCreateDialog = false"
            >Cancel</Button
          >
          <Button
            :disabled="apiKeysStore.createLoading || !newKeyName"
            @click="createKey"
          >
            <Loader2Icon
              v-if="apiKeysStore.createLoading"
              class="mr-2 h-4 w-4 animate-spin"
            />
            Create
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <Dialog v-model:open="showKeyDialog">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>API Key Created</DialogTitle>
          <DialogDescription>
            Copy this key now. You won't be able to see it again!
          </DialogDescription>
        </DialogHeader>
        <div class="space-y-4">
          <div class="p-3 border-2 border-dashed border-amber-500/50 rounded-md bg-amber-50/50">
            <div class="flex items-center gap-2">
              <code class="flex-1 text-sm font-mono break-all text-amber-900">
                {{ maskedKey }}
              </code>
              <Button variant="outline" size="sm" @click="copyKey">
                <CopyIcon class="w-4 h-4" />
              </Button>
            </div>
          </div>
          <p class="text-xs text-muted-foreground">
            Click the copy button to copy the full API key
          </p>
        </div>
        <DialogFooter>
          <Button @click="showKeyDialog = false">Done</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, h } from "vue";
import { useTimeAgo, useClipboard } from "@vueuse/core";
import { toast } from "vue-sonner";
import { Loader2Icon, Trash2Icon, BanIcon, ClockIcon, CopyIcon } from "lucide-vue-next";
import { createColumnHelper } from "@tanstack/vue-table";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import DataTable from "@/components/DataTable.vue";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { useAPIKeysStore } from "@/store/api_keys";
import { formatDateTime } from "@/lib/utils";
import type { ApiKey } from "@/types/api_keys";

const apiKeysStore = useAPIKeysStore();

const showDeleteDialog = ref(false);
const showRevokeDialog = ref(false);
const showCreateDialog = ref(false);
const showKeyDialog = ref(false);
const selectedKeyId = ref<string | null>(null);
const newKeyName = ref("");
const newKeyExpiresAt = ref("");
const newCreatedKey = ref("");

const { copy, isSupported: clipboardSupported } = useClipboard();

const maskedKey = computed(() => {
  const key = newCreatedKey.value;
  if (!key || key.length <= 16) return key;
  return `${key.substring(0, 8)}...${key.substring(key.length - 8)}`;
});

const apiKeys = computed(() => apiKeysStore.apiKeys);

const columnHelper = createColumnHelper<ApiKey>();

const columns = computed(() => [
  columnHelper.accessor("name", {
    header: "Name",
    cell: ({ row }) => h("span", { class: "font-medium" }, row.original.name),
  }),
  columnHelper.accessor("state", {
    header: "Status",
    cell: ({ row }) => {
      const state = row.original.state;
      const variant =
        state === "active"
          ? "default"
          : state === "expired"
            ? "secondary"
            : "destructive";
      return h(Badge, { variant }, () => state);
    },
  }),
  columnHelper.accessor("scopes", {
    header: "Scopes",
    cell: ({ row }) => {
      const scopes = row.original.scopes;
      if (!scopes || scopes.length === 0)
        return h(Badge, { variant: "secondary" }, () => "All");
      return h(
        "div",
        { class: "flex gap-1 flex-wrap" },
        scopes.map((s) =>
          h(Badge, { key: s, variant: "outline", class: "text-xs" }, () => s),
        ),
      );
    },
  }),
  columnHelper.accessor("expires_at", {
    header: "Expires",
    cell: ({ row }) => {
      const expires = row.original.expires_at;
      if (!expires) return h(Badge, { variant: "secondary" }, () => "Never");
      const isExpired = new Date(expires) < new Date();
      const timeAgo = useTimeAgo(expires);
      return h(Tooltip, {}, () => [
        h(TooltipTrigger, { asChild: true }, () =>
          h(
            Badge,
            {
              variant: isExpired ? "destructive" : "outline",
              class: "flex items-center gap-1 cursor-default",
            },
            () => [h(ClockIcon, { class: "w-3 h-3 mr-1" }), timeAgo.value],
          ),
        ),
        h(TooltipContent, {}, () => formatDateTime(expires as string)),
      ]);
    },
  }),
  columnHelper.accessor("created_at", {
    header: "Created",
    cell: ({ row }) => {
      const timeAgo = useTimeAgo(row.original.created_at);
      return h(Tooltip, {}, () => [
        h(TooltipTrigger, { asChild: true }, () =>
          h(
            Badge,
            {
              variant: "outline",
              class: "flex items-center gap-1 cursor-default",
            },
            () => [h(ClockIcon, { class: "w-3 h-3" }), timeAgo.value],
          ),
        ),
        h(TooltipContent, {}, () => formatDateTime(row.original.created_at)),
      ]);
    },
  }),
  columnHelper.accessor("last_used_at", {
    header: "Last Used",
    cell: ({ row }) => {
      if (!row.original.last_used_at)
        return h(Badge, { variant: "outline" }, () => "-");
      const timeAgo = useTimeAgo(row.original.last_used_at);
      return h(Tooltip, {}, () => [
        h(TooltipTrigger, { asChild: true }, () =>
          h(
            Badge,
            {
              variant: "outline",
              class: "flex items-center gap-1 cursor-default",
            },
            () => [h(ClockIcon, { class: "w-3 h-3" }), timeAgo.value],
          ),
        ),
        h(TooltipContent, {}, () =>
          formatDateTime(row.original.last_used_at as string),
        ),
      ]);
    },
  }),
  columnHelper.display({
    id: "actions",
    header: "",
    cell: ({ row }) => {
      const key = row.original;
      const isRevoked = key.state === "revoked";
      return h("div", { class: "flex gap-1" }, [
        !isRevoked &&
          h(
            Button,
            {
              variant: "outline",
              size: "sm",
              onClick: () => {
                selectedKeyId.value = key.id;
                showRevokeDialog.value = true;
              },
            },
            () => h(BanIcon, { class: "w-4 h-4" }),
          ),
        h(
          Button,
          {
            variant: "ghost",
            size: "sm",
            onClick: () => {
              selectedKeyId.value = key.id;
              showDeleteDialog.value = true;
            },
          },
          () => h(Trash2Icon, { class: "w-4 h-4" }),
        ),
      ]);
    },
  }),
]);

onMounted(() => {
  apiKeysStore.fetchAPIKeys();
});

async function confirmDelete() {
  if (selectedKeyId.value) {
    await apiKeysStore.deleteAPIKey(selectedKeyId.value);
  }
  showDeleteDialog.value = false;
  selectedKeyId.value = null;
}

async function confirmRevoke() {
  if (selectedKeyId.value) {
    await apiKeysStore.revokeAPIKey(selectedKeyId.value);
  }
  showRevokeDialog.value = false;
  selectedKeyId.value = null;
}

async function createKey() {
  if (newKeyName.value) {
    const request: { name: string; expires_at?: string } = {
      name: newKeyName.value,
    };
    if (newKeyExpiresAt.value) {
      request.expires_at = new Date(newKeyExpiresAt.value).toISOString();
    }
    const newKey = await apiKeysStore.createAPIKey(request);
    newCreatedKey.value = newKey;
    newKeyName.value = "";
    newKeyExpiresAt.value = "";
    showCreateDialog.value = false;
    showKeyDialog.value = true;
  }
}

async function copyKey() {
  if (!clipboardSupported.value) {
    toast.error("Clipboard API not supported");
    return;
  }
  try {
    await copy(newCreatedKey.value);
    toast.success("API key copied to clipboard");
  } catch (error) {
    toast.error("Failed to copy API key");
  }
}
</script>
