import axios from "@/lib/axios";
import type { ApiKey, CreateAPIKeyRequest } from "@/types/api_keys";

const apiKeysEndpoint = "/api/v1/api-keys";
const apiKeyEndpoint = "/api/v1/api-key";

export const listAPIKeys = async (): Promise<Array<ApiKey>> => {
  const response = await axios.get<Array<ApiKey>>(apiKeysEndpoint);
  return response.data;
};

export const createAPIKey = async (
  request: CreateAPIKeyRequest,
): Promise<string> => {
  const response = await axios.post<string>(apiKeyEndpoint, request);
  return response.data;
};

export const revokeAPIKey = async (id: string): Promise<void> => {
  await axios.post(`${apiKeyEndpoint}/${id}/revoke`);
};

export const deleteAPIKey = async (id: string): Promise<void> => {
  await axios.delete(`${apiKeyEndpoint}/${id}`);
};
