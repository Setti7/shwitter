import ApiError from "../models/errors/ApiError";
import User, { UserProfile } from "../models/user";
import { genericApiError, handleApiError } from "../utils/api";
import { apiService } from "./api";

export const getUser = async (id?: string): Promise<User | ApiError> => {
  const api = await apiService.getExecutor();

  try {
    const response = await api.get("users/" + (id === undefined ? "me" : id));
    return response.data["data"];
  } catch (error) {
    return handleApiError(error) ?? genericApiError;
  }
};

export const getUserProfile = async (id?: string): Promise<UserProfile | ApiError> => {
  const api = await apiService.getExecutor();

  try {
    const response = await api.get("users/" + (id === undefined ? "me" : id) + "/profile");
    return response.data["data"];
  } catch (error) {
    return handleApiError(error) ?? genericApiError;
  }
};

export const getIsFollowing = async (id: string): Promise<boolean | ApiError> => {
  const api = await apiService.getExecutor();

  try {
    const response = await api.get("users/" + id + "/follow");
    return response.data["data"];
  } catch (error) {
    return handleApiError(error) ?? genericApiError;
  }
};

export const followUser = async (id: string): Promise<undefined | ApiError> => {
  const api = await apiService.getExecutor();

  try {
    await api.post("users/" + id + "/follow");
    return;
  } catch (error) {
    return handleApiError(error) ?? genericApiError;
  }
};

export const unFollowUser = async (
  id: string
): Promise<undefined | ApiError> => {
  const api = await apiService.getExecutor();

  try {
    await api.post("users/" + id + "/unfollow");
    return;
  } catch (error) {
    return handleApiError(error) ?? genericApiError;
  }
};
