import type { GeoIP } from "@/types";
import axios from "@/lib/axios";

const ipEndpoint = "/api/v1/ip";

export const getMyIp = async (): Promise<GeoIP> => {
  const response = await axios.get<GeoIP>(`${ipEndpoint}`);
  return response.data;
};

export const lookupIp = async (ip: string): Promise<GeoIP> => {
  const response = await axios.get<GeoIP>(`${ipEndpoint}/${ip}`);
  return response.data;
};
