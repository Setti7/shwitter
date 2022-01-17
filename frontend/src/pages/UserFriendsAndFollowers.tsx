import {
  Box,
  AppBar,
  Toolbar,
  IconButton,
  Typography,
  Container,
} from "@mui/material";
import { useNavigate, useParams } from "react-router-dom";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import { FC, useEffect } from "react";
import { HOME_ROUTE } from "../config/routes";
import FriendOrFollowerList from "../components/FriendOrFollowerList";
import { FriendOrFollowerType } from "../models/user";

export const UserFollowersPage = () => {
  return <UserFoFPage type={FriendOrFollowerType.Follower} />;
};

export const UserFriendsPage = () => {
  return <UserFoFPage type={FriendOrFollowerType.Friend} />;
};

interface Props {
  type: FriendOrFollowerType;
}

const UserFoFPage: FC<Props> = ({ type }) => {
  const navigate = useNavigate();
  const { userId } = useParams();

  useEffect(() => {
    if (userId === undefined) {
      navigate(HOME_ROUTE);
    }
  }, [userId]);

  if (userId === undefined) {
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

            <Typography variant="h6">
              {type === FriendOrFollowerType.Friend ? "Following" : "Followers"}
            </Typography>
          </Toolbar>
        </AppBar>
      </Box>

      <Container component="main" maxWidth="sm">
        <Box
          sx={{
            marginTop: 2,
            display: "flex",
            flexDirection: "column",
            alignItems: "stretch",
          }}
        >
          <FriendOrFollowerList type={type} />
        </Box>
      </Container>
    </>
  );
};
