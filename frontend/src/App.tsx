import "@fontsource/signika";

import { useEffect } from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import {
  HOME_ROUTE,
  NEW_SHWEET_ROUTE,
  SHWEET_DETAILS_ROUTE,
  SIGN_IN_ROUTE,
  SIGN_UP_ROUTE,
  USER_FOLLOWERS_ROUTE,
  USER_FRIENDS_ROUTE,
  USER_ROUTE,
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
import UserProfilePage from "./pages/UserProfile";
import ShweetDetailsPage from "./pages/ShweetDetails";
import { UserFriendsPage, UserFollowersPage } from "./pages/UserFriendsAndFollowers";

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
          <Route path={USER_ROUTE} element={<UserProfilePage />} />
          <Route path={USER_FRIENDS_ROUTE} element={<UserFriendsPage />} />
          <Route path={USER_FOLLOWERS_ROUTE} element={<UserFollowersPage />} />
          <Route path={SHWEET_DETAILS_ROUTE} element={<ShweetDetailsPage />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}

export default App;
