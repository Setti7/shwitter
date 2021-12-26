import ApiError from "../models/errors/ApiError";
import AuthError from "../models/errors/AuthError";

export const genericApiError = new ApiError({
  error: "An unexpected error occurred.",
});

export const handleApiError = (error: any): ApiError | undefined => {
  if (!error.response) {
    console.log(error);
    return new ApiError({
      error: "Connection error. Are you connected to the internet?",
    });
  } else if (error.response.status === 500 || error.response.status === 503) {
    if (error.response.data.detail) {
      return new ApiError({
        error: error.response.data.detail,
      });
    } else {
      return genericApiError;
    }
  } else if (error.response.status === 400) {
    return new ApiError(error.response.data);
  } else if (error.response.status === 401) {
    return new AuthError(error.response.data);
  } else if (error.response.status === 403) {
    return new ApiError(error.response.data);
  }
};
