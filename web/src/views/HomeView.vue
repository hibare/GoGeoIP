<template>
  <div class="space-y-8">
    <div class="flex flex-col items-center justify-center space-y-4">
      <h1 class="text-3xl font-bold">IP Geolocation Lookup</h1>
      <p class="text-muted-foreground">Look up any IP address</p>
    </div>

    <Card class="w-full max-w-2xl mx-auto">
      <CardContent class="pt-6">
        <form @submit.prevent="handleLookup" class="flex gap-2">
          <Input
            v-model="ipInput"
            type="text"
            placeholder="Enter IP address (e.g., 8.8.8.8)"
            class="flex-1"
            :disabled="isLoading"
          />
          <Button type="submit" :disabled="isLoading">
            <Loader2Icon v-if="isLoading" class="w-4 h-4 mr-2 animate-spin" />
            <SearchIcon v-else class="w-4 h-4 mr-2" />
            Lookup
          </Button>
          <Button
            type="button"
            variant="secondary"
            @click="handleMyIp"
            :disabled="isLoading"
          >
            <Loader2Icon
              v-if="isLoading && !result"
              class="w-4 h-4 mr-2 animate-spin"
            />
            My IP
          </Button>
        </form>
      </CardContent>
    </Card>

    <div v-if="isLoading" class="flex justify-center py-12">
      <Loader2Icon class="h-8 w-8 animate-spin text-primary" />
    </div>

    <div v-else-if="error" class="text-destructive text-center">
      {{ error }}
    </div>

    <div v-else-if="result" class="grid gap-6 md:grid-cols-2">
      <Card>
        <CardHeader>
          <CardTitle class="flex items-center gap-2">
            <MapPinIcon class="w-5 h-5" />
            Location
          </CardTitle>
        </CardHeader>
        <CardContent class="space-y-2">
          <div class="flex justify-between items-center">
            <span class="text-muted-foreground">Country</span>
            <span class="font-medium flex items-center gap-2">
              <span v-if="getCountryFlag(result.iso_country_code)">{{ getCountryFlag(result.iso_country_code) }}</span>
              {{ result.country || "N/A" }}
            </span>
          </div>
          <div class="flex justify-between">
            <span class="text-muted-foreground">City</span>
            <span class="font-medium">{{ result.city || "N/A" }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-muted-foreground">Continent</span>
            <span class="font-medium">{{ result.continent || "N/A" }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-muted-foreground">Timezone</span>
            <span class="font-medium">{{ result.timezone || "N/A" }}</span>
          </div>
          <Separator />
          <div class="flex justify-between">
            <span class="text-muted-foreground">Coordinates</span>
            <span class="font-medium"
              >{{ result.latitude }}, {{ result.longitude }}</span
            >
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle class="flex items-center gap-2">
            <NetworkIcon class="w-5 h-5" />
            Network
          </CardTitle>
        </CardHeader>
        <CardContent class="space-y-2">
          <div class="flex justify-between">
            <span class="text-muted-foreground">ASN</span>
            <span class="font-medium">{{ result.asn || "N/A" }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-muted-foreground">Organization</span>
            <span class="font-medium">{{ result.organization || "N/A" }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-muted-foreground">IP Address</span>
            <span class="font-medium">{{ result.ip }}</span>
          </div>
        </CardContent>
      </Card>

      <Card class="md:col-span-2">
        <CardHeader>
          <CardTitle class="flex items-center gap-2">
            <ShieldIcon class="w-5 h-5" />
            Details
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div class="flex flex-wrap gap-2">
            <Badge
              :variant="result.is_anonymous_proxy ? 'destructive' : 'secondary'"
            >
              <SatelliteIcon class="w-3 h-3 mr-1" />
              Anonymous Proxy: {{ result.is_anonymous_proxy ? "Yes" : "No" }}
            </Badge>
            <Badge
              :variant="
                result.is_satellite_provider ? 'destructive' : 'secondary'
              "
            >
              <GlobeIcon class="w-3 h-3 mr-1" />
              Satellite: {{ result.is_satellite_provider ? "Yes" : "No" }}
            </Badge>
            <Badge variant="outline">
              {{ result.iso_continent_code }} / {{ result.iso_country_code }}
            </Badge>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { useRouter, useRoute } from "vue-router";
import {
  SearchIcon,
  Loader2Icon,
  MapPinIcon,
  GlobeIcon,
  NetworkIcon,
  ShieldIcon,
  SatelliteIcon,
} from "lucide-vue-next";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { useUserStore } from "@/store/auth";
import { useHistoryStore } from "@/store/history";
import type { GeoIP } from "@/types";
import { getMyIp, lookupIp } from "@/apis/ip";
import { getCountryFlag } from "@/lib/flags";

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();
const historyStore = useHistoryStore();

const ipInput = ref("");
const isLoading = ref(false);
const error = ref("");
const result = ref<GeoIP | null>(null);

onMounted(async () => {
  const queryIp = route.query.ip as string;
  if (queryIp) {
    ipInput.value = queryIp;
    await handleLookup();
  } else {
    await handleMyIp();
  }
});

watch(ipInput, (newVal) => {
  if (newVal) {
    router.replace({ query: { ip: newVal } });
  } else {
    router.replace({ query: {} });
  }
});

async function handleLookup() {
  if (!ipInput.value.trim()) return;

  if (!userStore.isAuthenticated) {
    router.push({ path: "/login", query: { redirect: route.fullPath } });
    return;
  }

  isLoading.value = true;
  error.value = "";
  result.value = null;

  try {
    const data = await lookupIp(ipInput.value);
    result.value = data;
    router.replace({ query: { ip: ipInput.value } });

    const location = data.city && data.country
      ? `${data.city}, ${data.country}`
      : data.country || "Unknown";
    historyStore.addEntry(data.ip, data.organization || "", location, data.iso_country_code);
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Failed to lookup IP";
  } finally {
    isLoading.value = false;
  }
}

async function handleMyIp() {
  isLoading.value = true;
  error.value = "";
  result.value = null;

  try {
    const data = await getMyIp();
    result.value = data;
    ipInput.value = data.ip;
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Failed to get my IP";
  } finally {
    isLoading.value = false;
  }
}
</script>
