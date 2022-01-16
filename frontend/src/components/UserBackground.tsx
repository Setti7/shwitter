import { Avatar, Box, Fab, IconButton, Snackbar } from "@mui/material";
import { FC, useCallback, useContext, useEffect, useState } from "react";
import { AuthContext } from "../contexts/auth";
import ApiError from "../models/errors/ApiError";
import { UserProfile } from "../models/user";
import { followUser, getIsFollowing, unFollowUser } from "../services/user";
import ShareIcon from "@mui/icons-material/Share";
import CloseIcon from "@mui/icons-material/Close";

const fabBlackStyle = {
  color: "common.white",
  bgcolor: "#000",
  border: "1px solid #666",
  "&:hover": {
    bgcolor: "#333",
  },
};

const UserBackground: FC<{ userProfile: UserProfile }> = ({ userProfile: user }) => {
  const { user: currentUser } = useContext(AuthContext);
  const [isFollowing, setIsFollowing] = useState(false);
  const [snackbarOpen, setSnackbarOpen] = useState(false);

  const handleClose = (
    event: React.SyntheticEvent | Event,
    reason?: string
  ) => {
    if (reason === "clickaway") {
      return;
    }

    setSnackbarOpen(false);
  };

  useEffect(() => {
    const getData = async () => {
      const result = await getIsFollowing(user.id);
      if (!(result instanceof ApiError)) {
        setIsFollowing(result);
      }
    };

    getData();
  }, [user]);

  const action = useCallback(async () => {
    if (isFollowing) {
      await unFollowUser(user.id);
      setIsFollowing(false);
    } else {
      followUser(user.id);
      setIsFollowing(true);
    }
  }, [user, isFollowing]);

  const copyToClipboard = useCallback(async () => {
    const url = window.location.origin + window.location.pathname;
    await navigator.clipboard.writeText(url);
    setSnackbarOpen(true);
  }, []);

  return (
    <>
      <Box position="relative" mb={6}>
        <img
          width={1500}
          height={500}
          style={{ maxWidth: "100%", height: "auto" }}
          alt=""
          // TODO: use user background image
          src="https://pbs.twimg.com/profile_banners/2451306812/1502508575/1500x500"
        ></img>
        <Box
          ml={2}
          width={80}
          height={80}
          sx={{
            background: "#000",
            borderRadius: "50%",
            position: "absolute",
            transform: "translateY(-50px)",
          }}
        >
          {/* TODO: use user profile picture */}
          <Avatar
            alt={user.name}
            sx={{
              width: 76,
              height: 76,
              position: "absolute",
              marginTop: "2px",
              marginLeft: "2px",
            }}
          />
        </Box>
        <Box sx={{ position: "absolute", right: 0 }}>
          <Box mr={2} flexGrow={1} display="flex">
            <Fab
              sx={fabBlackStyle}
              variant="extended"
              size="small"
              onClick={copyToClipboard}
            >
              <ShareIcon fontSize="small" />
            </Fab>

            <Box mr={1} />
            {currentUser?.id !== user.id ? (
              <Fab
                size="small"
                variant="extended"
                onClick={action}
                sx={isFollowing ? fabBlackStyle : undefined}
              >
                {isFollowing ? "Following" : "Follow"}
              </Fab>
            ) : undefined}
          </Box>
        </Box>
      </Box>
      <Snackbar
        open={snackbarOpen}
        autoHideDuration={6000}
        onClose={handleClose}
        message="User profile link copied to your clipboard."
        action={
          <IconButton
            size="small"
            aria-label="close"
            color="inherit"
            onClick={handleClose}
          >
            <CloseIcon fontSize="small" />
          </IconButton>
        }
      />
    </>
  );
};

export default UserBackground;
