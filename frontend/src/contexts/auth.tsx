import React, { createContext, useCallback, useEffect, useState } from "react";
import ApiError from "../models/errors/ApiError";
import AuthError from "../models/errors/AuthError";
import User from "../models/user";
import { apiService } from "../services/api";
import { loadUser, login, logout, saveUser } from "../services/auth";
import { getUser } from "../services/user";

export enum AuthStatus {
  Loading = "Loading",
  Authenticated = "Authenticated",
  NotAuthenticated = "NotAuthenticated",
}

type IAuthFetch = (args0?: { force: boolean }) => Promise<void>;

interface AuthContextData {
  user?: User;
  authStatus: AuthStatus;
  authLogin: (
    username: string,
    password: string
  ) => Promise<AuthError | undefined>;
  authLogout: () => void;
  authFetch: IAuthFetch;
}

export const AuthContext = createContext<AuthContextData>(
  {} as AuthContextData
);

export const AuthProvider: React.FC = ({ children }) => {
  const [user, setUser] = useState<User | undefined>(undefined);
  const [authStatus, setAuthStatus] = useState<AuthStatus>(AuthStatus.Loading);

  // Load user from LocalStorage on startup
  useEffect(() => {
    const _user = loadUser();
    if (_user) {
      updateUser(_user);

      authFetch({ force: true }); // Refresh user data from server
      setAuthStatus(AuthStatus.Authenticated);
    } else {
      setAuthStatus(AuthStatus.NotAuthenticated);
    }
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  const updateUser = (user: User | undefined) => {
    setUser(user);
    saveUser(user); // save to LocalStorage
  };

  const authLogin = useCallback(
    async (
      username: string,
      password: string
    ): Promise<ApiError | AuthError | undefined> => {
      const sessResult = await login(username, password);

      if (sessResult instanceof ApiError) {
        return sessResult;
      }

      apiService.authorize(sessResult);

      const userResult = await getUser();
      if (userResult instanceof ApiError) {
        authLogout();
      } else {
        updateUser(userResult);
        setAuthStatus(AuthStatus.Authenticated);
      }
    },
    [] // eslint-disable-line react-hooks/exhaustive-deps
  );

  const authFetch = useCallback<IAuthFetch>(
    async (params) => {
      if (authStatus === AuthStatus.Authenticated || params?.force) {
        const userResult = await getUser();

        if (!(userResult instanceof ApiError)) {
          updateUser(userResult);
        }
      }
    },
    [authStatus] // eslint-disable-line react-hooks/exhaustive-deps
  );

  const authLogout = useCallback(async () => {
    if (apiService.session !== undefined) {
      await logout(apiService.session);
    }

    setAuthStatus(AuthStatus.NotAuthenticated);
    updateUser(undefined);
    apiService.authorize(undefined);
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <AuthContext.Provider
      value={{
        user,
        authStatus,
        authLogin,
        authLogout,
        authFetch,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};
