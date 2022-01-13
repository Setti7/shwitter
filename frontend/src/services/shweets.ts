import ApiError from "../models/errors/ApiError";
import { Timeline } from "../models/shweet";
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
    const response = await api.post("shweets", { message });
    return response.data["data"];
  } catch (error) {
    return handleApiError(error) ?? genericApiError;
  }
};

export const getTimeline = async (): Promise<Timeline | ApiError> => {
  const api = await apiService.getExecutor();

  try {
    const response = await api.get("timeline");
    return response.data["data"];
  } catch (error) {
    return handleApiError(error) ?? genericApiError;
  }
};

export const getUserline = async (id: string): Promise<Timeline | ApiError> => {
  const api = await apiService.getExecutor();

  try {
    const response = await api.get("userline/" + id);
    return response.data["data"];
  } catch (error) {
    return handleApiError(error) ?? genericApiError;
  }
};
