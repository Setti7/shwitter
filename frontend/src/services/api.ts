import axios, { AxiosInstance } from "axios";
import Session from "../models/session";

import { API_ORIGIN } from "./../config/config";
import { loadSession, saveSession } from "./auth";

const HEADER = "X-Session-Token";

const api = axios.create({
  baseURL: API_ORIGIN,
  headers: {
    "Content-Type": "application/json",
    Accept: "application/json"
  },
});

enum Status {
  Idle,
  Initializing,
  Done,
}

class ApiService {
  private status = Status.Idle;
  private resolvers: ((args0: AxiosInstance) => void)[] = [];
  public session: Session | undefined;

  api = api;

  // Call this on app startup
  initialize = async (): Promise<void> => {
    if (this.status === Status.Idle) {
      this.status = Status.Initializing;

      this.authorize(loadSession());

      this.status = Status.Done;
      this.resolvers.forEach((r) => r(this.api));
    }
  };

  authorize = (value: Session | undefined) => {
    this.session = value;

    if (this.session) {
      this.api.defaults.headers.common[HEADER] = this.session.token;
    } else {
      delete this.api.defaults.headers.common[HEADER];
    }

    saveSession(this.session);
  };

  // Call this instead of direclty acessing the api instance
  getExecutor = async (): Promise<AxiosInstance> => {
    switch (this.status) {
      case Status.Done:
        return this.api;
      case Status.Initializing:
        return new Promise<AxiosInstance>((resolve) => {
          this.resolvers.push(resolve);
        });
      case Status.Idle:
        await this.initialize();
        return this.api;
    }
  };
}

export const apiService = new ApiService();
