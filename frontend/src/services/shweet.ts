import ApiError from "../models/errors/ApiError";
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
