import { Container, Box, Divider, Typography } from "@mui/material";
import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import Header from "../components/Header";
import { getUserline } from "../services/shweets";
import { Timeline } from "../models/shweet";
import ApiError from "../models/errors/ApiError";
import ShweetCard from "../components/ShweetCard";
import { HOME_ROUTE } from "../config/routes";

const UserProfile = () => {
  const navigate = useNavigate();
  const [error, setError] = useState<String | undefined>();
  const [userline, setUserline] = useState<Timeline>([]);

  const { userId } = useParams();

  useEffect(() => {
    const getData = async () => {
      if (!userId) {
        navigate(HOME_ROUTE);
        return;
      }

      const result = await getUserline(userId);
      if (result instanceof ApiError) {
        setError(result.getFormattedStatus());
      } else {
        setUserline(result);
      }
    };

    getData();
  }, [userId]);

  if (userId === undefined) {
    return <></>;
  }
  // TODO:
  // [ ] Add user details/bio and follower/following count
  // [ ] Add follow/unfollow button

  return (
    <>
      <Header />
      <Container component="main">
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

export default UserProfile;
