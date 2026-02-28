import { defineStore } from "pinia";
import { ref } from "vue";
import type { ApiKey, CreateAPIKeyRequest } from "@/types/api_keys";
import {
  listAPIKeys,
  createAPIKey,
  revokeAPIKey,
  deleteAPIKey,
} from "@/apis/api_keys";
import { toast } from "vue-sonner";

export const useAPIKeysStore = defineStore("api-keys", () => {
  const apiKeys = ref<ApiKey[]>([]);
  const listLoading = ref(false);
  const createLoading = ref(false);
  const revokeLoading = ref(false);
  const deleteLoading = ref(false);

  const fetchAPIKeys = async () => {
    listLoading.value = true;
    try {
      apiKeys.value = await listAPIKeys();
    } catch (error) {
      console.error("Failed to fetch API keys:", error);
      throw error;
    } finally {
      listLoading.value = false;
    }
  };

  const createAPIKeyAction = async (
    data: CreateAPIKeyRequest,
  ): Promise<string> => {
    createLoading.value = true;
    try {
      const result = await createAPIKey(data);
      toast.success("API Key Created");
      fetchAPIKeys();
      return result;
    } finally {
      createLoading.value = false;
    }
  };

  const revokeAPIKeyAction = async (id: string) => {
    revokeLoading.value = true;
    try {
      await revokeAPIKey(id);
      toast.success("API key revoked successfully");
      fetchAPIKeys();
    } catch (error) {
      console.error("Failed to revoke API key:", error);
      throw error;
    } finally {
      revokeLoading.value = false;
    }
  };

  const deleteAPIKeyAction = async (id: string) => {
    deleteLoading.value = true;
    try {
      await deleteAPIKey(id);
      toast.success("API key deleted successfully");
      fetchAPIKeys();
    } catch (error) {
      console.error("Failed to delete API key:", error);
      throw error;
    } finally {
      deleteLoading.value = false;
    }
  };

  return {
    // State
    apiKeys,
    listLoading,
    createLoading,
    revokeLoading,
    deleteLoading,

    // Actions
    fetchAPIKeys,
    createAPIKey: createAPIKeyAction,
    revokeAPIKey: revokeAPIKeyAction,
    deleteAPIKey: deleteAPIKeyAction,
  };
});
