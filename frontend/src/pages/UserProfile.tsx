import {
  Container,
  Box,
  Divider,
  Typography,
  AppBar,
  Toolbar,
  IconButton,
} from "@mui/material";
import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { getUserline } from "../services/shweets";
import { Timeline } from "../models/shweet";
import ApiError from "../models/errors/ApiError";
import ShweetCard from "../components/ShweetCard";
import { HOME_ROUTE } from "../config/routes";
import UserDetails from "../components/UserDetails";
import { getUserProfile } from "../services/user";
import { UserProfile } from "../models/user";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import UserBackground from "../components/UserBackground";

const UserProfilePage = () => {
  const navigate = useNavigate();
  const [error, setError] = useState<String | undefined>();
  const [userline, setUserline] = useState<Timeline>([]);
  const [profile, setProfile] = useState<UserProfile | undefined>();

  const { userId } = useParams();

  useEffect(() => {
    const getData = async () => {
      if (!userId) {
        navigate(HOME_ROUTE);
        return;
      }

      const [userResult, lineResult] = await Promise.all([
        getUserProfile(userId),
        getUserline(userId),
      ]);

      if (userResult instanceof ApiError) {
        setError(userResult.getFormattedStatus());
      } else {
        setProfile(userResult);
      }

      if (lineResult instanceof ApiError) {
        setError(lineResult.getFormattedStatus());
      } else {
        setUserline(lineResult);
      }
    };

    getData();
  }, [userId]);

  if (userId === undefined || profile === undefined) {
    // TODO loading
    return <></>;
  }

  return (
    <>
      <Box sx={{ flexGrow: 1 }}>
        <AppBar position="static">
          <Toolbar>
            <IconButton
              edge="start"
              color="inherit"
              aria-label="go back"
              sx={{ mr: 2 }}
              onClick={() => navigate(-1)}
            >
              <ArrowBackIcon />
            </IconButton>

            <Box display="flex" alignItems="baseline">
              <Box flexGrow={1}>
                <Typography variant="h6">{profile.name}</Typography>
              </Box>
              <Typography ml={2} variant="caption">
                {profile.shweets_count} shweets
              </Typography>
            </Box>
          </Toolbar>
        </AppBar>
      </Box>

      <Container maxWidth="sm" sx={{ padding: "0" }}>
          <UserBackground userProfile={profile} />
          <UserDetails userProfile={profile} />
          <Box mb={3} />
      </Container>

      <Container component="main" maxWidth="sm">
        <Box
          sx={{
            marginTop: 2,
            display: "flex",
            flexDirection: "column",
            alignItems: "stretch",
          }}
        >
          {error !== undefined ? (
            <Typography textAlign="center">{error}</Typography>
          ) : (
            <></>
          )}

          {userline.map((s) => {
            return (
              <Box mb={1}>
                <ShweetCard shweet={s} />
                <Divider sx={{ marginTop: 2 }} />
              </Box>
            );
          })}
        </Box>
      </Container>
    </>
  );
};

export default UserProfilePage;
