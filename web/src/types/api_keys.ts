export interface ApiKey {
  id: string;
  name: string;
  scopes: string[];
  expires_at: string | null;
  last_used_at: string | null;
  revoked_at: string | null;
  created_at: string;
  updated_at: string;
  state: "active" | "revoked" | "expired";
}

export interface CreateAPIKeyRequest {
  name: string;
  scopes?: string[];
  expires_at?: string;
}
