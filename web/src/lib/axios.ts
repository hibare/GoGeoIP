import axios, {
  type InternalAxiosRequestConfig,
  type AxiosResponse,
  type AxiosError,
  type AxiosRequestConfig,
} from "axios";
import { toast } from "vue-sonner";
import router from "@/router";

// Request interceptor
axios.interceptors.request.use(
  (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
    // Ensure all cookies are included in requests for authentication
    config.withCredentials = true;
    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

// Response interceptor
axios.interceptors.response.use(
  (response: AxiosResponse): AxiosResponse => {
    return response;
  },
  (error: AxiosError | Error): Promise<AxiosError> => {
    if (axios.isAxiosError(error)) {
      const { method, url } = error.config as AxiosRequestConfig;
      const { status } = (error.response as AxiosResponse) ?? {};

      // Handle specific status codes globally
      switch (status) {
        case 400:
          let errorMessage = "Bad Request";

          if (error.response?.data) {
            // Handle different error response formats
            if (typeof error.response.data === "string") {
              errorMessage = error.response.data;
            } else if (error.response.data?.message) {
              errorMessage = error.response.data.message;
            } else if (error.response.data?.error) {
              errorMessage = error.response.data.error;
            }
          }

          toast.error(`Validation Error: ${errorMessage}`);
          break;

        case 401:
          // Skip global handling for the initial auth check
          // This allows the router to handle the redirect gracefully without a full page reload
          if (url?.endsWith("/auth/me")) {
            return Promise.reject(error);
          }

          // Redirect to login page for GET requests on 401
          if (method?.toLowerCase() === "get") {
            router.push("/login");
            return Promise.reject(error);
          }

          toast.error("Authentication Required");
          break;

        case 403:
          toast.error("Access Denied");
          break;

        case 404:
          toast.warning("Not Found");
          break;

        case 422:
          let validationError = "Validation Error";
          if (error.response?.data) {
            if (typeof error.response.data === "string") {
              validationError = error.response.data;
            } else if (error.response.data?.message) {
              validationError = error.response.data.message;
            } else if (error.response.data?.error) {
              validationError = error.response.data.error;
            }
          }
          toast.error(validationError);
          break;

        case 500:
          toast.error("Server Error");
          break;

        case 503:
          toast.warning("Service Unavailable");
          break;

        default:
          // Generic error handling
          console.error(
            `Axios error on ${method?.toUpperCase()} ${url}: Status ${status}`,
            error,
          );
          toast.error("An error occurred. Please try again.");
          break;
      }
    } else {
      // Non-Axios error
      console.error("Non-Axios error:", error);
      toast.error("An unexpected error occurred.");
    }

    return Promise.reject(error);
  },
);

export const HTTP_STATUS_MULTIPLE_CHOICES = 300;
export const HTTP_STATUS_BAD_REQUEST = 400;

export default axios;
