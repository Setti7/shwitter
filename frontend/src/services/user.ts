import ApiError from "../models/errors/ApiError";
import User, { FriendOrFollower, UserProfile } from "../models/user";
import { genericApiError, handleApiError } from "../utils/api";
import { apiService } from "./api";

export const getUser = async (id?: string): Promise<User | ApiError> => {
  const api = await apiService.getExecutor();

  try {
    const response = await api.get("v1/users/" + (id === undefined ? "me" : id));
    return response.data["data"];
  } catch (error) {
    return handleApiError(error) ?? genericApiError;
  }
};

export const getUserProfile = async (
  id?: string
): Promise<UserProfile | ApiError> => {
  const api = await apiService.getExecutor();

  try {
    const response = await api.get(
      "v1/users/" + (id === undefined ? "me" : id) + "/profile"
    );
    return response.data["data"];
  } catch (error) {
    return handleApiError(error) ?? genericApiError;
  }
};

export const getIsFollowing = async (
  id: string
): Promise<boolean | ApiError> => {
  const api = await apiService.getExecutor();

  try {
    const response = await api.get("v1/users/" + id + "/follow");
    return response.data["data"];
  } catch (error) {
    return handleApiError(error) ?? genericApiError;
  }
};

export const getFollowers = async (
  id: string
): Promise<FriendOrFollower[] | ApiError> => {
  const api = await apiService.getExecutor();

  try {
    const response = await api.get("v1/users/" + id + "/followers");
    return response.data["data"];
  } catch (error) {
    return handleApiError(error) ?? genericApiError;
  }
};

export const getFriends = async (
  id: string
): Promise<FriendOrFollower[] | ApiError> => {
  const api = await apiService.getExecutor();

  try {
    const response = await api.get("v1/users/" + id + "/friends");
    return response.data["data"];
  } catch (error) {
    return handleApiError(error) ?? genericApiError;
  }
};

export const followOrUnfollowUser = async (id: string): Promise<undefined | ApiError> => {
  const api = await apiService.getExecutor();

  try {
    await api.post("v1/users/" + id + "/follow");
    return;
  } catch (error) {
    return handleApiError(error) ?? genericApiError;
  }
};
