import ApiError from "../models/errors/ApiError";
import AuthError from "../models/errors/AuthError";
import Session from "../models/session";
import User from "../models/user";
import { genericApiError, handleApiError } from "../utils/api";
import { apiService } from "./api";

const SESSION_KEY = "sess";
const USER_KEY = "user";

export const login = async (
  username: string,
  password: string
): Promise<Session | AuthError> => {
  const api = await apiService.getExecutor();

  try {
    const response = await api.post("sessions", {
      username,
      password,
    });
    const sess: Session = response.data["data"];
    return sess;
  } catch (error) {
    return handleApiError(error) ?? genericApiError;
  }
};

interface ICreateUser {
  username: string;
  name: string;
  email: string;
  password: string;
}

export const createUser = async ({
  username,
  name,
  email,
  password,
}: ICreateUser): Promise<User | ApiError> => {
  const api = await apiService.getExecutor();

  try {
    const response = await api.post("users", {
      username,
      name,
      email,
      password,
    });
    const user: User = response.data["data"];
    return user;
  } catch (error) {
    return handleApiError(error) ?? genericApiError;
  }
};

export const saveSession = (session: Session | undefined): void => {
  if (session) {
    localStorage.setItem(SESSION_KEY, JSON.stringify(session));
  } else {
    localStorage.removeItem(SESSION_KEY);
  }
};

export const loadSession = (): Session | undefined => {
  const sess = localStorage.getItem(SESSION_KEY);
  if (sess) {
    return JSON.parse(sess) as Session;
  }
};

export const saveUser = (user: User | undefined): void => {
  if (user) {
    localStorage.setItem(USER_KEY, JSON.stringify(user));
  } else {
    localStorage.removeItem(USER_KEY);
  }
};

export const loadUser = (): User | undefined => {
  const user = localStorage.getItem(USER_KEY);
  if (user) {
    return JSON.parse(user) as User;
  }
};
