import { type User } from "@/types/auth";

import axios, {
  HTTP_STATUS_MULTIPLE_CHOICES,
  HTTP_STATUS_BAD_REQUEST,
} from "@/lib/axios";

const authEndpoint = "/api/v1/auth";

export const login = async (
  redirect?: string,
): Promise<{ redirectUrl: string }> => {
  const params = redirect ? { redirect } : {};
  const response = await axios.get(`${authEndpoint}/login`, { params });
  if (response.data && response.data.redirect_url) {
    return { redirectUrl: response.data.redirect_url };
  }
  throw new Error("Login failed: Invalid response from server");
};

export const logout = async (): Promise<{ redirectUrl?: string }> => {
  const response = await axios.post(`${authEndpoint}/logout`, {});
  if (
    response.status >= HTTP_STATUS_MULTIPLE_CHOICES &&
    response.status < HTTP_STATUS_BAD_REQUEST &&
    response.headers.location
  ) {
    return { redirectUrl: response.headers.location };
  }
  return {};
};

export const getProfile = async (): Promise<User> => {
  const response = await axios.get<User>(`${authEndpoint}/me`);
  return response.data;
};
