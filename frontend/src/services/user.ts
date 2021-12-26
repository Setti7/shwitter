import ApiError from "../models/errors/ApiError";
import User from "../models/user";
import { genericApiError, handleApiError } from "../utils/api";
import { apiService } from "./api";

export const getUser = async (): Promise<User | ApiError> => {
  const api = await apiService.getExecutor();

  try {
    const response = await api.get("users/me");
    return response.data['data'];
  } catch (error) {
    return handleApiError(error) ?? genericApiError;
  }
};
