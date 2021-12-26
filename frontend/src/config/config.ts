import { strToBool } from "../utils/strings";

// The URL for the api. Example: app.fincashcorban.com.br/api/. Required.
export const API_ORIGIN = process.env.REACT_APP_API_ORIGIN;

// Add debug features. Default: false.
export const DEBUG: boolean = strToBool(process.env.REACT_APP_DEBUG ?? "false");
