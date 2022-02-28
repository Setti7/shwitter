import ApiError from "../models/errors/ApiError";
import { ShweetDetails, Timeline } from "../models/shweet";
import { genericApiError, handleApiError } from "../utils/api";
import { apiService } from "./api";

export type ShweetID = string;

interface CreateShweetProps {
  message: string;
}

export const createShweet = async ({
  message,
}: CreateShweetProps): Promise<ShweetID | ApiError> => {
  const api = await apiService.getExecutor();

  try {
    const response = await api.post("v1/timeline", { message });
    return response.data["data"];
  } catch (error) {
    return handleApiError(error) ?? genericApiError;
  }
};

export const getShweetDetails = async (
  shweetId: string
): Promise<ShweetDetails | ApiError> => {
  const api = await apiService.getExecutor();

  try {
    const response = await api.get("v1/shweets/" + shweetId);
    return response.data["data"];
  } catch (error) {
    return handleApiError(error) ?? genericApiError;
  }
};

export const getTimeline = async (): Promise<Timeline | ApiError> => {
  const api = await apiService.getExecutor();

  try {
    const response = await api.get("v1/timeline");
    return response.data["data"];
  } catch (error) {
    return handleApiError(error) ?? genericApiError;
  }
};

export const getUserline = async (id: string): Promise<Timeline | ApiError> => {
  const api = await apiService.getExecutor();

  try {
    const response = await api.get("v1/userline/" + id);
    return response.data["data"];
  } catch (error) {
    return handleApiError(error) ?? genericApiError;
  }
};

export const likeShweet = async (shweetID: ShweetID): Promise<ShweetID | ApiError> => {
  const api = await apiService.getExecutor();

  try {
    const response = await api.post("v1/shweets/" + shweetID + "/like");
    return response.data["data"];
  } catch (error) {
    return handleApiError(error) ?? genericApiError;
  }
};
