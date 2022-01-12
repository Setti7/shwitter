import "@fontsource/signika";

import React, { useEffect } from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import {
  HOME_ROUTE,
  NEW_SHWEET_ROUTE,
  SIGN_IN_ROUTE,
  SIGN_UP_ROUTE,
} from "./config/routes";
import { AuthProvider } from "./contexts/auth";
import HomePage from "./pages/Home";
import NewShweetPage from "./pages/NewShweet";
import SignInPage from "./pages/SignIn";
import SignUpPage from "./pages/SignUp";
import { apiService } from "./services/api";
import dayjs from "dayjs";
import duration from "dayjs/plugin/duration";
import relativeTime from "dayjs/plugin/relativeTime";

function App() {
  useEffect(() => {
    apiService.initialize();
    dayjs.extend(duration);
    dayjs.extend(relativeTime);
  }, []);

  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path={HOME_ROUTE} element={<HomePage />} />
          <Route path={SIGN_IN_ROUTE} element={<SignInPage />} />
          <Route path={SIGN_UP_ROUTE} element={<SignUpPage />} />
          <Route path={NEW_SHWEET_ROUTE} element={<NewShweetPage />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}

export default App;
