import {
  Container,
  Box,
  Typography,
  AppBar,
  Toolbar,
  IconButton,
} from "@mui/material";
import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import ApiError from "../models/errors/ApiError";
import { HOME_ROUTE } from "../config/routes";
import UserDetails from "../components/UserDetails";
import { getUserProfile } from "../services/user";
import { UserProfile } from "../models/user";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import UserBackground from "../components/UserBackground";
import UserLine from "../components/UserLine";

const UserProfilePage = () => {
  const navigate = useNavigate();
  const [error, setError] = useState<String | undefined>();
  const [profile, setProfile] = useState<UserProfile | undefined>();

  const { userId } = useParams();

  useEffect(() => {
    const getData = async () => {
      if (!userId) {
        navigate(HOME_ROUTE);
        return;
      }

      const userResult = await getUserProfile(userId);

      if (userResult instanceof ApiError) {
        setError(userResult.getFormattedStatus());
      } else {
        setProfile(userResult);
      }
    };

    getData();
  }, [userId]);

  if (userId === undefined || profile === undefined) {
    // TODO loading
    return (
      <>
        {error !== undefined ? (
          <Typography textAlign="center">{error}</Typography>
        ) : (
          <></>
        )}
      </>
    );
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
          <UserLine userId={userId} />
        </Box>
      </Container>
    </>
  );
};

export default UserProfilePage;
