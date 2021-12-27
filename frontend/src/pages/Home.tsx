import { Container, Box, Divider } from "@mui/material";
import { useContext, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import Header from "../components/Header";
import { AuthContext, AuthStatus } from "../contexts/auth";
import ShweetButton from "../components/ShweetButton";
import { SIGN_IN_ROUTE } from "../config/routes";
import { getTimeline } from "../services/shweets";
import { Timeline } from "../models/shweet";
import ApiError from "../models/errors/ApiError";
import ShweetCard from "../components/ShweetCard";

const HomePage = () => {
  const navigate = useNavigate();
  const [timeline, setTimeline] = useState<Timeline>([]);
  const { user, authStatus } = useContext(AuthContext);

  useEffect(() => {
    if (authStatus === AuthStatus.NotAuthenticated) {
      navigate(SIGN_IN_ROUTE);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [authStatus]);

  useEffect(() => {
    const getData = async () => {
      const tl = await getTimeline();
      if (tl instanceof ApiError) {
        // TODO show error
      } else {
        setTimeline(tl);
        console.log(tl);
      }
    };

    getData();
  }, []);

  if (user === undefined) {
    return <></>;
  }

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
          {/* TODO: 
          [ ] Fix timeline (we need to invert the order of shweets)
          [ ] Add userline and profile view
          */}

          {timeline.map((s) => {
            return (
              <Box mb={1}>
                <ShweetCard shweet={s} />
                <Divider sx={{marginTop: 2}} />
              </Box>
            );
          })}

          <ShweetButton />
        </Box>
      </Container>
    </>
  );
};

export default HomePage;
